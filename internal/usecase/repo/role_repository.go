package repo

import (
	"fiber-boiler-plate/internal/domain"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type roleRepository struct {
	db    *gorm.DB
	redis RedisRepository
}

func NewRoleRepository(db *gorm.DB, redis RedisRepository) RoleRepository {
	return &roleRepository{
		db:    db,
		redis: redis,
	}
}

func (r *roleRepository) GetAll(req *domain.RoleListRequest) ([]*domain.Role, int, error) {
	cacheKey := fmt.Sprintf("role:list:page:%d:per_page:%d", req.Page, req.PerPage)

	if req.Search != nil && *req.Search != "" {
		cacheKey += fmt.Sprintf(":search:%s", *req.Search)
	}
	if req.Status != nil && *req.Status != "" {
		cacheKey += fmt.Sprintf(":status:%s", *req.Status)
	}

	var roles []*domain.Role
	var total int64

	if exists, _ := r.redis.Exists(cacheKey); exists {
		var cachedResult struct {
			Roles []*domain.Role `json:"roles"`
			Total int            `json:"total"`
		}
		if err := r.redis.GetJSON(cacheKey, &cachedResult); err == nil {
			return cachedResult.Roles, cachedResult.Total, nil
		}
	}

	query := r.db.Model(&domain.Role{}).Preload("Permissions")

	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + strings.ToLower(*req.Search) + "%"
		query = query.Where("LOWER(nama) LIKE ?", searchTerm)
	}

	if req.Status != nil && *req.Status != "" {
		query = query.Where("status = ?", *req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PerPage
	if err := query.Order("nama ASC").Offset(offset).Limit(req.PerPage).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	cacheData := struct {
		Roles []*domain.Role `json:"roles"`
		Total int            `json:"total"`
	}{
		Roles: roles,
		Total: int(total),
	}

	r.redis.SetJSON(cacheKey, cacheData, 5*time.Minute)

	return roles, int(total), nil
}

func (r *roleRepository) GetByID(id string) (*domain.Role, error) {
	cacheKey := fmt.Sprintf("role:id:%s", id)

	var role domain.Role
	if exists, _ := r.redis.Exists(cacheKey); exists {
		if err := r.redis.GetJSON(cacheKey, &role); err == nil {
			return &role, nil
		}
	}

	if err := r.db.Preload("Permissions").First(&role, "id = ?", id).Error; err != nil {
		return nil, err
	}

	r.redis.SetJSON(cacheKey, role, 10*time.Minute)

	return &role, nil
}

func (r *roleRepository) GetByNama(nama string) (*domain.Role, error) {
	var role domain.Role
	if err := r.db.Preload("Permissions").First(&role, "nama = ?", nama).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Create(role *domain.Role) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(role).Error; err != nil {
			return err
		}

		if len(role.Permissions) > 0 {
			permissionIDs := make([]string, len(role.Permissions))
			for i, perm := range role.Permissions {
				permissionIDs[i] = perm.ID
			}

			if err := r.updateRolePermissionsInTx(tx, role.ID, permissionIDs); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	r.invalidateCache()
	return nil
}

func (r *roleRepository) Update(role *domain.Role) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("nama", "deskripsi", "status", "updated_at").Updates(role).Error; err != nil {
			return err
		}

		permissionIDs := make([]string, len(role.Permissions))
		for i, perm := range role.Permissions {
			permissionIDs[i] = perm.ID
		}

		if err := r.updateRolePermissionsInTx(tx, role.ID, permissionIDs); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	r.invalidateCache()
	r.redis.Delete(fmt.Sprintf("role:id:%s", role.ID))
	return nil
}

func (r *roleRepository) Delete(id string) error {
	if err := r.db.Select("Permissions").Delete(&domain.Role{}, "id = ?", id).Error; err != nil {
		return err
	}

	r.invalidateCache()
	r.redis.Delete(fmt.Sprintf("role:id:%s", id))
	return nil
}

func (r *roleRepository) IsNameExists(nama string, excludeID ...string) (bool, error) {
	query := r.db.Model(&domain.Role{}).Where("nama = ?", nama)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *roleRepository) IsUsedByUsers(id string) (bool, error) {
	var count int64
	if err := r.db.Table("users").Where("role_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *roleRepository) GetRolePermissions(roleID string, req *domain.RolePermissionListRequest) ([]*domain.Permission, int, error) {
	cacheKey := fmt.Sprintf("role:permissions:%s:page:%d:per_page:%d", roleID, req.Page, req.PerPage)

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

	query := r.db.Model(&domain.Permission{}).
		Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (req.Page - 1) * req.PerPage
	if err := query.Order("permissions.nama ASC").Offset(offset).Limit(req.PerPage).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	cacheData := struct {
		Permissions []*domain.Permission `json:"permissions"`
		Total       int                  `json:"total"`
	}{
		Permissions: permissions,
		Total:       int(total),
	}

	r.redis.SetJSON(cacheKey, cacheData, 10*time.Minute)

	return permissions, int(total), nil
}

func (r *roleRepository) UpdateRolePermissions(roleID string, permissionIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return r.updateRolePermissionsInTx(tx, roleID, permissionIDs)
	})
}

func (r *roleRepository) updateRolePermissionsInTx(tx *gorm.DB, roleID string, permissionIDs []string) error {
	if err := tx.Where("role_id = ?", roleID).Delete(&domain.RolePermission{}).Error; err != nil {
		return err
	}

	if len(permissionIDs) > 0 {
		for _, permissionID := range permissionIDs {
			rolePermission := domain.RolePermission{
				RoleID:       roleID,
				PermissionID: permissionID,
				CreatedAt:    time.Now(),
			}
			if err := tx.Create(&rolePermission).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *roleRepository) ValidatePermissions(permissionIDs []string) ([]string, error) {
	var existingIDs []string

	if err := r.db.Model(&domain.Permission{}).
		Where("id IN ?", permissionIDs).
		Pluck("id", &existingIDs).Error; err != nil {
		return nil, err
	}

	invalidIDs := make([]string, 0)
	for _, id := range permissionIDs {
		found := false
		for _, existing := range existingIDs {
			if id == existing {
				found = true
				break
			}
		}
		if !found {
			invalidIDs = append(invalidIDs, id)
		}
	}

	return invalidIDs, nil
}

func (r *roleRepository) invalidateCache() {
	keys, _ := r.redis.GetKeys("role:list:*")
	for _, key := range keys {
		r.redis.Delete(key)
	}

	permissionKeys, _ := r.redis.GetKeys("role:permissions:*")
	for _, key := range permissionKeys {
		r.redis.Delete(key)
	}
}
