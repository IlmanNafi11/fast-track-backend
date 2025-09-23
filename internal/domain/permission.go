package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Permission struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Nama      string    `json:"nama" gorm:"type:varchar(100);not null;uniqueIndex;index"`
	Kategori  string    `json:"kategori" gorm:"type:varchar(20);not null;check:kategori IN ('admin','aplikasi')"`
	Deskripsi *string   `json:"deskripsi" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Permission) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

type CreatePermissionRequest struct {
	Nama      string  `json:"nama" validate:"required,min=3,max=100,permissionname"`
	Kategori  string  `json:"kategori" validate:"required,oneof=admin aplikasi"`
	Deskripsi *string `json:"deskripsi" validate:"omitempty,max=255"`
}

type UpdatePermissionRequest struct {
	Nama      string  `json:"nama" validate:"required,min=3,max=100,permissionname"`
	Kategori  string  `json:"kategori" validate:"required,oneof=admin aplikasi"`
	Deskripsi *string `json:"deskripsi" validate:"omitempty,max=255"`
}

type PermissionListRequest struct {
	Search   *string `json:"search" query:"search"`
	Kategori *string `json:"kategori" query:"kategori" validate:"omitempty,oneof=admin aplikasi"`
	Page     int     `json:"page" query:"page" validate:"min=1"`
	PerPage  int     `json:"per_page" query:"per_page" validate:"min=1,max=100"`
}

func NewPermissionListRequest() *PermissionListRequest {
	return &PermissionListRequest{
		Page:    1,
		PerPage: 10,
	}
}

type PermissionResponse struct {
	ID        string    `json:"id"`
	Nama      string    `json:"nama"`
	Kategori  string    `json:"kategori"`
	Deskripsi *string   `json:"deskripsi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PermissionListResponse struct {
	Success   bool                 `json:"success"`
	Message   string               `json:"message"`
	Code      int                  `json:"code"`
	Data      []PermissionResponse `json:"data"`
	Meta      PaginationMeta       `json:"meta"`
	Timestamp time.Time            `json:"timestamp"`
}

type PermissionDetailResponse struct {
	Success   bool               `json:"success"`
	Message   string             `json:"message"`
	Code      int                `json:"code"`
	Data      PermissionResponse `json:"data"`
	Timestamp time.Time          `json:"timestamp"`
}

func ToPermissionResponse(permission *Permission) *PermissionResponse {
	return &PermissionResponse{
		ID:        permission.ID,
		Nama:      permission.Nama,
		Kategori:  permission.Kategori,
		Deskripsi: permission.Deskripsi,
		CreatedAt: permission.CreatedAt,
		UpdatedAt: permission.UpdatedAt,
	}
}

func ToPermissionResponseList(permissions []*Permission) []PermissionResponse {
	var responses []PermissionResponse
	for _, permission := range permissions {
		responses = append(responses, *ToPermissionResponse(permission))
	}
	return responses
}
