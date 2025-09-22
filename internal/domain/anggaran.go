package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Anggaran struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	KantongID   string    `json:"kantong_id" gorm:"type:uuid;not null;index:idx_anggarans_kantong_user"`
	UserID      uint      `json:"-" gorm:"not null;index:idx_anggarans_kantong_user,idx_anggarans_user_bulan_tahun"`
	Bulan       int       `json:"bulan" gorm:"not null;index:idx_anggarans_user_bulan_tahun;check:bulan >= 1 AND bulan <= 12"`
	Tahun       int       `json:"tahun" gorm:"not null;index:idx_anggarans_user_bulan_tahun;check:tahun >= 2020"`
	Rencana     *float64  `json:"rencana" gorm:"type:decimal(15,2)"`
	CarryIn     float64   `json:"carry_in" gorm:"type:decimal(15,2);not null;default:0"`
	Penyesuaian float64   `json:"penyesuaian" gorm:"type:decimal(15,2);not null;default:0"`
	Terpakai    float64   `json:"terpakai" gorm:"type:decimal(15,2);not null;default:0"`
	Sisa        float64   `json:"sisa" gorm:"type:decimal(15,2);not null;default:0"`
	Progres     float64   `json:"progres" gorm:"type:decimal(5,2);not null;default:0"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `json:"-" gorm:"foreignKey:UserID"`
	Kantong     Kantong   `json:"-" gorm:"foreignKey:KantongID"`
}

func (a *Anggaran) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}

func (a *Anggaran) TableName() string {
	return "anggarans"
}

type PenyesuaianAnggaran struct {
	ID         string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	AnggaranID string    `json:"anggaran_id" gorm:"type:uuid;not null;index"`
	Jenis      string    `json:"jenis" gorm:"type:varchar(10);not null;check:jenis IN ('tambah','kurangi')"`
	Jumlah     float64   `json:"jumlah" gorm:"type:decimal(15,2);not null;check:jumlah >= 0"`
	CreatedAt  time.Time `json:"created_at"`
	Anggaran   Anggaran  `json:"-" gorm:"foreignKey:AnggaranID"`
}

func (p *PenyesuaianAnggaran) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

func (p *PenyesuaianAnggaran) TableName() string {
	return "penyesuaian_anggarans"
}

type AnggaranItem struct {
	KantongID      string            `json:"kantong_id"`
	NamaKantong    string            `json:"nama_kantong"`
	Rencana        *float64          `json:"rencana"`
	CarryIn        float64           `json:"carry_in"`
	Penyesuaian    float64           `json:"penyesuaian"`
	Terpakai       float64           `json:"terpakai"`
	Sisa           float64           `json:"sisa"`
	Progres        float64           `json:"progres"`
	DetailKantong  *Kantong          `json:"detail_kantong,omitempty"`
	StatistikBulan []StatistikHarian `json:"statistik_bulan,omitempty"`
	Bulan          int               `json:"bulan"`
	Tahun          int               `json:"tahun"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
}

type StatistikHarian struct {
	Tanggal           time.Time `json:"tanggal"`
	JumlahTransaksi   int       `json:"jumlah_transaksi"`
	TotalPengeluaran  float64   `json:"total_pengeluaran"`
	AkumulasiTerpakai float64   `json:"akumulasi_terpakai"`
}

type AnggaranListRequest struct {
	Search        *string `json:"search" query:"search"`
	SortBy        string  `json:"sort_by" query:"sort_by" validate:"omitempty,oneof=nama_kantong rencana carry_in penyesuaian sisa"`
	SortDirection string  `json:"sort_direction" query:"sort_direction" validate:"omitempty,oneof=asc desc"`
	Page          int     `json:"page" query:"page" validate:"min=1"`
	PerPage       int     `json:"per_page" query:"per_page" validate:"min=1,max=100"`
	Bulan         *int    `json:"bulan" query:"bulan" validate:"omitempty,min=1,max=12"`
	Tahun         *int    `json:"tahun" query:"tahun" validate:"omitempty,min=2020"`
}

func NewAnggaranListRequest() *AnggaranListRequest {
	now := time.Now()
	return &AnggaranListRequest{
		SortBy:        "nama_kantong",
		SortDirection: "asc",
		Page:          1,
		PerPage:       10,
		Bulan:         &[]int{int(now.Month())}[0],
		Tahun:         &[]int{now.Year()}[0],
	}
}

type PenyesuaianAnggaranRequest struct {
	KantongID string  `json:"kantong_id" validate:"required,uuid"`
	Jenis     string  `json:"jenis" validate:"required,oneof=tambah kurangi"`
	Jumlah    float64 `json:"jumlah" validate:"required,min=0"`
	Bulan     int     `json:"bulan" validate:"required,min=1,max=12"`
	Tahun     int     `json:"tahun" validate:"required,min=2020"`
}

type AnggaranResponse struct {
	KantongID   string   `json:"kantong_id"`
	NamaKantong string   `json:"nama_kantong"`
	Rencana     *float64 `json:"rencana"`
	CarryIn     float64  `json:"carry_in"`
	Penyesuaian float64  `json:"penyesuaian"`
	Terpakai    float64  `json:"terpakai"`
	Sisa        float64  `json:"sisa"`
	Progres     float64  `json:"progres"`
	Bulan       int      `json:"bulan"`
	Tahun       int      `json:"tahun"`
}

type AnggaranDetailResponse struct {
	KantongID      string            `json:"kantong_id"`
	NamaKantong    string            `json:"nama_kantong"`
	DetailKantong  *KantongResponse  `json:"detail_kantong"`
	Rencana        *float64          `json:"rencana"`
	CarryIn        float64           `json:"carry_in"`
	Penyesuaian    float64           `json:"penyesuaian"`
	Terpakai       float64           `json:"terpakai"`
	Sisa           float64           `json:"sisa"`
	Progres        float64           `json:"progres"`
	StatistikBulan []StatistikHarian `json:"statistik_bulan"`
	Bulan          int               `json:"bulan"`
	Tahun          int               `json:"tahun"`
}

func ToAnggaranResponse(item *AnggaranItem) *AnggaranResponse {
	if item == nil {
		return nil
	}
	return &AnggaranResponse{
		KantongID:   item.KantongID,
		NamaKantong: item.NamaKantong,
		Rencana:     item.Rencana,
		CarryIn:     item.CarryIn,
		Penyesuaian: item.Penyesuaian,
		Terpakai:    item.Terpakai,
		Sisa:        item.Sisa,
		Progres:     item.Progres,
		Bulan:       item.Bulan,
		Tahun:       item.Tahun,
	}
}

func ToAnggaranDetailResponse(item *AnggaranItem) *AnggaranDetailResponse {
	if item == nil {
		return nil
	}

	var detailKantong *KantongResponse
	if item.DetailKantong != nil {
		detailKantong = ToKantongResponse(item.DetailKantong)
	}

	return &AnggaranDetailResponse{
		KantongID:      item.KantongID,
		NamaKantong:    item.NamaKantong,
		DetailKantong:  detailKantong,
		Rencana:        item.Rencana,
		CarryIn:        item.CarryIn,
		Penyesuaian:    item.Penyesuaian,
		Terpakai:       item.Terpakai,
		Sisa:           item.Sisa,
		Progres:        item.Progres,
		StatistikBulan: item.StatistikBulan,
		Bulan:          item.Bulan,
		Tahun:          item.Tahun,
	}
}

func ToAnggaranResponseList(items []*AnggaranItem) []*AnggaranResponse {
	responses := make([]*AnggaranResponse, len(items))
	for i, item := range items {
		responses[i] = ToAnggaranResponse(item)
	}
	return responses
}
