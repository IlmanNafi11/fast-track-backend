package repo

import (
	"fiber-boiler-plate/internal/domain"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type subscriptionPlanRepository struct {
	db    *gorm.DB
	redis RedisRepository
}

func NewSubscriptionPlanRepository(db *gorm.DB, redis RedisRepository) SubscriptionPlanRepository {
	return &subscriptionPlanRepository{
		db:    db,
		redis: redis,
	}
}

func (r *subscriptionPlanRepository) GetAll(req *domain.SubscriptionPlanListRequest) ([]*domain.SubscriptionPlan, int, error) {
	cacheKey := fmt.Sprintf("subscription_plan:list:page:%d:per_page:%d:sort:%s:%s",
		req.Page, req.PerPage, req.SortBy, req.SortDirection)

	if req.Search != "" {
		cacheKey += fmt.Sprintf(":search:%s", req.Search)
	}

	var plans []*domain.SubscriptionPlan
	var total int64

	query := r.db.Model(&domain.SubscriptionPlan{})

	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where("nama ILIKE ?", searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderClause := fmt.Sprintf("%s %s", req.SortBy, strings.ToUpper(req.SortDirection))
	offset := req.GetOffset()

	if err := query.Order(orderClause).Offset(offset).Limit(req.PerPage).Find(&plans).Error; err != nil {
		return nil, 0, err
	}

	if r.redis != nil {
		go func() {
			_ = r.redis.SetJSON(cacheKey, map[string]interface{}{
				"plans": plans,
				"total": total,
			}, 10*time.Minute)
		}()
	}

	return plans, int(total), nil
}

func (r *subscriptionPlanRepository) GetByID(id string) (*domain.SubscriptionPlan, error) {
	cacheKey := fmt.Sprintf("subscription_plan:id:%s", id)

	var plan domain.SubscriptionPlan
	if r.redis != nil {
		if err := r.redis.GetJSON(cacheKey, &plan); err == nil {
			return &plan, nil
		}
	}

	if err := r.db.Where("id = ?", id).First(&plan).Error; err != nil {
		return nil, err
	}

	if r.redis != nil {
		go func() {
			_ = r.redis.SetJSON(cacheKey, plan, 30*time.Minute)
		}()
	}

	return &plan, nil
}

func (r *subscriptionPlanRepository) GetByKode(kode string) (*domain.SubscriptionPlan, error) {
	cacheKey := fmt.Sprintf("subscription_plan:kode:%s", kode)

	var plan domain.SubscriptionPlan
	if r.redis != nil {
		if err := r.redis.GetJSON(cacheKey, &plan); err == nil {
			return &plan, nil
		}
	}

	if err := r.db.Where("kode = ?", kode).First(&plan).Error; err != nil {
		return nil, err
	}

	if r.redis != nil {
		go func() {
			_ = r.redis.SetJSON(cacheKey, plan, 30*time.Minute)
		}()
	}

	return &plan, nil
}

func (r *subscriptionPlanRepository) Create(plan *domain.SubscriptionPlan) error {
	if err := r.db.Create(plan).Error; err != nil {
		return err
	}

	if r.redis != nil {
		go func() {
			r.clearRelatedCache(plan.ID.String())
		}()
	}

	return nil
}

func (r *subscriptionPlanRepository) Update(plan *domain.SubscriptionPlan) error {
	if err := r.db.Save(plan).Error; err != nil {
		return err
	}

	if r.redis != nil {
		go func() {
			r.clearRelatedCache(plan.ID.String())
		}()
	}

	return nil
}

func (r *subscriptionPlanRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&domain.SubscriptionPlan{}).Error; err != nil {
		return err
	}

	if r.redis != nil {
		go func() {
			r.clearRelatedCache(id)
		}()
	}

	return nil
}

func (r *subscriptionPlanRepository) IsNameExists(nama string, excludeID ...string) (bool, error) {
	query := r.db.Model(&domain.SubscriptionPlan{}).Where("nama = ?", nama)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *subscriptionPlanRepository) IsKodeExists(kode string, excludeID ...string) (bool, error) {
	query := r.db.Model(&domain.SubscriptionPlan{}).Where("kode = ?", kode)

	if len(excludeID) > 0 && excludeID[0] != "" {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *subscriptionPlanRepository) CountActiveUsers(planID string) (int64, error) {
	cacheKey := fmt.Sprintf("subscription_plan:active_users:%s", planID)

	if r.redis != nil {
		if cachedCount, err := r.redis.Get(cacheKey); err == nil && cachedCount != "" {
			var count int64
			if err := r.redis.GetJSON(cacheKey, &count); err == nil {
				return count, nil
			}
		}
	}

	var count int64
	if err := r.db.Table("user_subscriptions").
		Where("subscription_plan_id = ? AND status = 'active'", planID).
		Count(&count).Error; err != nil {
		return 0, err
	}

	if r.redis != nil {
		go func() {
			_ = r.redis.SetJSON(cacheKey, count, 5*time.Minute)
		}()
	}

	return count, nil
}

func (r *subscriptionPlanRepository) clearRelatedCache(planID string) {
	if r.redis != nil {
		patterns := []string{
			"subscription_plan:list:*",
			fmt.Sprintf("subscription_plan:id:%s", planID),
			fmt.Sprintf("subscription_plan:kode:*"),
			fmt.Sprintf("subscription_plan:active_users:%s", planID),
		}

		for _, pattern := range patterns {
			if keys, err := r.redis.GetKeys(pattern); err == nil {
				for _, key := range keys {
					_ = r.redis.Delete(key)
				}
			}
		}
	}
}
