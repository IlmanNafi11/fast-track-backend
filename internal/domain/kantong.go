package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Kantong struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IDKartu   string    `json:"id_kartu" gorm:"type:varchar(6);uniqueIndex;not null"`
	UserID    uint      `json:"-" gorm:"not null;index"`
	Nama      string    `json:"nama" gorm:"type:varchar(100);not null;index"`
	Kategori  string    `json:"kategori" gorm:"type:varchar(20);not null;check:kategori IN ('Pengeluaran','Tabungan','Darurat','Transport','Tidak Spesifik')"`
	Deskripsi *string   `json:"deskripsi" gorm:"type:varchar(500)"`
	Limit     *float64  `json:"limit" gorm:"column:limit_amount;type:decimal(15,2);check:limit_amount >= 0"`
	Saldo     float64   `json:"saldo" gorm:"type:decimal(15,2);not null;default:0;check:saldo >= 0"`
	Warna     string    `json:"warna" gorm:"type:varchar(10);not null;check:warna IN ('Navy','Glass','Purple','Green','Red')"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"-" gorm:"foreignKey:UserID"`
}

type KantongResponse struct {
	ID        string    `json:"id"`
	IDKartu   string    `json:"id_kartu"`
	Nama      string    `json:"nama"`
	Kategori  string    `json:"kategori"`
	Deskripsi *string   `json:"deskripsi"`
	Limit     *float64  `json:"limit"`
	Saldo     float64   `json:"saldo"`
	Warna     string    `json:"warna"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func ToKantongResponse(kantong *Kantong) *KantongResponse {
	if kantong == nil {
		return nil
	}
	return &KantongResponse{
		ID:        kantong.ID,
		IDKartu:   kantong.IDKartu,
		Nama:      kantong.Nama,
		Kategori:  kantong.Kategori,
		Deskripsi: kantong.Deskripsi,
		Limit:     kantong.Limit,
		Saldo:     kantong.Saldo,
		Warna:     kantong.Warna,
		CreatedAt: kantong.CreatedAt,
		UpdatedAt: kantong.UpdatedAt,
	}
}

func ToKantongResponseList(kantongs []*Kantong) []*KantongResponse {
	responses := make([]*KantongResponse, len(kantongs))
	for i, kantong := range kantongs {
		responses[i] = ToKantongResponse(kantong)
	}
	return responses
}

func (k *Kantong) BeforeCreate(tx *gorm.DB) error {
	if k.ID == "" {
		k.ID = uuid.New().String()
	}
	return nil
}

type CreateKantongRequest struct {
	Nama      string   `json:"nama" validate:"required,min=1,max=100"`
	Kategori  string   `json:"kategori" validate:"required,oneof=Pengeluaran Tabungan Darurat Transport 'Tidak Spesifik'"`
	Deskripsi *string  `json:"deskripsi" validate:"omitempty,max=500"`
	Limit     *float64 `json:"limit" validate:"omitempty,min=0"`
	Saldo     *float64 `json:"saldo" validate:"omitempty,min=0"`
	Warna     string   `json:"warna" validate:"required,oneof=Navy Glass Purple Green Red"`
}

type UpdateKantongRequest struct {
	Nama      string   `json:"nama" validate:"required,min=1,max=100"`
	Kategori  string   `json:"kategori" validate:"required,oneof=Pengeluaran Tabungan Darurat Transport 'Tidak Spesifik'"`
	Deskripsi *string  `json:"deskripsi" validate:"omitempty,max=500"`
	Limit     *float64 `json:"limit" validate:"omitempty,min=0"`
	Saldo     float64  `json:"saldo" validate:"min=0"`
	Warna     string   `json:"warna" validate:"required,oneof=Navy Glass Purple Green Red"`
}

type PatchKantongRequest struct {
	Nama      *string  `json:"nama" validate:"omitempty,min=1,max=100"`
	Kategori  *string  `json:"kategori" validate:"omitempty,oneof=Pengeluaran Tabungan Darurat Transport 'Tidak Spesifik'"`
	Deskripsi *string  `json:"deskripsi" validate:"omitempty,max=500"`
	Limit     *float64 `json:"limit" validate:"omitempty,min=0"`
	Saldo     *float64 `json:"saldo" validate:"omitempty,min=0"`
	Warna     *string  `json:"warna" validate:"omitempty,oneof=Navy Glass Purple Green Red"`
}

type KantongListRequest struct {
	Search        *string `json:"search" query:"search"`
	SortBy        string  `json:"sort_by" query:"sort_by" validate:"omitempty,oneof=nama saldo"`
	SortDirection string  `json:"sort_direction" query:"sort_direction" validate:"omitempty,oneof=asc desc"`
	Page          int     `json:"page" query:"page" validate:"min=1"`
	PerPage       int     `json:"per_page" query:"per_page" validate:"min=1,max=100"`
}

func NewKantongListRequest() *KantongListRequest {
	return &KantongListRequest{
		SortBy:        "nama",
		SortDirection: "asc",
		Page:          1,
		PerPage:       10,
	}
}

type TransferKantongRequest struct {
	KantongAsalID   string  `json:"kantong_asal_id" validate:"required,uuid"`
	KantongTujuanID string  `json:"kantong_tujuan_id" validate:"required,uuid"`
	Jumlah          float64 `json:"jumlah" validate:"required,gt=0"`
	Catatan         *string `json:"catatan" validate:"omitempty,max=500"`
}

type TransferKantongDetail struct {
	ID           string  `json:"id"`
	Nama         string  `json:"nama"`
	SaldoSebelum float64 `json:"saldo_sebelum"`
	SaldoSesudah float64 `json:"saldo_sesudah"`
}

type TransferResult struct {
	TransferID      string                `json:"transfer_id"`
	KantongAsal     TransferKantongDetail `json:"kantong_asal"`
	KantongTujuan   TransferKantongDetail `json:"kantong_tujuan"`
	Jumlah          float64               `json:"jumlah"`
	Catatan         *string               `json:"catatan"`
	TanggalTransfer time.Time             `json:"tanggal_transfer"`
}

type TransferKantongResponse struct {
	Success   bool            `json:"success"`
	Message   string          `json:"message"`
	Code      int             `json:"code"`
	Data      *TransferResult `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
}
