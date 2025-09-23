package usecase

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlanUsecase interface {
	GetAll(req *domain.SubscriptionPlanListRequest) ([]*domain.SubscriptionPlan, *domain.PaginationMeta, error)
	GetByID(id string) (*domain.SubscriptionPlan, error)
	Create(req *domain.CreateSubscriptionPlanRequest) (*domain.SubscriptionPlan, error)
	Update(id string, req *domain.UpdateSubscriptionPlanRequest) (*domain.SubscriptionPlan, error)
	Patch(id string, req *domain.PatchSubscriptionPlanRequest) (*domain.SubscriptionPlan, error)
	Delete(id string) error
}

type subscriptionPlanUsecase struct {
	subscriptionPlanRepo repo.SubscriptionPlanRepository
}

func NewSubscriptionPlanUsecase(
	subscriptionPlanRepo repo.SubscriptionPlanRepository,
) SubscriptionPlanUsecase {
	return &subscriptionPlanUsecase{
		subscriptionPlanRepo: subscriptionPlanRepo,
	}
}

func (uc *subscriptionPlanUsecase) GetAll(req *domain.SubscriptionPlanListRequest) ([]*domain.SubscriptionPlan, *domain.PaginationMeta, error) {
	req.SetDefaults()

	plans, total, err := uc.subscriptionPlanRepo.GetAll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal mengambil daftar subscription plan: %w", err)
	}

	totalPages := (total + req.PerPage - 1) / req.PerPage
	meta := &domain.PaginationMeta{
		CurrentPage:  req.Page,
		TotalPages:   totalPages,
		TotalRecords: total,
		PerPage:      req.PerPage,
	}

	return plans, meta, nil
}

func (uc *subscriptionPlanUsecase) GetByID(id string) (*domain.SubscriptionPlan, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("format ID tidak valid")
	}

	plan, err := uc.subscriptionPlanRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil subscription plan: %w", err)
	}

	return plan, nil
}

func (uc *subscriptionPlanUsecase) Create(req *domain.CreateSubscriptionPlanRequest) (*domain.SubscriptionPlan, error) {
	exists, err := uc.subscriptionPlanRepo.IsNameExists(req.Nama)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa nama subscription plan: %w", err)
	}
	if exists {
		return nil, errors.New("subscription plan dengan nama tersebut sudah ada")
	}

	plan := &domain.SubscriptionPlan{
		ID:            uuid.New(),
		Nama:          req.Nama,
		Harga:         req.Harga,
		Interval:      req.Interval,
		HariPercobaan: req.HariPercobaan,
		Status:        req.Status,
	}

	plan.Kode = uc.generateUniqueKode(req.Nama, req.Interval)

	if err := uc.subscriptionPlanRepo.Create(plan); err != nil {
		return nil, fmt.Errorf("gagal membuat subscription plan: %w", err)
	}

	return plan, nil
}

func (uc *subscriptionPlanUsecase) Update(id string, req *domain.UpdateSubscriptionPlanRequest) (*domain.SubscriptionPlan, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("format ID tidak valid")
	}

	existingPlan, err := uc.subscriptionPlanRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil subscription plan: %w", err)
	}

	if req.Nama != existingPlan.Nama {
		exists, err := uc.subscriptionPlanRepo.IsNameExists(req.Nama, id)
		if err != nil {
			return nil, fmt.Errorf("gagal memeriksa nama subscription plan: %w", err)
		}
		if exists {
			return nil, errors.New("subscription plan dengan nama tersebut sudah ada")
		}
	}

	existingPlan.Nama = req.Nama
	existingPlan.Harga = req.Harga
	existingPlan.Interval = req.Interval
	existingPlan.HariPercobaan = req.HariPercobaan
	existingPlan.Status = req.Status

	if req.Nama != existingPlan.Nama || req.Interval != existingPlan.Interval {
		existingPlan.Kode = uc.generateUniqueKode(req.Nama, req.Interval)
	}

	if err := uc.subscriptionPlanRepo.Update(existingPlan); err != nil {
		return nil, fmt.Errorf("gagal mengupdate subscription plan: %w", err)
	}

	return existingPlan, nil
}

func (uc *subscriptionPlanUsecase) Patch(id string, req *domain.PatchSubscriptionPlanRequest) (*domain.SubscriptionPlan, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("format ID tidak valid")
	}

	existingPlan, err := uc.subscriptionPlanRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subscription plan tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil subscription plan: %w", err)
	}

	needsKodeUpdate := false

	if req.Nama != nil && *req.Nama != existingPlan.Nama {
		exists, err := uc.subscriptionPlanRepo.IsNameExists(*req.Nama, id)
		if err != nil {
			return nil, fmt.Errorf("gagal memeriksa nama subscription plan: %w", err)
		}
		if exists {
			return nil, errors.New("subscription plan dengan nama tersebut sudah ada")
		}
		existingPlan.Nama = *req.Nama
		needsKodeUpdate = true
	}

	if req.Harga != nil {
		existingPlan.Harga = *req.Harga
	}

	if req.Interval != nil && *req.Interval != existingPlan.Interval {
		existingPlan.Interval = *req.Interval
		needsKodeUpdate = true
	}

	if req.HariPercobaan != nil {
		existingPlan.HariPercobaan = *req.HariPercobaan
	}

	if req.Status != nil {
		existingPlan.Status = *req.Status
	}

	if needsKodeUpdate {
		existingPlan.Kode = uc.generateUniqueKode(existingPlan.Nama, existingPlan.Interval)
	}

	if err := uc.subscriptionPlanRepo.Update(existingPlan); err != nil {
		return nil, fmt.Errorf("gagal mengupdate subscription plan: %w", err)
	}

	return existingPlan, nil
}

func (uc *subscriptionPlanUsecase) Delete(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return errors.New("format ID tidak valid")
	}

	_, err := uc.subscriptionPlanRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subscription plan tidak ditemukan")
		}
		return fmt.Errorf("gagal mengambil subscription plan: %w", err)
	}

	activeUsers, err := uc.subscriptionPlanRepo.CountActiveUsers(id)
	if err != nil {
		return fmt.Errorf("gagal memeriksa pengguna aktif: %w", err)
	}

	if activeUsers > 0 {
		return errors.New("subscription plan tidak dapat dihapus karena sedang digunakan oleh pengguna aktif")
	}

	if err := uc.subscriptionPlanRepo.Delete(id); err != nil {
		return fmt.Errorf("gagal menghapus subscription plan: %w", err)
	}

	return nil
}

func (uc *subscriptionPlanUsecase) generateUniqueKode(nama, interval string) string {
	baseKode := uc.generateBaseKode(nama, interval)

	for {
		randomSuffix := uc.generateRandomSuffix()
		kode := baseKode + "_" + randomSuffix

		exists, err := uc.subscriptionPlanRepo.IsKodeExists(kode)
		if err != nil || !exists {
			return kode
		}
	}
}

func (uc *subscriptionPlanUsecase) generateBaseKode(nama, interval string) string {
	namaUpper := strings.ToUpper(nama)
	namaUpper = strings.ReplaceAll(namaUpper, " ", "_")

	intervalUpper := strings.ToUpper(interval)
	if intervalUpper == "BULAN" {
		intervalUpper = "MONTHLY"
	} else if intervalUpper == "TAHUN" {
		intervalUpper = "YEARLY"
	}

	return namaUpper + "_" + intervalUpper
}

func (uc *subscriptionPlanUsecase) generateRandomSuffix() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 3)
	for i := range result {
		result[i] = charset[i%len(charset)]
	}
	return string(result)
}
