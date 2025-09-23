package domain_test

import (
	"fiber-boiler-plate/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRoleDomain(t *testing.T) {
	t.Run("should create role with valid data", func(t *testing.T) {
		deskripsi := "Administrator sistem"
		role := &domain.Role{
			Nama:      "admin",
			Deskripsi: &deskripsi,
			Status:    "aktif",
		}

		assert.Equal(t, "admin", role.Nama)
		assert.Equal(t, "Administrator sistem", *role.Deskripsi)
		assert.Equal(t, "aktif", role.Status)
	})

	t.Run("should convert to role response correctly", func(t *testing.T) {
		deskripsi := "User biasa"
		permissions := []domain.Permission{
			{ID: "perm-1", Nama: "read", Kategori: "aplikasi"},
			{ID: "perm-2", Nama: "write", Kategori: "aplikasi"},
		}

		role := &domain.Role{
			ID:          "role-1",
			Nama:        "user",
			Deskripsi:   &deskripsi,
			Status:      "aktif",
			Permissions: permissions,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		response := domain.ToRoleResponse(role)

		assert.Equal(t, "role-1", response.ID)
		assert.Equal(t, "user", response.Nama)
		assert.Equal(t, "User biasa", *response.Deskripsi)
		assert.Equal(t, "aktif", response.Status)
		assert.Len(t, response.Permissions, 2)
		assert.Equal(t, "read", response.Permissions[0].Nama)
	})

	t.Run("should convert to role list item correctly", func(t *testing.T) {
		deskripsi := "Moderator"
		role := &domain.Role{
			ID:        "role-2",
			Nama:      "moderator",
			Deskripsi: &deskripsi,
			Status:    "non_aktif",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		listItem := domain.ToRoleListItem(role, 5)

		assert.Equal(t, "role-2", listItem.ID)
		assert.Equal(t, "moderator", listItem.Nama)
		assert.Equal(t, "Moderator", *listItem.Deskripsi)
		assert.Equal(t, "non_aktif", listItem.Status)
		assert.Equal(t, 5, listItem.PermissionsCount)
	})
}

func TestCreateRoleRequest(t *testing.T) {
	t.Run("should create valid request", func(t *testing.T) {
		deskripsi := "Test role"
		permissionIDs := []string{"perm-1", "perm-2"}

		req := &domain.CreateRoleRequest{
			Nama:          "test_role",
			Status:        "aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: permissionIDs,
		}

		assert.Equal(t, "test_role", req.Nama)
		assert.Equal(t, "aktif", req.Status)
		assert.Equal(t, "Test role", *req.Deskripsi)
		assert.Len(t, req.PermissionIDs, 2)
	})
}

func TestUpdateRoleRequest(t *testing.T) {
	t.Run("should create valid update request", func(t *testing.T) {
		deskripsi := "Updated role"
		permissionIDs := []string{"perm-1", "perm-2", "perm-3"}

		req := &domain.UpdateRoleRequest{
			Nama:          "updated_role",
			Status:        "non_aktif",
			Deskripsi:     &deskripsi,
			PermissionIDs: permissionIDs,
		}

		assert.Equal(t, "updated_role", req.Nama)
		assert.Equal(t, "non_aktif", req.Status)
		assert.Equal(t, "Updated role", *req.Deskripsi)
		assert.Len(t, req.PermissionIDs, 3)
	})
}

func TestRoleListRequest(t *testing.T) {
	t.Run("should create default request", func(t *testing.T) {
		req := domain.NewRoleListRequest()

		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 10, req.PerPage)
		assert.Nil(t, req.Search)
		assert.Nil(t, req.Status)
	})

	t.Run("should create request with parameters", func(t *testing.T) {
		search := "admin"
		status := "aktif"
		req := &domain.RoleListRequest{
			Search:  &search,
			Status:  &status,
			Page:    2,
			PerPage: 20,
		}

		assert.Equal(t, "admin", *req.Search)
		assert.Equal(t, "aktif", *req.Status)
		assert.Equal(t, 2, req.Page)
		assert.Equal(t, 20, req.PerPage)
	})
}

func TestRolePermissionListRequest(t *testing.T) {
	t.Run("should create default request", func(t *testing.T) {
		req := domain.NewRolePermissionListRequest()

		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 20, req.PerPage)
	})
}

func TestRolePermission(t *testing.T) {
	t.Run("should create role permission relation", func(t *testing.T) {
		rolePermission := &domain.RolePermission{
			RoleID:       "role-1",
			PermissionID: "perm-1",
			CreatedAt:    time.Now(),
		}

		assert.Equal(t, "role-1", rolePermission.RoleID)
		assert.Equal(t, "perm-1", rolePermission.PermissionID)
		assert.False(t, rolePermission.CreatedAt.IsZero())
	})
}
