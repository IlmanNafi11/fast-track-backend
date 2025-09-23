package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fiber-boiler-plate/internal/controller/http"
	"fiber-boiler-plate/internal/domain"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionPlanUsecase struct {
	mock.Mock
}

func (m *MockSubscriptionPlanUsecase) GetAll(req *domain.SubscriptionPlanListRequest) ([]*domain.SubscriptionPlan, *domain.PaginationMeta, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*domain.SubscriptionPlan), args.Get(1).(*domain.PaginationMeta), args.Error(2)
}

func (m *MockSubscriptionPlanUsecase) GetByID(id string) (*domain.SubscriptionPlan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionPlanUsecase) Create(req *domain.CreateSubscriptionPlanRequest) (*domain.SubscriptionPlan, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionPlanUsecase) Update(id string, req *domain.UpdateSubscriptionPlanRequest) (*domain.SubscriptionPlan, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionPlanUsecase) Patch(id string, req *domain.PatchSubscriptionPlanRequest) (*domain.SubscriptionPlan, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.SubscriptionPlan), args.Error(1)
}

func (m *MockSubscriptionPlanUsecase) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestSubscriptionPlanController_GetAll(t *testing.T) {
	t.Run("should return subscription plans successfully", func(t *testing.T) {
		app := fiber.New()
		mockUsecase := new(MockSubscriptionPlanUsecase)
		controller := http.NewSubscriptionPlanController(mockUsecase)
		app.Get("/subscription-plans", controller.GetAll)

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

		expectedMeta := &domain.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   1,
			TotalRecords: 1,
			PerPage:      10,
		}

		mockUsecase.On("GetAll", mock.AnythingOfType("*domain.SubscriptionPlanListRequest")).Return(expectedPlans, expectedMeta, nil)

		req := httptest.NewRequest("GET", "/subscription-plans", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return error when usecase fails", func(t *testing.T) {
		app := fiber.New()
		mockUsecase := new(MockSubscriptionPlanUsecase)
		controller := http.NewSubscriptionPlanController(mockUsecase)
		app.Get("/subscription-plans", controller.GetAll)

		mockUsecase.On("GetAll", mock.AnythingOfType("*domain.SubscriptionPlanListRequest")).Return(([]*domain.SubscriptionPlan)(nil), (*domain.PaginationMeta)(nil), errors.New("usecase error"))

		req := httptest.NewRequest("GET", "/subscription-plans", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestSubscriptionPlanController_GetByID(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockSubscriptionPlanUsecase)
	controller := http.NewSubscriptionPlanController(mockUsecase)

	app.Get("/subscription-plans/:id", controller.GetByID)

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

		mockUsecase.On("GetByID", planID).Return(expectedPlan, nil)

		req := httptest.NewRequest("GET", "/subscription-plans/"+planID, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 404 when subscription plan not found", func(t *testing.T) {
		planID := uuid.New().String()

		mockUsecase.On("GetByID", planID).Return(nil, errors.New("subscription plan tidak ditemukan"))

		req := httptest.NewRequest("GET", "/subscription-plans/"+planID, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestSubscriptionPlanController_Create(t *testing.T) {
	t.Run("should create subscription plan successfully", func(t *testing.T) {
		app := fiber.New()
		mockUsecase := new(MockSubscriptionPlanUsecase)
		controller := http.NewSubscriptionPlanController(mockUsecase)
		app.Post("/subscription-plans", controller.Create)

		requestBody := domain.CreateSubscriptionPlanRequest{
			Nama:          "PRO Monthly",
			Harga:         99000,
			Interval:      "bulan",
			HariPercobaan: 7,
			Status:        "aktif",
		}

		expectedPlan := &domain.SubscriptionPlan{
			ID:            uuid.New(),
			Nama:          requestBody.Nama,
			Harga:         requestBody.Harga,
			Interval:      requestBody.Interval,
			HariPercobaan: requestBody.HariPercobaan,
			Status:        requestBody.Status,
		}

		mockUsecase.On("Create", mock.AnythingOfType("*domain.CreateSubscriptionPlanRequest")).Return(expectedPlan, nil)

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/subscription-plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 409 when name already exists", func(t *testing.T) {
		app := fiber.New()
		mockUsecase := new(MockSubscriptionPlanUsecase)
		controller := http.NewSubscriptionPlanController(mockUsecase)
		app.Post("/subscription-plans", controller.Create)

		requestBody := domain.CreateSubscriptionPlanRequest{
			Nama:          "PRO Monthly",
			Harga:         99000,
			Interval:      "bulan",
			HariPercobaan: 7,
			Status:        "aktif",
		}

		mockUsecase.On("Create", mock.AnythingOfType("*domain.CreateSubscriptionPlanRequest")).Return((*domain.SubscriptionPlan)(nil), errors.New("subscription plan dengan nama tersebut sudah ada"))

		body, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/subscription-plans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestSubscriptionPlanController_Delete(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(MockSubscriptionPlanUsecase)
	controller := http.NewSubscriptionPlanController(mockUsecase)

	app.Delete("/subscription-plans/:id", controller.Delete)

	t.Run("should delete subscription plan successfully", func(t *testing.T) {
		planID := uuid.New().String()

		mockUsecase.On("Delete", planID).Return(nil)

		req := httptest.NewRequest("DELETE", "/subscription-plans/"+planID, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 409 when subscription plan has active users", func(t *testing.T) {
		planID := uuid.New().String()

		mockUsecase.On("Delete", planID).Return(errors.New("subscription plan tidak dapat dihapus karena sedang digunakan oleh pengguna aktif"))

		req := httptest.NewRequest("DELETE", "/subscription-plans/"+planID, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}
