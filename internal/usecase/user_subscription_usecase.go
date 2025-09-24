package usecase

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"fmt"

	"gorm.io/gorm"
)

type UserSubscriptionUsecase interface {
	GetAll(req *domain.UserSubscriptionListRequest) ([]*domain.UserSubscriptionResponse, *domain.PaginationMeta, error)
	GetByID(id string) (*domain.UserSubscriptionDetailResponse, error)
	UpdateStatus(id string, req *domain.UpdateUserSubscriptionRequest) (*domain.UserSubscriptionDetailResponse, error)
	UpdatePaymentMethod(id string, req *domain.UpdatePaymentMethodRequest) (*domain.UserSubscriptionDetailResponse, error)
	GetStatistics() (*domain.UserSubscriptionStatistics, error)
}

type userSubscriptionUsecase struct {
	userSubscriptionRepo repo.UserSubscriptionRepository
	userRepo             repo.UserRepository
	subscriptionPlanRepo repo.SubscriptionPlanRepository
}

func NewUserSubscriptionUsecase(
	userSubscriptionRepo repo.UserSubscriptionRepository,
	userRepo repo.UserRepository,
	subscriptionPlanRepo repo.SubscriptionPlanRepository,
) UserSubscriptionUsecase {
	return &userSubscriptionUsecase{
		userSubscriptionRepo: userSubscriptionRepo,
		userRepo:             userRepo,
		subscriptionPlanRepo: subscriptionPlanRepo,
	}
}

func (uc *userSubscriptionUsecase) GetAll(req *domain.UserSubscriptionListRequest) ([]*domain.UserSubscriptionResponse, *domain.PaginationMeta, error) {
	req.SetDefaults()

	subscriptions, total, err := uc.userSubscriptionRepo.GetAll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal mengambil daftar subscription pengguna: %w", err)
	}

	totalPages := (total + req.PerPage - 1) / req.PerPage
	meta := &domain.PaginationMeta{
		CurrentPage:  req.Page,
		TotalPages:   totalPages,
		TotalRecords: total,
		PerPage:      req.PerPage,
	}

	responses := make([]*domain.UserSubscriptionResponse, len(subscriptions))
	for i, subscription := range subscriptions {
		responses[i] = uc.mapToResponse(subscription)
	}

	return responses, meta, nil
}

func (uc *userSubscriptionUsecase) GetByID(id string) (*domain.UserSubscriptionDetailResponse, error) {
	subscription, err := uc.userSubscriptionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subscription pengguna tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil detail subscription pengguna: %w", err)
	}

	return uc.mapToDetailResponse(subscription), nil
}

func (uc *userSubscriptionUsecase) UpdateStatus(id string, req *domain.UpdateUserSubscriptionRequest) (*domain.UserSubscriptionDetailResponse, error) {
	subscription, err := uc.userSubscriptionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subscription pengguna tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil subscription pengguna: %w", err)
	}

	if subscription.Status == "canceled" || subscription.Status == "ended" {
		return nil, fmt.Errorf("subscription yang sudah dibatalkan atau berakhir tidak dapat diubah")
	}

	var newStatus string
	switch req.Action {
	case "pause":
		if subscription.Status == "paused" {
			return nil, fmt.Errorf("subscription sudah dalam status pause")
		}
		newStatus = "paused"
	case "activate":
		if subscription.Status == "active" {
			return nil, fmt.Errorf("subscription sudah dalam status active")
		}
		newStatus = "active"
	case "cancel":
		newStatus = "canceled"
	default:
		return nil, fmt.Errorf("action tidak valid: %s", req.Action)
	}

	if err := uc.userSubscriptionRepo.UpdateStatus(id, newStatus, req.Reason); err != nil {
		return nil, fmt.Errorf("gagal mengupdate status subscription: %w", err)
	}

	updatedSubscription, err := uc.userSubscriptionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil subscription yang sudah diupdate: %w", err)
	}

	return uc.mapToDetailResponse(updatedSubscription), nil
}

func (uc *userSubscriptionUsecase) UpdatePaymentMethod(id string, req *domain.UpdatePaymentMethodRequest) (*domain.UserSubscriptionDetailResponse, error) {
	subscription, err := uc.userSubscriptionRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("subscription pengguna tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil subscription pengguna: %w", err)
	}

	if subscription.Status == "canceled" || subscription.Status == "ended" {
		return nil, fmt.Errorf("subscription yang sudah dibatalkan atau berakhir tidak dapat diubah")
	}

	if err := uc.userSubscriptionRepo.UpdatePaymentMethod(id, req.PaymentMethod); err != nil {
		return nil, fmt.Errorf("gagal mengupdate metode pembayaran: %w", err)
	}

	updatedSubscription, err := uc.userSubscriptionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil subscription yang sudah diupdate: %w", err)
	}

	return uc.mapToDetailResponse(updatedSubscription), nil
}

func (uc *userSubscriptionUsecase) GetStatistics() (*domain.UserSubscriptionStatistics, error) {
	stats, err := uc.userSubscriptionRepo.GetStatistics()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil statistik subscription pengguna: %w", err)
	}

	return stats, nil
}

func (uc *userSubscriptionUsecase) mapToResponse(subscription *domain.UserSubscription) *domain.UserSubscriptionResponse {
	return &domain.UserSubscriptionResponse{
		ID: subscription.ID,
		User: domain.UserInfo{
			ID:    subscription.User.ID,
			Nama:  subscription.User.Name,
			Email: subscription.User.Email,
			Role:  "User",
		},
		SubscriptionPlan: domain.SubscriptionPlanInfo{
			ID:    subscription.SubscriptionPlan.ID,
			Nama:  subscription.SubscriptionPlan.Nama,
			Harga: subscription.SubscriptionPlan.Harga,
		},
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		PaymentMethod:      subscription.PaymentMethod,
		CreatedAt:          subscription.CreatedAt,
		UpdatedAt:          subscription.UpdatedAt,
	}
}

func (uc *userSubscriptionUsecase) mapToDetailResponse(subscription *domain.UserSubscription) *domain.UserSubscriptionDetailResponse {
	return &domain.UserSubscriptionDetailResponse{
		ID: subscription.ID,
		User: domain.UserInfoDetail{
			ID:        subscription.User.ID,
			Nama:      subscription.User.Name,
			Email:     subscription.User.Email,
			Role:      "User",
			IsActive:  subscription.User.IsActive,
			CreatedAt: subscription.User.CreatedAt,
		},
		SubscriptionPlan: domain.SubscriptionPlanDetail{
			ID:            subscription.SubscriptionPlan.ID,
			Kode:          subscription.SubscriptionPlan.Kode,
			Nama:          subscription.SubscriptionPlan.Nama,
			Harga:         subscription.SubscriptionPlan.Harga,
			Interval:      subscription.SubscriptionPlan.Interval,
			HariPercobaan: subscription.SubscriptionPlan.HariPercobaan,
			Status:        subscription.SubscriptionPlan.Status,
		},
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		TrialEnd:           subscription.TrialEnd,
		PaymentMethod:      subscription.PaymentMethod,
		PaymentStatus:      subscription.PaymentStatus,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
		CanceledAt:         subscription.CanceledAt,
		EndedAt:            subscription.EndedAt,
		CreatedAt:          subscription.CreatedAt,
		UpdatedAt:          subscription.UpdatedAt,
	}
}
