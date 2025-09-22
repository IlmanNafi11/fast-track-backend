package http_test

import (
	"bytes"
	"encoding/json"
	"fiber-boiler-plate/internal/controller/http"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockKantongUsecase struct {
	mock.Mock
}

func (m *MockKantongUsecase) GetKantongList(userID uint, req *domain.KantongListRequest) ([]*domain.KantongResponse, *domain.PaginationMeta, error) {
	args := m.Called(userID, req)
	return args.Get(0).([]*domain.KantongResponse), args.Get(1).(*domain.PaginationMeta), args.Error(2)
}

func (m *MockKantongUsecase) GetKantongByID(id string, userID uint) (*domain.KantongResponse, error) {
	args := m.Called(id, userID)
	return args.Get(0).(*domain.KantongResponse), args.Error(1)
}

func (m *MockKantongUsecase) CreateKantong(req *domain.CreateKantongRequest, userID uint) (*domain.KantongResponse, error) {
	args := m.Called(req, userID)
	return args.Get(0).(*domain.KantongResponse), args.Error(1)
}

func (m *MockKantongUsecase) UpdateKantong(id string, req *domain.UpdateKantongRequest, userID uint) (*domain.KantongResponse, error) {
	args := m.Called(id, req, userID)
	return args.Get(0).(*domain.KantongResponse), args.Error(1)
}

func (m *MockKantongUsecase) PatchKantong(id string, req *domain.PatchKantongRequest, userID uint) (*domain.KantongResponse, error) {
	args := m.Called(id, req, userID)
	return args.Get(0).(*domain.KantongResponse), args.Error(1)
}

func (m *MockKantongUsecase) DeleteKantong(id string, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

func (m *MockKantongUsecase) SetAnggaranUsecase(anggaranUsecase usecase.AnggaranUsecase) {
	// Mock method - tidak perlu implementasi
}

func setupKantongController() (*fiber.App, *MockKantongUsecase) {
	app := fiber.New()
	mockUsecase := new(MockKantongUsecase)
	controller := http.NewKantongController(mockUsecase)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return c.Next()
	})

	app.Get("/kantong", controller.GetKantongList)
	app.Get("/kantong/:id", controller.GetKantongByID)
	app.Post("/kantong", controller.CreateKantong)
	app.Put("/kantong/:id", controller.UpdateKantong)
	app.Patch("/kantong/:id", controller.PatchKantong)
	app.Delete("/kantong/:id", controller.DeleteKantong)

	return app, mockUsecase
}

func TestGetKantongList_Success(t *testing.T) {
	app, mockUsecase := setupKantongController()

	expectedKantongs := []*domain.KantongResponse{
		{
			ID:       "test-id-1",
			IDKartu:  "ABC123",
			Nama:     "Kantong Utama",
			Kategori: "Pengeluaran",
			Saldo:    100000,
			Warna:    "Navy",
		},
	}

	expectedMeta := &domain.PaginationMeta{
		CurrentPage:  1,
		PerPage:      10,
		TotalRecords: 1,
		TotalPages:   1,
	}

	mockUsecase.On("GetKantongList", uint(1), mock.AnythingOfType("*domain.KantongListRequest")).Return(expectedKantongs, expectedMeta, nil)

	req := httptest.NewRequest("GET", "/kantong", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestGetKantongByID_Success(t *testing.T) {
	app, mockUsecase := setupKantongController()

	expectedKantong := &domain.KantongResponse{
		ID:       "test-id-1",
		IDKartu:  "ABC123",
		Nama:     "Kantong Utama",
		Kategori: "Pengeluaran",
		Saldo:    100000,
		Warna:    "Navy",
	}

	mockUsecase.On("GetKantongByID", "test-id-1", uint(1)).Return(expectedKantong, nil)

	req := httptest.NewRequest("GET", "/kantong/test-id-1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestGetKantongByID_NotFound(t *testing.T) {
	app, mockUsecase := setupKantongController()

	mockUsecase.On("GetKantongByID", "non-existent", uint(1)).Return((*domain.KantongResponse)(nil), fmt.Errorf("kantong tidak ditemukan"))

	req := httptest.NewRequest("GET", "/kantong/non-existent", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestGetKantongByID_InvalidUUID(t *testing.T) {
	app, mockUsecase := setupKantongController()

	mockUsecase.On("GetKantongByID", "invalid-uuid", uint(1)).Return((*domain.KantongResponse)(nil), fmt.Errorf("kantong tidak ditemukan"))

	req := httptest.NewRequest("GET", "/kantong/invalid-uuid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestCreateKantong_Success(t *testing.T) {
	app, mockUsecase := setupKantongController()

	deskripsi := "Deskripsi kantong baru"
	saldo := 50000.0

	kantongRequest := domain.CreateKantongRequest{
		Nama:      "Kantong Baru",
		Kategori:  "Pengeluaran",
		Deskripsi: &deskripsi,
		Saldo:     &saldo,
		Warna:     "Navy",
	}

	expectedKantong := &domain.KantongResponse{
		ID:       "new-id",
		IDKartu:  "XYZ456",
		Nama:     "Kantong Baru",
		Kategori: "Pengeluaran",
		Saldo:    50000,
		Warna:    "Navy",
	}

	mockUsecase.On("CreateKantong", mock.AnythingOfType("*domain.CreateKantongRequest"), uint(1)).Return(expectedKantong, nil)

	body, _ := json.Marshal(kantongRequest)
	req := httptest.NewRequest("POST", "/kantong", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 201, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestCreateKantong_DuplicateName(t *testing.T) {
	app, mockUsecase := setupKantongController()

	deskripsi := "Deskripsi kantong"
	saldo := 50000.0

	kantongRequest := domain.CreateKantongRequest{
		Nama:      "Kantong Existing",
		Kategori:  "Pengeluaran",
		Deskripsi: &deskripsi,
		Saldo:     &saldo,
		Warna:     "Navy",
	}

	mockUsecase.On("CreateKantong", mock.AnythingOfType("*domain.CreateKantongRequest"), uint(1)).Return((*domain.KantongResponse)(nil), fmt.Errorf("nama kantong sudah ada"))

	body, _ := json.Marshal(kantongRequest)
	req := httptest.NewRequest("POST", "/kantong", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req)

	assert.Equal(t, 409, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}
