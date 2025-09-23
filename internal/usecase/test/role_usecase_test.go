package usecase_test

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) GetAll(req *domain.RoleListRequest) ([]*domain.Role, int, error) {
	args := m.Called(req)
	return args.Get(0).([]*domain.Role), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) GetByID(id string) (*domain.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) GetByNama(nama string) (*domain.Role, error) {
	args := m.Called(nama)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Role), args.Error(1)
}

func (m *MockRoleRepository) Create(role *domain.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) CreateWithPermissions(role *domain.Role, permissionIDs []string) error {
	args := m.Called(role, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) Update(role *domain.Role) error {
	args := m.Called(role)
	return args.Error(0)
}

func (m *MockRoleRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRoleRepository) IsNameExists(nama string, excludeID ...string) (bool, error) {
	args := m.Called(nama, excludeID)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) IsUsedByUsers(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockRoleRepository) GetRolePermissions(roleID string, req *domain.RolePermissionListRequest) ([]*domain.Permission, int, error) {
	args := m.Called(roleID, req)
	return args.Get(0).([]*domain.Permission), args.Int(1), args.Error(2)
}

func (m *MockRoleRepository) UpdateRolePermissions(roleID string, permissionIDs []string) error {
	args := m.Called(roleID, permissionIDs)
	return args.Error(0)
}

func (m *MockRoleRepository) ValidatePermissions(permissionIDs []string) ([]string, error) {
	args := m.Called(permissionIDs)
	return args.Get(0).([]string), args.Error(1)
}

func TestRoleUsecase_GetRoleList(t *testing.T) {
	t.Run("should return role list successfully", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		req := &domain.RoleListRequest{Page: 1, PerPage: 10}

		deskripsi := "Admin role"
		roles := []*domain.Role{
			{
				ID:        "role-1",
				Nama:      "admin",
				Deskripsi: &deskripsi,
				Status:    "aktif",
				Permissions: []domain.Permission{
					{ID: "perm-1", Nama: "read"},
					{ID: "perm-2", Nama: "write"},
				},
			},
		}

		mockRoleRepo.On("GetAll", req).Return(roles, 1, nil)

		result, meta, err := roleUsecase.GetRoleList(req)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "admin", result[0].Nama)
		assert.Equal(t, 2, result[0].PermissionsCount)
		assert.Equal(t, 1, meta.CurrentPage)
		assert.Equal(t, 1, meta.TotalRecords)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("should handle repository error", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		req := &domain.RoleListRequest{Page: 1, PerPage: 10}
		mockRoleRepo.On("GetAll", req).Return([]*domain.Role{}, 0, errors.New("database error"))

		result, meta, err := roleUsecase.GetRoleList(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Nil(t, meta)
		mockRoleRepo.AssertExpectations(t)
	})
}

func TestRoleUsecase_GetRoleByID(t *testing.T) {
	t.Run("should return role by id successfully", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		deskripsi := "Admin role"
		role := &domain.Role{
			ID:        "role-1",
			Nama:      "admin",
			Deskripsi: &deskripsi,
			Status:    "aktif",
			Permissions: []domain.Permission{
				{ID: "perm-1", Nama: "read"},
			},
		}

		mockRoleRepo.On("GetByID", "role-1").Return(role, nil)

		result, err := roleUsecase.GetRoleByID("role-1")

		assert.NoError(t, err)
		assert.Equal(t, "role-1", result.ID)
		assert.Equal(t, "admin", result.Nama)
		assert.Len(t, result.Permissions, 1)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("should return error when role not found", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		mockRoleRepo.On("GetByID", "role-1").Return((*domain.Role)(nil), gorm.ErrRecordNotFound)

		result, err := roleUsecase.GetRoleByID("role-1")

		assert.Error(t, err)
		assert.Equal(t, "role tidak ditemukan", err.Error())
		assert.Nil(t, result)
		mockRoleRepo.AssertExpectations(t)
	})
}

func TestRoleUsecase_CreateRole(t *testing.T) {
	t.Run("should create role successfully", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		deskripsi := "New role"
		req := &domain.CreateRoleRequest{
			Nama:          "new_role",
			Status:        "aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"perm-1", "perm-2"},
		}

		mockRoleRepo.On("IsNameExists", "new_role", mock.Anything).Return(false, nil)
		mockRoleRepo.On("ValidatePermissions", []string{"perm-1", "perm-2"}).Return([]string{}, nil)
		mockRoleRepo.On("CreateWithPermissions", mock.AnythingOfType("*domain.Role"), []string{"perm-1", "perm-2"}).Return(nil)

		createdRole := &domain.Role{
			ID:        "role-1",
			Nama:      "new_role",
			Deskripsi: &deskripsi,
			Status:    "aktif",
			Permissions: []domain.Permission{
				{ID: "perm-1", Nama: "read"},
				{ID: "perm-2", Nama: "write"},
			},
		}
		mockRoleRepo.On("GetByID", mock.AnythingOfType("string")).Return(createdRole, nil)

		result, err := roleUsecase.CreateRole(req)

		assert.NoError(t, err)
		assert.Equal(t, "new_role", result.Nama)
		assert.Equal(t, "aktif", result.Status)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("should return error when role name already exists", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		deskripsi := "Existing role"
		req := &domain.CreateRoleRequest{
			Nama:          "existing_role",
			Status:        "aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"perm-1"},
		}

		mockRoleRepo.On("IsNameExists", "existing_role", mock.Anything).Return(true, nil)

		result, err := roleUsecase.CreateRole(req)

		assert.Error(t, err)
		assert.Equal(t, "role dengan nama tersebut sudah ada", err.Error())
		assert.Nil(t, result)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("should return error when permissions are invalid", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		deskripsi := "Invalid role"
		req := &domain.CreateRoleRequest{
			Nama:          "invalid_role",
			Status:        "aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"invalid-perm"},
		}

		mockRoleRepo.On("IsNameExists", "invalid_role", mock.Anything).Return(false, nil)
		mockRoleRepo.On("ValidatePermissions", []string{"invalid-perm"}).Return([]string{"invalid-perm"}, nil)

		result, err := roleUsecase.CreateRole(req)

		assert.Error(t, err)
		assert.Equal(t, "beberapa permission tidak ditemukan atau tidak valid", err.Error())
		assert.Nil(t, result)
		mockRoleRepo.AssertExpectations(t)
	})
}

func TestRoleUsecase_UpdateRole(t *testing.T) {
	t.Run("should update role successfully", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		existingRole := &domain.Role{
			ID:     "role-1",
			Nama:   "old_name",
			Status: "aktif",
		}

		deskripsi := "Updated role"
		req := &domain.UpdateRoleRequest{
			Nama:          "new_name",
			Status:        "non_aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: []string{"perm-1"},
		}

		mockRoleRepo.On("GetByID", "role-1").Return(existingRole, nil)
		mockRoleRepo.On("IsNameExists", "new_name", mock.Anything).Return(false, nil)
		mockRoleRepo.On("ValidatePermissions", []string{"perm-1"}).Return([]string{}, nil)
		mockRoleRepo.On("Update", mock.AnythingOfType("*domain.Role")).Return(nil)
		mockRoleRepo.On("UpdateRolePermissions", "role-1", []string{"perm-1"}).Return(nil)

		updatedRole := &domain.Role{
			ID:        "role-1",
			Nama:      "new_name",
			Deskripsi: &deskripsi,
			Status:    "non_aktif",
		}
		mockRoleRepo.On("GetByID", "role-1").Return(updatedRole, nil)

		result, err := roleUsecase.UpdateRole("role-1", req)

		assert.NoError(t, err)
		assert.Equal(t, "new_name", result.Nama)
		assert.Equal(t, "non_aktif", result.Status)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("should return error when role not found", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		req := &domain.UpdateRoleRequest{
			Nama:          "new_name",
			Status:        "aktif",
			PermissionIDs: []string{"perm-1"},
		}

		mockRoleRepo.On("GetByID", "role-1").Return((*domain.Role)(nil), gorm.ErrRecordNotFound)

		result, err := roleUsecase.UpdateRole("role-1", req)

		assert.Error(t, err)
		assert.Equal(t, "role tidak ditemukan", err.Error())
		assert.Nil(t, result)
		mockRoleRepo.AssertExpectations(t)
	})
}

func TestRoleUsecase_DeleteRole(t *testing.T) {
	t.Run("should delete role successfully", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		existingRole := &domain.Role{
			ID:   "role-1",
			Nama: "test_role",
		}

		mockRoleRepo.On("GetByID", "role-1").Return(existingRole, nil)
		mockRoleRepo.On("IsUsedByUsers", "role-1").Return(false, nil)
		mockRoleRepo.On("Delete", "role-1").Return(nil)

		err := roleUsecase.DeleteRole("role-1")

		assert.NoError(t, err)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("should return error when role is used by users", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		existingRole := &domain.Role{
			ID:   "role-1",
			Nama: "test_role",
		}

		mockRoleRepo.On("GetByID", "role-1").Return(existingRole, nil)
		mockRoleRepo.On("IsUsedByUsers", "role-1").Return(true, nil)

		err := roleUsecase.DeleteRole("role-1")

		assert.Error(t, err)
		assert.Equal(t, "role tidak dapat dihapus karena masih digunakan oleh pengguna lain", err.Error())
		mockRoleRepo.AssertExpectations(t)
	})
}

func TestRoleUsecase_GetRolePermissions(t *testing.T) {
	t.Run("should return role permissions successfully", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		role := &domain.Role{ID: "role-1"}
		req := &domain.RolePermissionListRequest{Page: 1, PerPage: 20}

		permissions := []*domain.Permission{
			{ID: "perm-1", Nama: "read", Kategori: "aplikasi"},
			{ID: "perm-2", Nama: "write", Kategori: "aplikasi"},
		}

		mockRoleRepo.On("GetByID", "role-1").Return(role, nil)
		mockRoleRepo.On("GetRolePermissions", "role-1", req).Return(permissions, 2, nil)

		result, meta, err := roleUsecase.GetRolePermissions("role-1", req)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "read", result[0].Nama)
		assert.Equal(t, 1, meta.CurrentPage)
		assert.Equal(t, 2, meta.TotalRecords)
		mockRoleRepo.AssertExpectations(t)
	})

	t.Run("should return error when role not found", func(t *testing.T) {
		mockRoleRepo := new(MockRoleRepository)
		mockPermissionRepo := new(MockPermissionRepository)
		roleUsecase := usecase.NewRoleUsecase(mockRoleRepo, mockPermissionRepo)

		req := &domain.RolePermissionListRequest{Page: 1, PerPage: 20}

		mockRoleRepo.On("GetByID", "role-1").Return((*domain.Role)(nil), gorm.ErrRecordNotFound)

		result, meta, err := roleUsecase.GetRolePermissions("role-1", req)

		assert.Error(t, err)
		assert.Equal(t, "role tidak ditemukan", err.Error())
		assert.Nil(t, result)
		assert.Nil(t, meta)
		mockRoleRepo.AssertExpectations(t)
	})
}
