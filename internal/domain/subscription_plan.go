package domain

import (
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionPlan struct {
	ID            uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Kode          string    `json:"kode" gorm:"uniqueIndex;not null"`
	Nama          string    `json:"nama" gorm:"not null"`
	Harga         float64   `json:"harga" gorm:"not null;default:0"`
	Interval      string    `json:"interval" gorm:"not null"`
	HariPercobaan int       `json:"hari_percobaan" gorm:"default:0"`
	Status        string    `json:"status" gorm:"not null;default:'aktif'"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (s *SubscriptionPlan) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	if s.Kode == "" {
		s.generateKode()
	}
	return nil
}

func (s *SubscriptionPlan) generateKode() {
	baseKode := generateBaseKode(s.Nama, s.Interval)
	randomSuffix := generateRandomSuffix()
	s.Kode = baseKode + "_" + randomSuffix
}

type CreateSubscriptionPlanRequest struct {
	Nama          string  `json:"nama" validate:"required,min=1,max=100"`
	Harga         float64 `json:"harga" validate:"gte=0"`
	Interval      string  `json:"interval" validate:"required,oneof=bulan tahun"`
	HariPercobaan int     `json:"hari_percobaan" validate:"gte=0"`
	Status        string  `json:"status" validate:"required,oneof=aktif 'non aktif'"`
}

type UpdateSubscriptionPlanRequest struct {
	Nama          string  `json:"nama" validate:"required,min=1,max=100"`
	Harga         float64 `json:"harga" validate:"required,gte=0"`
	Interval      string  `json:"interval" validate:"required,oneof=bulan tahun"`
	HariPercobaan int     `json:"hari_percobaan" validate:"required,gte=0"`
	Status        string  `json:"status" validate:"required,oneof=aktif 'non aktif'"`
}

type PatchSubscriptionPlanRequest struct {
	Nama          *string  `json:"nama" validate:"omitempty,min=1,max=100"`
	Harga         *float64 `json:"harga" validate:"omitempty,gte=0"`
	Interval      *string  `json:"interval" validate:"omitempty,oneof=bulan tahun"`
	HariPercobaan *int     `json:"hari_percobaan" validate:"omitempty,gte=0"`
	Status        *string  `json:"status" validate:"omitempty,oneof=aktif 'non aktif'"`
}

type SubscriptionPlanListRequest struct {
	Search        string `json:"search" query:"search"`
	SortBy        string `json:"sort_by" query:"sort_by" validate:"omitempty,oneof=nama harga"`
	SortDirection string `json:"sort_direction" query:"sort_direction" validate:"omitempty,oneof=asc desc"`
	Page          int    `json:"page" query:"page" validate:"omitempty,gte=1"`
	PerPage       int    `json:"per_page" query:"per_page" validate:"omitempty,gte=1,lte=100"`
}

func (req *SubscriptionPlanListRequest) SetDefaults() {
	if req.SortBy == "" {
		req.SortBy = "nama"
	}
	if req.SortDirection == "" {
		req.SortDirection = "asc"
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PerPage < 1 {
		req.PerPage = 10
	}
}

func (req *SubscriptionPlanListRequest) GetOffset() int {
	return (req.Page - 1) * req.PerPage
}

func generateBaseKode(nama, interval string) string {
	namaUpper := strings.ToUpper(nama)
	namaUpper = strings.ReplaceAll(namaUpper, " ", "_")

	intervalUpper := strings.ToUpper(interval)
	if intervalUpper == "BULAN" {
		intervalUpper = "MONTHLY"
	} else if intervalUpper == "TAHUN" {
		intervalUpper = "YEARLY"
	}

	return namaUpper + "_" + intervalUpper
}

func generateRandomSuffix() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	result := make([]byte, 3)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
