package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaksi struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uint      `json:"-" gorm:"not null;index"`
	KantongID string    `json:"kantong_id" gorm:"type:uuid;not null;index"`
	Tanggal   time.Time `json:"tanggal" gorm:"type:date;not null;index"`
	Jenis     string    `json:"jenis" gorm:"type:varchar(20);not null;check:jenis IN ('Pemasukan','Pengeluaran')"`
	Jumlah    float64   `json:"jumlah" gorm:"type:decimal(15,2);not null;check:jumlah > 0"`
	Catatan   *string   `json:"catatan" gorm:"type:varchar(500)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
	Kantong   Kantong   `json:"-" gorm:"foreignKey:KantongID"`
}

func (t *Transaksi) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

type TransaksiResponse struct {
	ID          string    `json:"id"`
	Tanggal     string    `json:"tanggal"`
	Jenis       string    `json:"jenis"`
	Jumlah      float64   `json:"jumlah"`
	KantongID   string    `json:"kantong_id"`
	KantongNama string    `json:"kantong_nama"`
	Catatan     *string   `json:"catatan"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateTransaksiRequest struct {
	KantongID string  `json:"kantong_id" validate:"required,uuid"`
	Tanggal   string  `json:"tanggal" validate:"required"`
	Jenis     string  `json:"jenis" validate:"required,oneof=Pemasukan Pengeluaran"`
	Jumlah    float64 `json:"jumlah" validate:"required,gt=0"`
	Catatan   *string `json:"catatan" validate:"omitempty,max=500"`
}

type UpdateTransaksiRequest struct {
	KantongID string  `json:"kantong_id" validate:"required,uuid"`
	Tanggal   string  `json:"tanggal" validate:"required"`
	Jenis     string  `json:"jenis" validate:"required,oneof=Pemasukan Pengeluaran"`
	Jumlah    float64 `json:"jumlah" validate:"required,gt=0"`
	Catatan   *string `json:"catatan" validate:"omitempty,max=500"`
}

type PatchTransaksiRequest struct {
	KantongID *string  `json:"kantong_id,omitempty" validate:"omitempty,uuid"`
	Tanggal   *string  `json:"tanggal,omitempty"`
	Jenis     *string  `json:"jenis,omitempty" validate:"omitempty,oneof=Pemasukan Pengeluaran"`
	Jumlah    *float64 `json:"jumlah,omitempty" validate:"omitempty,gt=0"`
	Catatan   *string  `json:"catatan,omitempty" validate:"omitempty,max=500"`
}

type TransaksiListRequest struct {
	Search         *string `json:"search" query:"search"`
	Jenis          *string `json:"jenis" query:"jenis" validate:"omitempty,oneof=Pemasukan Pengeluaran"`
	KantongNama    *string `json:"kantong_nama" query:"kantong_nama"`
	TanggalMulai   *string `json:"tanggal_mulai" query:"tanggal_mulai"`
	TanggalSelesai *string `json:"tanggal_selesai" query:"tanggal_selesai"`
	SortBy         string  `json:"sort_by" query:"sort_by" validate:"oneof=tanggal jumlah" default:"tanggal"`
	SortDirection  string  `json:"sort_direction" query:"sort_direction" validate:"oneof=asc desc" default:"desc"`
	Page           int     `json:"page" query:"page" validate:"min=1" default:"1"`
	PerPage        int     `json:"per_page" query:"per_page" validate:"min=1,max=100" default:"10"`
}

type TransaksiListResponse struct {
	Success   bool                `json:"success"`
	Message   string              `json:"message"`
	Code      int                 `json:"code"`
	Data      []TransaksiResponse `json:"data"`
	Meta      PaginationMeta      `json:"meta"`
	Timestamp time.Time           `json:"timestamp"`
}

type TransaksiDetailResponse struct {
	Success   bool              `json:"success"`
	Message   string            `json:"message"`
	Code      int               `json:"code"`
	Data      TransaksiResponse `json:"data"`
	Timestamp time.Time         `json:"timestamp"`
}
