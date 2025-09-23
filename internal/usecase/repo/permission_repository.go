package repo

import (
	"fiber-boiler-plate/internal/domain"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db    *gorm.DB
	redis RedisRepository
}

func NewPermissionRepository(db *gorm.DB, redis RedisRepository) PermissionRepository {
	return &permissionRepository{
		db:    db,
		redis: redis,
	}
}

func (r *permissionRepository) GetAll(req *domain.PermissionListRequest) ([]*domain.Permission, int, error) {
	cacheKey := fmt.Sprintf("permission:list:page:%d:per_page:%d", req.Page, req.PerPage)

	if req.Search != nil && *req.Search != "" {
		cacheKey += fmt.Sprintf(":search:%s", *req.Search)
	}
	if req.Kategori != nil && *req.Kategori != "" {
		cacheKey += fmt.Sprintf(":kategori:%s", *req.Kategori)
	}

	var permissions []*domain.Permission
	var total int64

	if exists, _ := r.redis.Exists(cacheKey); exists {
		var cachedResult struct {
			Permissions []*domain.Permission `json:"permissions"`
			Total       int                  `json:"total"`
		}
		if err := r.redis.GetJSON(cacheKey, &cachedResult); err == nil {
			return cachedResult.Permissions, cachedResult.Total, nil
		}
	}

	query := r.db.Model(&domain.Permission{})

	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + strings.ToLower(*req.Search) + "%"
		query = query.Where("LOWER(nama) LIKE ?", searchTerm)
	}

	if req.Kategori != nil && *req.Kategori != "" {
		query = query.Where("kategori = ?", *req.Kategori)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PerPage
	if err := query.Order("nama ASC").Offset(offset).Limit(req.PerPage).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	cacheData := struct {
		Permissions []*domain.Permission `json:"permissions"`
		Total       int                  `json:"total"`
	}{
		Permissions: permissions,
		Total:       int(total),
	}

	r.redis.SetJSON(cacheKey, cacheData, 5*time.Minute)

	return permissions, int(total), nil
}

func (r *permissionRepository) GetByID(id string) (*domain.Permission, error) {
	cacheKey := fmt.Sprintf("permission:id:%s", id)

	var permission *domain.Permission
	if exists, _ := r.redis.Exists(cacheKey); exists {
		if err := r.redis.GetJSON(cacheKey, &permission); err == nil && permission != nil {
			return permission, nil
		}
	}

	permission = &domain.Permission{}
	err := r.db.Where("id = ?", id).First(permission).Error
	if err != nil {
		return nil, err
	}

	r.redis.SetJSON(cacheKey, permission, 10*time.Minute)

	return permission, nil
}

func (r *permissionRepository) GetByNama(nama string) (*domain.Permission, error) {
	cacheKey := fmt.Sprintf("permission:nama:%s", strings.ToLower(nama))

	var permission *domain.Permission
	if exists, _ := r.redis.Exists(cacheKey); exists {
		if err := r.redis.GetJSON(cacheKey, &permission); err == nil && permission != nil {
			return permission, nil
		}
	}

	permission = &domain.Permission{}
	err := r.db.Where("nama = ?", nama).First(permission).Error
	if err != nil {
		return nil, err
	}

	r.redis.SetJSON(cacheKey, permission, 10*time.Minute)

	return permission, nil
}

func (r *permissionRepository) Create(permission *domain.Permission) error {
	if err := r.db.Create(permission).Error; err != nil {
		return err
	}

	r.invalidateCache()
	return nil
}

func (r *permissionRepository) Update(permission *domain.Permission) error {
	if err := r.db.Save(permission).Error; err != nil {
		return err
	}

	r.invalidateCache()
	r.redis.Delete(fmt.Sprintf("permission:id:%s", permission.ID))
	r.redis.Delete(fmt.Sprintf("permission:nama:%s", strings.ToLower(permission.Nama)))

	return nil
}

func (r *permissionRepository) Delete(id string) error {
	permission, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if err := r.db.Delete(&domain.Permission{}, "id = ?", id).Error; err != nil {
		return err
	}

	r.invalidateCache()
	r.redis.Delete(fmt.Sprintf("permission:id:%s", id))
	r.redis.Delete(fmt.Sprintf("permission:nama:%s", strings.ToLower(permission.Nama)))

	return nil
}

func (r *permissionRepository) IsNameExists(nama string, excludeID ...string) (bool, error) {
	query := r.db.Model(&domain.Permission{}).Where("nama = ?", nama)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}

func (r *permissionRepository) IsUsedByRoles(id string) (bool, error) {
	var count int64
	err := r.db.Table("role_permissions").Where("permission_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *permissionRepository) invalidateCache() {
	keys, _ := r.redis.GetKeys("permission:list:*")
	for _, key := range keys {
		r.redis.Delete(key)
	}
}
