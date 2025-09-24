package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Invoice struct {
	ID                 string            `json:"invoice_id" gorm:"primaryKey;type:varchar(50)"`
	UserID             uint              `json:"user_id" gorm:"not null;index"`
	User               User              `json:"user" gorm:"foreignKey:UserID"`
	Jumlah             float64           `json:"jumlah" gorm:"not null;check:jumlah > 0"`
	Status             string            `json:"status" gorm:"not null;default:'pending';check:status IN ('sukses','gagal','pending')"`
	DibayarPada        *time.Time        `json:"dibayar_pada" gorm:"index"`
	MetodePembayaran   *string           `json:"metode_pembayaran"`
	Keterangan         *string           `json:"keterangan" gorm:"type:text"`
	SubscriptionPlanID *uuid.UUID        `json:"subscription_plan_id" gorm:"type:uuid;index"`
	SubscriptionPlan   *SubscriptionPlan `json:"subscription_plan,omitempty" gorm:"foreignKey:SubscriptionPlanID"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

func (i *Invoice) BeforeCreate(tx *gorm.DB) error {
	if i.ID == "" {
		i.ID = generateInvoiceID()
	}
	return nil
}

func generateInvoiceID() string {
	now := time.Now()
	timestamp := now.Format("200601021504")
	return "INV-" + timestamp + "-" + uuid.New().String()[:8]
}

type InvoiceListRequest struct {
	Search         *string `json:"search" query:"search"`
	Status         *string `json:"status" query:"status" validate:"omitempty,oneof=sukses gagal pending"`
	NamaUser       *string `json:"nama_user" query:"nama_user"`
	TanggalMulai   *string `json:"tanggal_mulai" query:"tanggal_mulai" validate:"omitempty,datetime=2006-01-02"`
	TanggalSelesai *string `json:"tanggal_selesai" query:"tanggal_selesai" validate:"omitempty,datetime=2006-01-02"`
	SortBy         string  `json:"sort_by" query:"sort_by" validate:"omitempty,oneof=invoice_id nama_user jumlah status dibayar_pada"`
	SortDirection  string  `json:"sort_direction" query:"sort_direction" validate:"omitempty,oneof=asc desc"`
	Page           int     `json:"page" query:"page" validate:"min=1"`
	PerPage        int     `json:"per_page" query:"per_page" validate:"min=1,max=100"`
}

func (r *InvoiceListRequest) GetOffset() int {
	return (r.Page - 1) * r.PerPage
}

func (req *InvoiceListRequest) SetDefaults() {
	if req.SortBy == "" {
		req.SortBy = "dibayar_pada"
	}
	if req.SortDirection == "" {
		req.SortDirection = "desc"
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PerPage < 1 {
		req.PerPage = 10
	}
}

type UpdateInvoiceStatusRequest struct {
	Status     string  `json:"status" validate:"required,oneof=sukses gagal pending"`
	Keterangan *string `json:"keterangan" validate:"omitempty,max=500"`
}

type InvoiceResponse struct {
	InvoiceID            string     `json:"invoice_id"`
	NamaUser             string     `json:"nama_user"`
	UserID               uint       `json:"user_id"`
	UserEmail            string     `json:"user_email"`
	Jumlah               float64    `json:"jumlah"`
	Status               string     `json:"status"`
	DibayarPada          *time.Time `json:"dibayar_pada"`
	MetodePembayaran     *string    `json:"metode_pembayaran"`
	Keterangan           *string    `json:"keterangan"`
	SubscriptionPlanID   *uuid.UUID `json:"subscription_plan_id"`
	SubscriptionPlanNama *string    `json:"subscription_plan_nama"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type InvoiceDetailResponse struct {
	InvoiceID            string     `json:"invoice_id"`
	NamaUser             string     `json:"nama_user"`
	UserID               uint       `json:"user_id"`
	UserEmail            string     `json:"user_email"`
	Jumlah               float64    `json:"jumlah"`
	Status               string     `json:"status"`
	DibayarPada          *time.Time `json:"dibayar_pada"`
	MetodePembayaran     *string    `json:"metode_pembayaran"`
	Keterangan           *string    `json:"keterangan"`
	SubscriptionPlanID   *uuid.UUID `json:"subscription_plan_id"`
	SubscriptionPlanNama *string    `json:"subscription_plan_nama"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type InvoiceStatistics struct {
	TotalInvoice         int64                 `json:"total_invoice"`
	TotalSukses          int64                 `json:"total_sukses"`
	TotalGagal           int64                 `json:"total_gagal"`
	TotalPending         int64                 `json:"total_pending"`
	TotalPendapatan      float64               `json:"total_pendapatan"`
	RataRataPembayaran   float64               `json:"rata_rata_pembayaran"`
	InvoiceBulanan       []InvoiceStatsBulanan `json:"invoice_bulanan"`
	TopSubscriptionPlans []TopSubscriptionPlan `json:"top_subscription_plans"`
}

type InvoiceStatsBulanan struct {
	Bulan           string  `json:"bulan"`
	TotalInvoice    int64   `json:"total_invoice"`
	TotalSukses     int64   `json:"total_sukses"`
	TotalPendapatan float64 `json:"total_pendapatan"`
}

type TopSubscriptionPlan struct {
	SubscriptionPlanNama string  `json:"subscription_plan_nama"`
	JumlahInvoice        int64   `json:"jumlah_invoice"`
	TotalPendapatan      float64 `json:"total_pendapatan"`
}

type InvoiceStatisticsRequest struct {
	Bulan *int `json:"bulan" query:"bulan" validate:"omitempty,min=1,max=12"`
	Tahun *int `json:"tahun" query:"tahun" validate:"omitempty,min=2020"`
}
