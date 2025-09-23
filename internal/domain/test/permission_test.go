package domain_test

import (
	"fiber-boiler-plate/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPermission_Creation(t *testing.T) {
	now := time.Now()
	permission := &domain.Permission{
		ID:        "550e8400-e29b-41d4-a716-446655440001",
		Nama:      "admin.dashboard.read",
		Kategori:  "admin",
		Deskripsi: &[]string{"Akses untuk membaca dashboard admin"}[0],
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.NotNil(t, permission)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", permission.ID)
	assert.Equal(t, "admin.dashboard.read", permission.Nama)
	assert.Equal(t, "admin", permission.Kategori)
	assert.Equal(t, "Akses untuk membaca dashboard admin", *permission.Deskripsi)
	assert.Equal(t, now, permission.CreatedAt)
	assert.Equal(t, now, permission.UpdatedAt)
}

func TestCreatePermissionRequest_Structure(t *testing.T) {
	req := &domain.CreatePermissionRequest{
		Nama:      "aplikasi.transaksi.create",
		Kategori:  "aplikasi",
		Deskripsi: &[]string{"Permission untuk membuat transaksi"}[0],
	}

	assert.NotNil(t, req)
	assert.Equal(t, "aplikasi.transaksi.create", req.Nama)
	assert.Equal(t, "aplikasi", req.Kategori)
	assert.Equal(t, "Permission untuk membuat transaksi", *req.Deskripsi)
}

func TestUpdatePermissionRequest_Structure(t *testing.T) {
	req := &domain.UpdatePermissionRequest{
		Nama:      "admin.user.manage",
		Kategori:  "admin",
		Deskripsi: &[]string{"Permission untuk mengelola pengguna"}[0],
	}

	assert.NotNil(t, req)
	assert.Equal(t, "admin.user.manage", req.Nama)
	assert.Equal(t, "admin", req.Kategori)
	assert.Equal(t, "Permission untuk mengelola pengguna", *req.Deskripsi)
}

func TestPermissionListRequest_Defaults(t *testing.T) {
	req := domain.NewPermissionListRequest()

	assert.NotNil(t, req)
	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.PerPage)
	assert.Nil(t, req.Search)
	assert.Nil(t, req.Kategori)
}

func TestPermissionResponse_Structure(t *testing.T) {
	now := time.Now()
	response := &domain.PermissionResponse{
		ID:        "550e8400-e29b-41d4-a716-446655440001",
		Nama:      "admin.dashboard.read",
		Kategori:  "admin",
		Deskripsi: &[]string{"Akses untuk membaca dashboard admin"}[0],
		CreatedAt: now,
		UpdatedAt: now,
	}

	assert.NotNil(t, response)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", response.ID)
	assert.Equal(t, "admin.dashboard.read", response.Nama)
	assert.Equal(t, "admin", response.Kategori)
	assert.Equal(t, "Akses untuk membaca dashboard admin", *response.Deskripsi)
	assert.Equal(t, now, response.CreatedAt)
	assert.Equal(t, now, response.UpdatedAt)
}

func TestToPermissionResponse(t *testing.T) {
	now := time.Now()
	permission := &domain.Permission{
		ID:        "550e8400-e29b-41d4-a716-446655440001",
		Nama:      "admin.dashboard.read",
		Kategori:  "admin",
		Deskripsi: &[]string{"Akses untuk membaca dashboard admin"}[0],
		CreatedAt: now,
		UpdatedAt: now,
	}

	response := domain.ToPermissionResponse(permission)

	assert.NotNil(t, response)
	assert.Equal(t, permission.ID, response.ID)
	assert.Equal(t, permission.Nama, response.Nama)
	assert.Equal(t, permission.Kategori, response.Kategori)
	assert.Equal(t, permission.Deskripsi, response.Deskripsi)
	assert.Equal(t, permission.CreatedAt, response.CreatedAt)
	assert.Equal(t, permission.UpdatedAt, response.UpdatedAt)
}

func TestToPermissionResponseList(t *testing.T) {
	now := time.Now()
	permissions := []*domain.Permission{
		{
			ID:        "550e8400-e29b-41d4-a716-446655440001",
			Nama:      "admin.dashboard.read",
			Kategori:  "admin",
			Deskripsi: &[]string{"Akses untuk membaca dashboard admin"}[0],
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        "550e8400-e29b-41d4-a716-446655440002",
			Nama:      "aplikasi.transaksi.create",
			Kategori:  "aplikasi",
			Deskripsi: nil,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	responses := domain.ToPermissionResponseList(permissions)

	assert.NotNil(t, responses)
	assert.Len(t, responses, 2)
	assert.Equal(t, permissions[0].ID, responses[0].ID)
	assert.Equal(t, permissions[1].ID, responses[1].ID)
}

func TestPermissionListResponse_Structure(t *testing.T) {
	now := time.Now()
	response := &domain.PermissionListResponse{
		Success: true,
		Message: "Daftar permission berhasil diambil",
		Code:    200,
		Data: []domain.PermissionResponse{
			{
				ID:        "550e8400-e29b-41d4-a716-446655440001",
				Nama:      "admin.dashboard.read",
				Kategori:  "admin",
				Deskripsi: &[]string{"Akses untuk membaca dashboard admin"}[0],
				CreatedAt: now,
				UpdatedAt: now,
			},
		},
		Meta: domain.PaginationMeta{
			CurrentPage:  1,
			TotalPages:   1,
			TotalRecords: 1,
			PerPage:      10,
		},
		Timestamp: now,
	}

	assert.NotNil(t, response)
	assert.True(t, response.Success)
	assert.Equal(t, "Daftar permission berhasil diambil", response.Message)
	assert.Equal(t, 200, response.Code)
	assert.Len(t, response.Data, 1)
	assert.Equal(t, 1, response.Meta.CurrentPage)
	assert.Equal(t, now, response.Timestamp)
}

func TestPermissionDetailResponse_Structure(t *testing.T) {
	now := time.Now()
	response := &domain.PermissionDetailResponse{
		Success: true,
		Message: "Detail permission berhasil diambil",
		Code:    200,
		Data: domain.PermissionResponse{
			ID:        "550e8400-e29b-41d4-a716-446655440001",
			Nama:      "admin.dashboard.read",
			Kategori:  "admin",
			Deskripsi: &[]string{"Akses untuk membaca dashboard admin"}[0],
			CreatedAt: now,
			UpdatedAt: now,
		},
		Timestamp: now,
	}

	assert.NotNil(t, response)
	assert.True(t, response.Success)
	assert.Equal(t, "Detail permission berhasil diambil", response.Message)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", response.Data.ID)
	assert.Equal(t, now, response.Timestamp)
}
