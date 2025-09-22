package repo

import (
	"crypto/rand"
	"fiber-boiler-plate/internal/domain"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type kantongRepository struct {
	db    *gorm.DB
	redis RedisRepository
}

func NewKantongRepository(db *gorm.DB, redis RedisRepository) KantongRepository {
	return &kantongRepository{
		db:    db,
		redis: redis,
	}
}

func (r *kantongRepository) GetByUserID(userID uint, req *domain.KantongListRequest) ([]*domain.Kantong, int, error) {
	cacheKey := fmt.Sprintf("kantong:list:user:%d:page:%d:per_page:%d:sort:%s:%s",
		userID, req.Page, req.PerPage, req.SortBy, req.SortDirection)

	if req.Search != nil && *req.Search != "" {
		cacheKey += fmt.Sprintf(":search:%s", *req.Search)
	}

	var kantongs []*domain.Kantong
	var total int64

	query := r.db.Where("user_id = ?", userID)

	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		query = query.Where("nama ILIKE ?", searchTerm)
	}

	if err := query.Model(&domain.Kantong{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := fmt.Sprintf("%s %s", req.SortBy, strings.ToUpper(req.SortDirection))
	offset := (req.Page - 1) * req.PerPage

	if err := query.Order(orderClause).Offset(offset).Limit(req.PerPage).Find(&kantongs).Error; err != nil {
		return nil, 0, err
	}

	if r.redis != nil {
		_ = r.redis.Set(cacheKey, kantongs, 5*time.Minute)
	}

	return kantongs, int(total), nil
}

func (r *kantongRepository) GetByID(id string, userID uint) (*domain.Kantong, error) {
	cacheKey := fmt.Sprintf("kantong:id:%s:user:%d", id, userID)

	var kantong domain.Kantong
	if r.redis != nil {
		if err := r.redis.GetJSON(cacheKey, &kantong); err == nil {
			return &kantong, nil
		}
	}

	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&kantong).Error; err != nil {
		return nil, err
	}

	if r.redis != nil {
		_ = r.redis.Set(cacheKey, kantong, 10*time.Minute)
	}

	return &kantong, nil
}

func (r *kantongRepository) GetByIDKartu(idKartu string, userID uint) (*domain.Kantong, error) {
	cacheKey := fmt.Sprintf("kantong:id_kartu:%s:user:%d", idKartu, userID)

	var kantong domain.Kantong
	if r.redis != nil {
		if err := r.redis.GetJSON(cacheKey, &kantong); err == nil {
			return &kantong, nil
		}
	}

	if err := r.db.Where("id_kartu = ? AND user_id = ?", idKartu, userID).First(&kantong).Error; err != nil {
		return nil, err
	}

	if r.redis != nil {
		_ = r.redis.Set(cacheKey, kantong, 10*time.Minute)
	}

	return &kantong, nil
}

func (r *kantongRepository) Create(kantong *domain.Kantong) error {
	if err := r.db.Create(kantong).Error; err != nil {
		return err
	}

	if r.redis != nil {
		cacheKey := fmt.Sprintf("kantong:id:%s:user:%d", kantong.ID, kantong.UserID)
		_ = r.redis.Set(cacheKey, kantong, 10*time.Minute)

		r.clearUserListCache(kantong.UserID)
	}

	return nil
}

func (r *kantongRepository) Update(kantong *domain.Kantong) error {
	if err := r.db.Save(kantong).Error; err != nil {
		return err
	}

	if r.redis != nil {
		cacheKey := fmt.Sprintf("kantong:id:%s:user:%d", kantong.ID, kantong.UserID)
		_ = r.redis.Set(cacheKey, kantong, 10*time.Minute)

		r.clearUserListCache(kantong.UserID)
	}

	return nil
}

func (r *kantongRepository) Delete(id string, userID uint) error {
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Kantong{}).Error; err != nil {
		return err
	}

	if r.redis != nil {
		cacheKey := fmt.Sprintf("kantong:id:%s:user:%d", id, userID)
		_ = r.redis.Delete(cacheKey)

		r.clearUserListCache(userID)
	}

	return nil
}

func (r *kantongRepository) IsNameExistForUser(nama string, userID uint, excludeID ...string) (bool, error) {
	query := r.db.Where("nama = ? AND user_id = ?", nama, userID)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Model(&domain.Kantong{}).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *kantongRepository) GenerateUniqueIDKartu() (string, error) {
	maxAttempts := 10

	for i := 0; i < maxAttempts; i++ {
		idKartu := r.generateRandomIDKartu()

		var count int64
		if err := r.db.Model(&domain.Kantong{}).Where("id_kartu = ?", idKartu).Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return idKartu, nil
		}
	}

	return "", fmt.Errorf("tidak dapat menghasilkan ID kartu unik setelah %d percobaan", maxAttempts)
}

func (r *kantongRepository) generateRandomIDKartu() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)

	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return r.generateFallbackIDKartu()
	}

	for i := range result {
		result[i] = charset[bytes[i]%byte(len(charset))]
	}

	return string(result)
}

func (r *kantongRepository) generateFallbackIDKartu() string {
	timestamp := time.Now().UnixNano()
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, 6)
	for i := range result {
		result[i] = charset[timestamp%int64(len(charset))]
		timestamp = timestamp / int64(len(charset))
	}

	return string(result)
}

func (r *kantongRepository) clearUserListCache(userID uint) {
	if r.redis == nil {
		return
	}
}

func (r *kantongRepository) Transfer(kantongAsalID, kantongTujuanID string, jumlah float64, userID uint) (*domain.Kantong, *domain.Kantong, error) {
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var kantongAsal, kantongTujuan domain.Kantong

	if err := tx.Where("id = ? AND user_id = ?", kantongAsalID, userID).First(&kantongAsal).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Where("id = ? AND user_id = ?", kantongTujuanID, userID).First(&kantongTujuan).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if kantongAsal.Saldo < jumlah {
		tx.Rollback()
		return nil, nil, fmt.Errorf("saldo tidak mencukupi")
	}

	kantongAsal.Saldo -= jumlah
	kantongTujuan.Saldo += jumlah
	kantongAsal.UpdatedAt = time.Now()
	kantongTujuan.UpdatedAt = time.Now()

	if err := tx.Save(&kantongAsal).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Save(&kantongTujuan).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	r.clearUserListCache(userID)

	return &kantongAsal, &kantongTujuan, nil
}
