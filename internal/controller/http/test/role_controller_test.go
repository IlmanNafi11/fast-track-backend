package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fiber-boiler-plate/internal/controller/http"
	"fiber-boiler-plate/internal/domain"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRoleUsecase struct {
	mock.Mock
}

func (m *MockRoleUsecase) GetRoleList(req *domain.RoleListRequest) ([]domain.RoleListItem, *domain.PaginationMeta, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]domain.RoleListItem), args.Get(1).(*domain.PaginationMeta), args.Error(2)
}

func (m *MockRoleUsecase) GetRoleByID(id string) (*domain.RoleResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RoleResponse), args.Error(1)
}

func (m *MockRoleUsecase) CreateRole(req *domain.CreateRoleRequest) (*domain.RoleResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RoleResponse), args.Error(1)
}

func (m *MockRoleUsecase) UpdateRole(id string, req *domain.UpdateRoleRequest) (*domain.RoleResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RoleResponse), args.Error(1)
}

func (m *MockRoleUsecase) DeleteRole(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoleUsecase) GetRolePermissions(roleID string, req *domain.RolePermissionListRequest) ([]domain.PermissionResponse, *domain.PaginationMeta, error) {
	args := m.Called(roleID, req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]domain.PermissionResponse), args.Get(1).(*domain.PaginationMeta), args.Error(2)
}

func TestRoleController_GetRoleList(t *testing.T) {
	t.Run("should return role list successfully", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Get("/roles", controller.GetRoleList)

		roles := []domain.RoleListItem{
			{
				ID:               "role-1",
				Nama:             "admin",
				Status:           "aktif",
				PermissionsCount: 5,
			},
		}
		meta := &domain.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   1,
			TotalRecords: 1,
			PerPage:      10,
		}

		mockUsecase.On("GetRoleList", mock.AnythingOfType("*domain.RoleListRequest")).Return(roles, meta, nil)

		req := httptest.NewRequest("GET", "/roles?page=1&per_page=10", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		json.Unmarshal(body, &response)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Daftar role berhasil diambil", response["message"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should handle usecase error", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Get("/roles", controller.GetRoleList)

		mockUsecase.On("GetRoleList", mock.AnythingOfType("*domain.RoleListRequest")).Return(nil, nil, errors.New("database error"))

		req := httptest.NewRequest("GET", "/roles", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestRoleController_GetRoleByID(t *testing.T) {
	t.Run("should return role by id successfully", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Get("/roles/:id", controller.GetRoleByID)

		deskripsi := "Admin role"
		role := &domain.RoleResponse{
			ID:        "role-1",
			Nama:      "admin",
			Deskripsi: &deskripsi,
			Status:    "aktif",
		}

		mockUsecase.On("GetRoleByID", "role-1").Return(role, nil)

		req := httptest.NewRequest("GET", "/roles/role-1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		json.Unmarshal(body, &response)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Detail role berhasil diambil", response["message"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 404 when role not found", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Get("/roles/:id", controller.GetRoleByID)

		mockUsecase.On("GetRoleByID", "role-1").Return(nil, errors.New("role tidak ditemukan"))

		req := httptest.NewRequest("GET", "/roles/role-1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestRoleController_CreateRole(t *testing.T) {
	t.Run("should create role successfully", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Post("/roles", controller.CreateRole)

		deskripsi := "New role"
		requestBody := domain.CreateRoleRequest{
			Nama:          "new_role",
			Status:        "aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"550e8400-e29b-41d4-a716-446655440001"},
		}

		createdRole := &domain.RoleResponse{
			ID:        "role-1",
			Nama:      "new_role",
			Deskripsi: &deskripsi,
			Status:    "aktif",
		}

		mockUsecase.On("CreateRole", mock.AnythingOfType("*domain.CreateRoleRequest")).Return(createdRole, nil)

		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/roles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		json.Unmarshal(body, &response)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Role berhasil dibuat", response["message"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 409 when role name already exists", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Post("/roles", controller.CreateRole)

		deskripsi := "Existing role"
		requestBody := domain.CreateRoleRequest{
			Nama:          "existing_role",
			Status:        "aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"550e8400-e29b-41d4-a716-446655440001"},
		}

		mockUsecase.On("CreateRole", mock.AnythingOfType("*domain.CreateRoleRequest")).Return(nil, errors.New("role dengan nama tersebut sudah ada"))

		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/roles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 409, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 422 when permissions are invalid", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Post("/roles", controller.CreateRole)

		deskripsi := "Invalid role"
		requestBody := domain.CreateRoleRequest{
			Nama:          "invalid_role",
			Status:        "aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"550e8400-e29b-41d4-a716-446655440001"},
		}

		mockUsecase.On("CreateRole", mock.AnythingOfType("*domain.CreateRoleRequest")).Return(nil, errors.New("beberapa permission tidak ditemukan atau tidak valid"))

		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("POST", "/roles", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 422, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestRoleController_UpdateRole(t *testing.T) {
	t.Run("should update role successfully", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Put("/roles/:id", controller.UpdateRole)

		deskripsi := "Updated role"
		requestBody := domain.UpdateRoleRequest{
			Nama:          "updated_role",
			Status:        "non_aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"550e8400-e29b-41d4-a716-446655440001", "550e8400-e29b-41d4-a716-446655440002"},
		}

		updatedRole := &domain.RoleResponse{
			ID:        "role-1",
			Nama:      "updated_role",
			Deskripsi: &deskripsi,
			Status:    "non_aktif",
		}

		mockUsecase.On("UpdateRole", "role-1", mock.AnythingOfType("*domain.UpdateRoleRequest")).Return(updatedRole, nil)

		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("PUT", "/roles/role-1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		json.Unmarshal(body, &response)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Role berhasil diupdate", response["message"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 404 when role not found", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Put("/roles/:id", controller.UpdateRole)

		deskripsi := "Not found role"
		requestBody := domain.UpdateRoleRequest{
			Nama:          "not_found_role",
			Status:        "aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"550e8400-e29b-41d4-a716-446655440001"},
		}

		mockUsecase.On("UpdateRole", "role-1", mock.AnythingOfType("*domain.UpdateRoleRequest")).Return(nil, errors.New("role tidak ditemukan"))

		jsonBody, _ := json.Marshal(requestBody)
		req := httptest.NewRequest("PUT", "/roles/role-1", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestRoleController_DeleteRole(t *testing.T) {
	t.Run("should delete role successfully", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Delete("/roles/:id", controller.DeleteRole)

		mockUsecase.On("DeleteRole", "role-1").Return(nil)

		req := httptest.NewRequest("DELETE", "/roles/role-1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		json.Unmarshal(body, &response)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Role berhasil dihapus", response["message"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 409 when role is used by users", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Delete("/roles/:id", controller.DeleteRole)

		mockUsecase.On("DeleteRole", "role-1").Return(errors.New("role tidak dapat dihapus karena masih digunakan oleh pengguna lain"))

		req := httptest.NewRequest("DELETE", "/roles/role-1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 409, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}

func TestRoleController_GetRolePermissions(t *testing.T) {
	t.Run("should return role permissions successfully", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Get("/roles/:id/permissions", controller.GetRolePermissions)

		permissions := []domain.PermissionResponse{
			{
				ID:       "perm-1",
				Nama:     "read",
				Kategori: "aplikasi",
			},
		}
		meta := &domain.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   1,
			TotalRecords: 1,
			PerPage:      20,
		}

		mockUsecase.On("GetRolePermissions", "role-1", mock.AnythingOfType("*domain.RolePermissionListRequest")).Return(permissions, meta, nil)

		req := httptest.NewRequest("GET", "/roles/role-1/permissions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		body, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		json.Unmarshal(body, &response)

		assert.True(t, response["success"].(bool))
		assert.Equal(t, "Daftar permissions role berhasil diambil", response["message"])
		mockUsecase.AssertExpectations(t)
	})

	t.Run("should return 404 when role not found", func(t *testing.T) {
		mockUsecase := new(MockRoleUsecase)
		controller := http.NewRoleController(mockUsecase)

		app := fiber.New()
		app.Get("/roles/:id/permissions", controller.GetRolePermissions)

		mockUsecase.On("GetRolePermissions", "role-1", mock.AnythingOfType("*domain.RolePermissionListRequest")).Return(nil, nil, errors.New("role tidak ditemukan"))

		req := httptest.NewRequest("GET", "/roles/role-1/permissions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}
