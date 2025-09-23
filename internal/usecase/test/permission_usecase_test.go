package usecase_test

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockPermissionRepository struct {
	mock.Mock
}

func (m *MockPermissionRepository) GetAll(req *domain.PermissionListRequest) ([]*domain.Permission, int, error) {
	args := m.Called(req)
	return args.Get(0).([]*domain.Permission), args.Int(1), args.Error(2)
}

func (m *MockPermissionRepository) GetByID(id string) (*domain.Permission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) GetByNama(nama string) (*domain.Permission, error) {
	args := m.Called(nama)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Permission), args.Error(1)
}

func (m *MockPermissionRepository) Create(permission *domain.Permission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) Update(permission *domain.Permission) error {
	args := m.Called(permission)
	return args.Error(0)
}

func (m *MockPermissionRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPermissionRepository) IsNameExists(nama string, excludeID ...string) (bool, error) {
	args := m.Called(nama)
	return args.Bool(0), args.Error(1)
}

func (m *MockPermissionRepository) IsUsedByRoles(id string) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func TestNewPermissionUsecase(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	assert.NotNil(t, uc)
}

func TestPermissionUsecase_GetPermissionList(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	permissions := []*domain.Permission{
		{
			ID:       "1",
			Nama:     "admin.dashboard.read",
			Kategori: "admin",
		},
		{
			ID:       "2",
			Nama:     "aplikasi.transaksi.create",
			Kategori: "aplikasi",
		},
	}

	req := &domain.PermissionListRequest{
		Page:    1,
		PerPage: 10,
	}

	mockRepo.On("GetAll", req).Return(permissions, 2, nil)

	result, meta, err := uc.GetPermissionList(req)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 1, meta.CurrentPage)
	assert.Equal(t, 2, meta.TotalRecords)
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_GetPermissionByID_Success(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	permission := &domain.Permission{
		ID:       "1",
		Nama:     "admin.dashboard.read",
		Kategori: "admin",
	}

	mockRepo.On("GetByID", "1").Return(permission, nil)

	result, err := uc.GetPermissionByID("1")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin.dashboard.read", result.Nama)
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_GetPermissionByID_NotFound(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	mockRepo.On("GetByID", "1").Return(nil, gorm.ErrRecordNotFound)

	result, err := uc.GetPermissionByID("1")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "permission tidak ditemukan", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_CreatePermission_Success(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	req := &domain.CreatePermissionRequest{
		Nama:     "admin.user.create",
		Kategori: "admin",
	}

	mockRepo.On("IsNameExists", "admin.user.create").Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*domain.Permission")).Return(nil)

	result, err := uc.CreatePermission(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin.user.create", result.Nama)
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_CreatePermission_NameExists(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	req := &domain.CreatePermissionRequest{
		Nama:     "admin.user.create",
		Kategori: "admin",
	}

	mockRepo.On("IsNameExists", "admin.user.create").Return(true, nil)

	result, err := uc.CreatePermission(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "permission dengan nama tersebut sudah ada", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_UpdatePermission_Success(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	permission := &domain.Permission{
		ID:        "1",
		Nama:      "admin.dashboard.read",
		Kategori:  "admin",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	req := &domain.UpdatePermissionRequest{
		Nama:     "admin.dashboard.write",
		Kategori: "admin",
	}

	mockRepo.On("GetByID", "1").Return(permission, nil)
	mockRepo.On("IsNameExists", "admin.dashboard.write").Return(false, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Permission")).Return(nil)

	result, err := uc.UpdatePermission("1", req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "admin.dashboard.write", result.Nama)
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_UpdatePermission_NotFound(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	req := &domain.UpdatePermissionRequest{
		Nama:     "admin.dashboard.write",
		Kategori: "admin",
	}

	mockRepo.On("GetByID", "1").Return(nil, gorm.ErrRecordNotFound)

	result, err := uc.UpdatePermission("1", req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "permission tidak ditemukan", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_DeletePermission_Success(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	permission := &domain.Permission{
		ID:       "1",
		Nama:     "admin.dashboard.read",
		Kategori: "admin",
	}

	mockRepo.On("GetByID", "1").Return(permission, nil)
	mockRepo.On("IsUsedByRoles", "1").Return(false, nil)
	mockRepo.On("Delete", "1").Return(nil)

	err := uc.DeletePermission("1")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_DeletePermission_UsedByRoles(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	permission := &domain.Permission{
		ID:       "1",
		Nama:     "admin.dashboard.read",
		Kategori: "admin",
	}

	mockRepo.On("GetByID", "1").Return(permission, nil)
	mockRepo.On("IsUsedByRoles", "1").Return(true, nil)

	err := uc.DeletePermission("1")

	assert.Error(t, err)
	assert.Equal(t, "permission tidak dapat dihapus karena masih digunakan oleh role lain", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestPermissionUsecase_DeletePermission_NotFound(t *testing.T) {
	mockRepo := &MockPermissionRepository{}
	uc := usecase.NewPermissionUsecase(mockRepo)

	mockRepo.On("GetByID", "1").Return(nil, gorm.ErrRecordNotFound)

	err := uc.DeletePermission("1")

	assert.Error(t, err)
	assert.Equal(t, "permission tidak ditemukan", err.Error())
	mockRepo.AssertExpectations(t)
}
