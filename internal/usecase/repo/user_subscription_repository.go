package repo

import (
	"fiber-boiler-plate/internal/domain"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type userSubscriptionRepository struct {
	db    *gorm.DB
	redis RedisRepository
}

func NewUserSubscriptionRepository(db *gorm.DB, redis RedisRepository) UserSubscriptionRepository {
	return &userSubscriptionRepository{
		db:    db,
		redis: redis,
	}
}

func (r *userSubscriptionRepository) GetAll(req *domain.UserSubscriptionListRequest) ([]*domain.UserSubscription, int, error) {
	cacheKey := fmt.Sprintf("user_subscription:list:page:%d:per_page:%d:sort:%s:%s",
		req.Page, req.PerPage, req.SortBy, req.SortDirection)

	if req.Search != nil && *req.Search != "" {
		cacheKey += fmt.Sprintf(":search:%s", *req.Search)
	}
	if req.Status != nil {
		cacheKey += fmt.Sprintf(":status:%s", *req.Status)
	}
	if req.PaymentMethod != nil {
		cacheKey += fmt.Sprintf(":payment_method:%s", *req.PaymentMethod)
	}

	var subscriptions []*domain.UserSubscription
	var total int64

	query := r.db.Model(&domain.UserSubscription{}).
		Preload("User").
		Preload("SubscriptionPlan")

	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		query = query.Joins("JOIN users ON users.id = user_subscriptions.user_id").
			Where("users.name ILIKE ? OR users.email ILIKE ?", searchTerm, searchTerm)
	}

	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	if req.PaymentMethod != nil {
		query = query.Where("payment_method = ?", *req.PaymentMethod)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orderClause string
	switch req.SortBy {
	case "nama":
		orderClause = fmt.Sprintf("users.name %s", strings.ToUpper(req.SortDirection))
		query = query.Joins("JOIN users ON users.id = user_subscriptions.user_id").
			Order(orderClause)
	case "email":
		orderClause = fmt.Sprintf("users.email %s", strings.ToUpper(req.SortDirection))
		query = query.Joins("JOIN users ON users.id = user_subscriptions.user_id").
			Order(orderClause)
	default:
		orderClause = fmt.Sprintf("%s %s", req.SortBy, strings.ToUpper(req.SortDirection))
		query = query.Order(orderClause)
	}

	offset := req.GetOffset()
	if err := query.Offset(offset).Limit(req.PerPage).Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, int(total), nil
}

func (r *userSubscriptionRepository) GetByID(id string) (*domain.UserSubscription, error) {
	cacheKey := fmt.Sprintf("user_subscription:id:%s", id)

	if r.redis != nil {
		var cachedSubscription domain.UserSubscription
		if err := r.redis.GetJSON(cacheKey, &cachedSubscription); err == nil && cachedSubscription.ID.String() != "" {
			return &cachedSubscription, nil
		}
	}

	var subscription domain.UserSubscription
	if err := r.db.Preload("User").
		Preload("SubscriptionPlan").
		Where("id = ?", id).
		First(&subscription).Error; err != nil {
		return nil, err
	}

	if r.redis != nil {
		go func() {
			_ = r.redis.SetJSON(cacheKey, subscription, 5*time.Minute)
		}()
	}

	return &subscription, nil
}

func (r *userSubscriptionRepository) UpdateStatus(id string, status string, reason *string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == "canceled" {
		updates["canceled_at"] = time.Now()
	}

	if status == "ended" {
		updates["ended_at"] = time.Now()
	}

	if err := r.db.Model(&domain.UserSubscription{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return err
	}

	r.clearRelatedCache(id)
	return nil
}

func (r *userSubscriptionRepository) UpdatePaymentMethod(id string, paymentMethod string) error {
	updates := map[string]interface{}{
		"payment_method": paymentMethod,
		"updated_at":     time.Now(),
	}

	if err := r.db.Model(&domain.UserSubscription{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return err
	}

	r.clearRelatedCache(id)
	return nil
}

func (r *userSubscriptionRepository) GetStatistics() (*domain.UserSubscriptionStatistics, error) {
	cacheKey := "user_subscription:statistics"

	if r.redis != nil {
		var cachedStats domain.UserSubscriptionStatistics
		if err := r.redis.GetJSON(cacheKey, &cachedStats); err == nil && cachedStats.TotalSubscriptions > 0 {
			return &cachedStats, nil
		}
	}

	var stats domain.UserSubscriptionStatistics

	if err := r.db.Model(&domain.UserSubscription{}).Count(&stats.TotalSubscriptions).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&domain.UserSubscription{}).Where("status = ?", "active").Count(&stats.ActiveSubscriptions).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&domain.UserSubscription{}).Where("status = ?", "paused").Count(&stats.PausedSubscriptions).Error; err != nil {
		return nil, err
	}

	if err := r.db.Model(&domain.UserSubscription{}).Where("status = ?", "trialing").Count(&stats.TrialingSubscriptions).Error; err != nil {
		return nil, err
	}

	var paymentMethodRows []struct {
		PaymentMethod string `json:"payment_method"`
		Count         int64  `json:"count"`
	}
	if err := r.db.Model(&domain.UserSubscription{}).
		Select("payment_method, COUNT(*) as count").
		Group("payment_method").
		Find(&paymentMethodRows).Error; err != nil {
		return nil, err
	}

	stats.PaymentMethods = make([]domain.PaymentMethodStatistic, len(paymentMethodRows))
	for i, row := range paymentMethodRows {
		percentage := float64(row.Count) / float64(stats.TotalSubscriptions) * 100
		stats.PaymentMethods[i] = domain.PaymentMethodStatistic{
			Method:     row.PaymentMethod,
			Count:      row.Count,
			Percentage: percentage,
		}
	}

	var monthlyRevenue float64
	if err := r.db.Model(&domain.UserSubscription{}).
		Joins("JOIN subscription_plans ON subscription_plans.id = user_subscriptions.subscription_plan_id").
		Where("user_subscriptions.status IN (?) AND subscription_plans.interval = ?", []string{"active", "trialing"}, "bulan").
		Select("COALESCE(SUM(subscription_plans.harga), 0)").
		Scan(&monthlyRevenue).Error; err != nil {
		return nil, err
	}

	var yearlyRevenue float64
	if err := r.db.Model(&domain.UserSubscription{}).
		Joins("JOIN subscription_plans ON subscription_plans.id = user_subscriptions.subscription_plan_id").
		Where("user_subscriptions.status IN (?) AND subscription_plans.interval = ?", []string{"active", "trialing"}, "tahun").
		Select("COALESCE(SUM(subscription_plans.harga), 0)").
		Scan(&yearlyRevenue).Error; err != nil {
		return nil, err
	}

	stats.MonthlyRevenue = monthlyRevenue
	stats.YearlyRevenue = yearlyRevenue + (monthlyRevenue * 12)

	if r.redis != nil {
		go func() {
			_ = r.redis.SetJSON(cacheKey, stats, 10*time.Minute)
		}()
	}

	return &stats, nil
}

func (r *userSubscriptionRepository) Create(subscription *domain.UserSubscription) error {
	if err := r.db.Create(subscription).Error; err != nil {
		return err
	}

	r.clearRelatedCache("")
	return nil
}

func (r *userSubscriptionRepository) Update(subscription *domain.UserSubscription) error {
	if err := r.db.Save(subscription).Error; err != nil {
		return err
	}

	r.clearRelatedCache(subscription.ID.String())
	return nil
}

func (r *userSubscriptionRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&domain.UserSubscription{}).Error; err != nil {
		return err
	}

	r.clearRelatedCache(id)
	return nil
}

func (r *userSubscriptionRepository) clearRelatedCache(subscriptionID string) {
	if r.redis != nil {
		patterns := []string{
			"user_subscription:list:*",
			"user_subscription:statistics",
		}

		if subscriptionID != "" {
			patterns = append(patterns, fmt.Sprintf("user_subscription:id:%s", subscriptionID))
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
