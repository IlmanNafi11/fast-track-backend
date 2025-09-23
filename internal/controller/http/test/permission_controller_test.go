package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fiber-boiler-plate/internal/controller/http"
	"fiber-boiler-plate/internal/domain"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPermissionUsecase struct {
	mock.Mock
}

func (m *MockPermissionUsecase) GetPermissionList(req *domain.PermissionListRequest) ([]domain.PermissionResponse, *domain.PaginationMeta, error) {
	args := m.Called(req)
	return args.Get(0).([]domain.PermissionResponse), args.Get(1).(*domain.PaginationMeta), args.Error(2)
}

func (m *MockPermissionUsecase) GetPermissionByID(id string) (*domain.PermissionResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionUsecase) CreatePermission(req *domain.CreatePermissionRequest) (*domain.PermissionResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionUsecase) UpdatePermission(id string, req *domain.UpdatePermissionRequest) (*domain.PermissionResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PermissionResponse), args.Error(1)
}

func (m *MockPermissionUsecase) DeletePermission(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestApp() *fiber.App {
	app := fiber.New()
	return app
}

func TestPermissionController_GetPermissionList_Success(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	permissions := []domain.PermissionResponse{
		{
			ID:       "1",
			Nama:     "admin.dashboard.read",
			Kategori: "admin",
		},
	}
	meta := &domain.PaginationMeta{
		CurrentPage:  1,
		TotalPages:   1,
		TotalRecords: 1,
		PerPage:      10,
	}

	mockUsecase.On("GetPermissionList", mock.AnythingOfType("*domain.PermissionListRequest")).Return(permissions, meta, nil)

	app.Get("/permission", controller.GetPermissionList)

	req := httptest.NewRequest("GET", "/permission?page=1&per_page=10", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_GetPermissionList_WithSearch(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	permissions := []domain.PermissionResponse{}
	meta := &domain.PaginationMeta{
		CurrentPage:  1,
		TotalPages:   0,
		TotalRecords: 0,
		PerPage:      10,
	}

	mockUsecase.On("GetPermissionList", mock.AnythingOfType("*domain.PermissionListRequest")).Return(permissions, meta, nil)

	app.Get("/permission", controller.GetPermissionList)

	req := httptest.NewRequest("GET", "/permission?search=admin&kategori=admin", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_GetPermissionByID_Success(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	permission := &domain.PermissionResponse{
		ID:       "1",
		Nama:     "admin.dashboard.read",
		Kategori: "admin",
	}

	mockUsecase.On("GetPermissionByID", "1").Return(permission, nil)

	app.Get("/permission/:id", controller.GetPermissionByID)

	req := httptest.NewRequest("GET", "/permission/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_GetPermissionByID_NotFound(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	mockUsecase.On("GetPermissionByID", "1").Return(nil, errors.New("permission tidak ditemukan"))

	app.Get("/permission/:id", controller.GetPermissionByID)

	req := httptest.NewRequest("GET", "/permission/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_CreatePermission_Success(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	reqBody := domain.CreatePermissionRequest{
		Nama:     "admin.user.create",
		Kategori: "admin",
	}

	permission := &domain.PermissionResponse{
		ID:       "1",
		Nama:     "admin.user.create",
		Kategori: "admin",
	}

	mockUsecase.On("CreatePermission", &reqBody).Return(permission, nil)

	app.Post("/permission", controller.CreatePermission)

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/permission", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_CreatePermission_Conflict(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	reqBody := domain.CreatePermissionRequest{
		Nama:     "admin.user.create",
		Kategori: "admin",
	}

	mockUsecase.On("CreatePermission", &reqBody).Return(nil, errors.New("permission dengan nama tersebut sudah ada"))

	app.Post("/permission", controller.CreatePermission)

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/permission", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_UpdatePermission_Success(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	reqBody := domain.UpdatePermissionRequest{
		Nama:     "admin.dashboard.write",
		Kategori: "admin",
	}

	permission := &domain.PermissionResponse{
		ID:       "1",
		Nama:     "admin.dashboard.write",
		Kategori: "admin",
	}

	mockUsecase.On("UpdatePermission", "1", &reqBody).Return(permission, nil)

	app.Put("/permission/:id", controller.UpdatePermission)

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/permission/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_UpdatePermission_NotFound(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	reqBody := domain.UpdatePermissionRequest{
		Nama:     "admin.dashboard.write",
		Kategori: "admin",
	}

	mockUsecase.On("UpdatePermission", "1", &reqBody).Return(nil, errors.New("permission tidak ditemukan"))

	app.Put("/permission/:id", controller.UpdatePermission)

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("PUT", "/permission/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_DeletePermission_Success(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	mockUsecase.On("DeletePermission", "1").Return(nil)

	app.Delete("/permission/:id", controller.DeletePermission)

	req := httptest.NewRequest("DELETE", "/permission/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestPermissionController_DeletePermission_Conflict(t *testing.T) {
	mockUsecase := &MockPermissionUsecase{}
	controller := http.NewPermissionController(mockUsecase)
	app := setupTestApp()

	mockUsecase.On("DeletePermission", "1").Return(errors.New("permission tidak dapat dihapus karena masih digunakan oleh role lain"))

	app.Delete("/permission/:id", controller.DeletePermission)

	req := httptest.NewRequest("DELETE", "/permission/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusConflict, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}
