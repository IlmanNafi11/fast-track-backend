package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          string       `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nama        string       `json:"nama" gorm:"type:varchar(50);not null;uniqueIndex;index"`
	Deskripsi   *string      `json:"deskripsi" gorm:"type:varchar(255)"`
	Status      string       `json:"status" gorm:"type:varchar(20);not null;check:status IN ('aktif','non_aktif');default:'aktif'"`
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

type RolePermission struct {
	RoleID       string `gorm:"primaryKey"`
	PermissionID string `gorm:"primaryKey"`
	CreatedAt    time.Time
}

type CreateRoleRequest struct {
	Nama          string   `json:"nama" validate:"required,min=2,max=50,rolename"`
	Status        string   `json:"status" validate:"required,oneof=aktif non_aktif"`
	Deskripsi     *string  `json:"deskripsi" validate:"omitempty,max=255"`
	PermissionIDs []string `json:"permission_ids" validate:"required,min=1,dive,uuid"`
}

type UpdateRoleRequest struct {
	Nama          string   `json:"nama" validate:"required,min=2,max=50,rolename"`
	Status        string   `json:"status" validate:"required,oneof=aktif non_aktif"`
	Deskripsi     *string  `json:"deskripsi" validate:"omitempty,max=255"`
	PermissionIDs []string `json:"permission_ids" validate:"required,min=1,dive,uuid"`
}

type RoleListRequest struct {
	Search  *string `json:"search" query:"search"`
	Status  *string `json:"status" query:"status" validate:"omitempty,oneof=aktif non_aktif"`
	Page    int     `json:"page" query:"page" validate:"min=1"`
	PerPage int     `json:"per_page" query:"per_page" validate:"min=1,max=100"`
}

func NewRoleListRequest() *RoleListRequest {
	return &RoleListRequest{
		Page:    1,
		PerPage: 10,
	}
}

type RolePermissionListRequest struct {
	Page    int `json:"page" query:"page" validate:"min=1"`
	PerPage int `json:"per_page" query:"per_page" validate:"min=1,max=100"`
}

func NewRolePermissionListRequest() *RolePermissionListRequest {
	return &RolePermissionListRequest{
		Page:    1,
		PerPage: 20,
	}
}

type RoleResponse struct {
	ID          string               `json:"id"`
	Nama        string               `json:"nama"`
	Deskripsi   *string              `json:"deskripsi"`
	Status      string               `json:"status"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type RoleListItem struct {
	ID               string    `json:"id"`
	Nama             string    `json:"nama"`
	Deskripsi        *string   `json:"deskripsi"`
	Status           string    `json:"status"`
	PermissionsCount int       `json:"permissions_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type RoleListResponse struct {
	Success   bool           `json:"success"`
	Message   string         `json:"message"`
	Code      int            `json:"code"`
	Data      []RoleListItem `json:"data"`
	Meta      PaginationMeta `json:"meta"`
	Timestamp time.Time      `json:"timestamp"`
}

type RoleDetailResponse struct {
	Success   bool         `json:"success"`
	Message   string       `json:"message"`
	Code      int          `json:"code"`
	Data      RoleResponse `json:"data"`
	Timestamp time.Time    `json:"timestamp"`
}

func ToRoleResponse(role *Role) *RoleResponse {
	permissions := ToPermissionResponseList(convertPermissionSlice(role.Permissions))
	return &RoleResponse{
		ID:          role.ID,
		Nama:        role.Nama,
		Deskripsi:   role.Deskripsi,
		Status:      role.Status,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}
}

func ToRoleListItem(role *Role, permissionsCount int) *RoleListItem {
	return &RoleListItem{
		ID:               role.ID,
		Nama:             role.Nama,
		Deskripsi:        role.Deskripsi,
		Status:           role.Status,
		PermissionsCount: permissionsCount,
		CreatedAt:        role.CreatedAt,
		UpdatedAt:        role.UpdatedAt,
	}
}

func convertPermissionSlice(permissions []Permission) []*Permission {
	result := make([]*Permission, len(permissions))
	for i := range permissions {
		result[i] = &permissions[i]
	}
	return result
}
