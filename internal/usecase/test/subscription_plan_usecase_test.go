package usecase

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockSubscriptionPlanRepository struct {
	mock.Mock
}

func (m *MockSubscriptionPlanRepository) GetAll(req *domain.SubscriptionPlanListRequest) ([]*domain.SubscriptionPlan, int, error) {
	args := m.Called(req)
	return args.Get(0).([]*domain.SubscriptionPlan), args.Int(1), args.Error(2)
}

func (m *MockSubscriptionPlanRepository) GetByID(id string) (*domain.SubscriptionPlan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionPlanRepository) GetByKode(kode string) (*domain.SubscriptionPlan, error) {
	args := m.Called(kode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionPlanRepository) Create(plan *domain.SubscriptionPlan) error {
	args := m.Called(plan)
	return args.Error(0)
}

func (m *MockSubscriptionPlanRepository) Update(plan *domain.SubscriptionPlan) error {
	args := m.Called(plan)
	return args.Error(0)
}

func (m *MockSubscriptionPlanRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriptionPlanRepository) IsNameExists(nama string, excludeID ...string) (bool, error) {
	args := m.Called(nama, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSubscriptionPlanRepository) IsKodeExists(kode string, excludeID ...string) (bool, error) {
	args := m.Called(kode, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSubscriptionPlanRepository) CountActiveUsers(planID string) (int64, error) {
	args := m.Called(planID)
	return args.Get(0).(int64), args.Error(1)
}

func TestSubscriptionPlanUsecase_GetAll(t *testing.T) {
	mockRepo := new(MockSubscriptionPlanRepository)
	uc := usecase.NewSubscriptionPlanUsecase(mockRepo)

	t.Run("should return subscription plans successfully", func(t *testing.T) {
		req := &domain.SubscriptionPlanListRequest{
			Page:    1,
			PerPage: 10,
		}

		expectedPlans := []*domain.SubscriptionPlan{
			{
				ID:            uuid.New(),
				Nama:          "PRO Monthly",
				Harga:         99000,
				Interval:      "bulan",
				HariPercobaan: 7,
				Status:        "aktif",
			},
		}

		mockRepo.On("GetAll", req).Return(expectedPlans, 1, nil)

		plans, meta, err := uc.GetAll(req)

		assert.NoError(t, err)
		assert.Equal(t, expectedPlans, plans)
		assert.Equal(t, 1, meta.CurrentPage)
		assert.Equal(t, 1, meta.TotalRecords)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when repository fails", func(t *testing.T) {
		mockRepo2 := new(MockSubscriptionPlanRepository)
		uc2 := usecase.NewSubscriptionPlanUsecase(mockRepo2)

		req := &domain.SubscriptionPlanListRequest{
			Page:    1,
			PerPage: 10,
		}

		mockRepo2.On("GetAll", req).Return(([]*domain.SubscriptionPlan)(nil), 0, errors.New("database error"))

		plans, meta, err := uc2.GetAll(req)

		assert.Error(t, err)
		assert.Nil(t, plans)
		assert.Nil(t, meta)
		assert.Contains(t, err.Error(), "gagal mengambil daftar subscription plan")
		mockRepo2.AssertExpectations(t)
	})
}

func TestSubscriptionPlanUsecase_GetByID(t *testing.T) {
	mockRepo := new(MockSubscriptionPlanRepository)
	uc := usecase.NewSubscriptionPlanUsecase(mockRepo)

	t.Run("should return subscription plan successfully", func(t *testing.T) {
		planID := uuid.New().String()
		expectedPlan := &domain.SubscriptionPlan{
			ID:            uuid.MustParse(planID),
			Nama:          "PRO Monthly",
			Harga:         99000,
			Interval:      "bulan",
			HariPercobaan: 7,
			Status:        "aktif",
		}

		mockRepo.On("GetByID", planID).Return(expectedPlan, nil)

		plan, err := uc.GetByID(planID)

		assert.NoError(t, err)
		assert.Equal(t, expectedPlan, plan)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error for invalid ID format", func(t *testing.T) {
		invalidID := "invalid-uuid"

		plan, err := uc.GetByID(invalidID)

		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Contains(t, err.Error(), "format ID tidak valid")
	})

	t.Run("should return error when subscription plan not found", func(t *testing.T) {
		planID := uuid.New().String()

		mockRepo.On("GetByID", planID).Return(nil, gorm.ErrRecordNotFound)

		plan, err := uc.GetByID(planID)

		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Contains(t, err.Error(), "subscription plan tidak ditemukan")
		mockRepo.AssertExpectations(t)
	})
}

func TestSubscriptionPlanUsecase_Create(t *testing.T) {
	mockRepo := new(MockSubscriptionPlanRepository)
	uc := usecase.NewSubscriptionPlanUsecase(mockRepo)

	t.Run("should create subscription plan successfully", func(t *testing.T) {
		req := &domain.CreateSubscriptionPlanRequest{
			Nama:          "PRO Monthly",
			Harga:         99000,
			Interval:      "bulan",
			HariPercobaan: 7,
			Status:        "aktif",
		}

		mockRepo.On("IsNameExists", req.Nama, []string(nil)).Return(false, nil)
		mockRepo.On("IsKodeExists", mock.AnythingOfType("string"), []string(nil)).Return(false, nil)
		mockRepo.On("Create", mock.AnythingOfType("*domain.SubscriptionPlan")).Return(nil)

		plan, err := uc.Create(req)

		assert.NoError(t, err)
		assert.Equal(t, req.Nama, plan.Nama)
		assert.Equal(t, req.Harga, plan.Harga)
		assert.Equal(t, req.Interval, plan.Interval)
		assert.NotEmpty(t, plan.Kode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when name already exists", func(t *testing.T) {
		mockRepo2 := new(MockSubscriptionPlanRepository)
		uc2 := usecase.NewSubscriptionPlanUsecase(mockRepo2)

		req := &domain.CreateSubscriptionPlanRequest{
			Nama:          "PRO Monthly",
			Harga:         99000,
			Interval:      "bulan",
			HariPercobaan: 7,
			Status:        "aktif",
		}

		mockRepo2.On("IsNameExists", req.Nama, []string(nil)).Return(true, nil)

		plan, err := uc2.Create(req)

		assert.Error(t, err)
		assert.Nil(t, plan)
		assert.Contains(t, err.Error(), "subscription plan dengan nama tersebut sudah ada")
		mockRepo2.AssertExpectations(t)
	})
}

func TestSubscriptionPlanUsecase_Delete(t *testing.T) {
	mockRepo := new(MockSubscriptionPlanRepository)
	uc := usecase.NewSubscriptionPlanUsecase(mockRepo)

	t.Run("should delete subscription plan successfully", func(t *testing.T) {
		planID := uuid.New().String()
		existingPlan := &domain.SubscriptionPlan{
			ID:            uuid.MustParse(planID),
			Nama:          "PRO Monthly",
			Harga:         99000,
			Interval:      "bulan",
			HariPercobaan: 7,
			Status:        "aktif",
		}

		mockRepo.On("GetByID", planID).Return(existingPlan, nil)
		mockRepo.On("CountActiveUsers", planID).Return(int64(0), nil)
		mockRepo.On("Delete", planID).Return(nil)

		err := uc.Delete(planID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when subscription plan has active users", func(t *testing.T) {
		planID := uuid.New().String()
		existingPlan := &domain.SubscriptionPlan{
			ID:            uuid.MustParse(planID),
			Nama:          "PRO Monthly",
			Harga:         99000,
			Interval:      "bulan",
			HariPercobaan: 7,
			Status:        "aktif",
		}

		mockRepo.On("GetByID", planID).Return(existingPlan, nil)
		mockRepo.On("CountActiveUsers", planID).Return(int64(5), nil)

		err := uc.Delete(planID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "subscription plan tidak dapat dihapus karena sedang digunakan oleh pengguna aktif")
		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error for invalid ID format", func(t *testing.T) {
		invalidID := "invalid-uuid"

		err := uc.Delete(invalidID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "format ID tidak valid")
	})
}
