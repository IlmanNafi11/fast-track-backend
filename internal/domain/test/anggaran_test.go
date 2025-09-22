package domain_test

import (
	"fiber-boiler-plate/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnggaranItem_Structure(t *testing.T) {
	now := time.Now()
	anggaran := &domain.AnggaranItem{
		KantongID:   "550e8400-e29b-41d4-a716-446655440001",
		NamaKantong: "Kantong Test",
		Rencana:     &[]float64{1000000}[0],
		CarryIn:     150000,
		Penyesuaian: -50000,
		Terpakai:    450000,
		Sisa:        650000,
		Progres:     45.0,
		Bulan:       9,
		Tahun:       2024,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	assert.NotNil(t, anggaran)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", anggaran.KantongID)
	assert.Equal(t, "Kantong Test", anggaran.NamaKantong)
	assert.Equal(t, float64(1000000), *anggaran.Rencana)
	assert.Equal(t, float64(150000), anggaran.CarryIn)
	assert.Equal(t, float64(-50000), anggaran.Penyesuaian)
	assert.Equal(t, float64(450000), anggaran.Terpakai)
	assert.Equal(t, float64(650000), anggaran.Sisa)
	assert.Equal(t, float64(45.0), anggaran.Progres)
	assert.Equal(t, 9, anggaran.Bulan)
	assert.Equal(t, 2024, anggaran.Tahun)
}

func TestAnggaranListRequest_DefaultValues(t *testing.T) {
	req := domain.NewAnggaranListRequest()

	assert.NotNil(t, req)
	assert.Equal(t, "nama_kantong", req.SortBy)
	assert.Equal(t, "asc", req.SortDirection)
	assert.Equal(t, 1, req.Page)
	assert.Equal(t, 10, req.PerPage)
	assert.NotNil(t, req.Bulan)
	assert.NotNil(t, req.Tahun)
	assert.Equal(t, int(time.Now().Month()), *req.Bulan)
	assert.Equal(t, time.Now().Year(), *req.Tahun)
}

func TestPenyesuaianAnggaranRequest_Structure(t *testing.T) {
	req := &domain.PenyesuaianAnggaranRequest{
		KantongID: "550e8400-e29b-41d4-a716-446655440001",
		Jenis:     "kurangi",
		Jumlah:    50000,
		Bulan:     9,
		Tahun:     2024,
	}

	assert.NotNil(t, req)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", req.KantongID)
	assert.Equal(t, "kurangi", req.Jenis)
	assert.Equal(t, float64(50000), req.Jumlah)
	assert.Equal(t, 9, req.Bulan)
	assert.Equal(t, 2024, req.Tahun)
}

func TestToAnggaranResponse_Conversion(t *testing.T) {
	now := time.Now()
	anggaran := &domain.AnggaranItem{
		KantongID:   "550e8400-e29b-41d4-a716-446655440001",
		NamaKantong: "Kantong Test",
		Rencana:     &[]float64{1000000}[0],
		CarryIn:     150000,
		Penyesuaian: -50000,
		Terpakai:    450000,
		Sisa:        650000,
		Progres:     45.0,
		Bulan:       9,
		Tahun:       2024,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	response := domain.ToAnggaranResponse(anggaran)

	assert.NotNil(t, response)
	assert.Equal(t, anggaran.KantongID, response.KantongID)
	assert.Equal(t, anggaran.NamaKantong, response.NamaKantong)
	assert.Equal(t, anggaran.Rencana, response.Rencana)
	assert.Equal(t, anggaran.CarryIn, response.CarryIn)
	assert.Equal(t, anggaran.Penyesuaian, response.Penyesuaian)
	assert.Equal(t, anggaran.Terpakai, response.Terpakai)
	assert.Equal(t, anggaran.Sisa, response.Sisa)
	assert.Equal(t, anggaran.Progres, response.Progres)
	assert.Equal(t, anggaran.Bulan, response.Bulan)
	assert.Equal(t, anggaran.Tahun, response.Tahun)
}

func TestToAnggaranResponse_NilInput(t *testing.T) {
	response := domain.ToAnggaranResponse(nil)
	assert.Nil(t, response)
}

func TestToAnggaranDetailResponse_Conversion(t *testing.T) {
	now := time.Now()
	kantong := &domain.Kantong{
		ID:       "550e8400-e29b-41d4-a716-446655440001",
		IDKartu:  "K4N7G1",
		Nama:     "Kantong Test",
		Kategori: "Pengeluaran",
		Saldo:    750000,
		Warna:    "Navy",
	}

	statistik := []domain.StatistikHarian{
		{
			Tanggal:           now,
			JumlahTransaksi:   3,
			TotalPengeluaran:  125000,
			AkumulasiTerpakai: 125000,
		},
	}

	anggaran := &domain.AnggaranItem{
		KantongID:      "550e8400-e29b-41d4-a716-446655440001",
		NamaKantong:    "Kantong Test",
		Rencana:        &[]float64{1000000}[0],
		CarryIn:        150000,
		Penyesuaian:    -50000,
		Terpakai:       450000,
		Sisa:           650000,
		Progres:        45.0,
		DetailKantong:  kantong,
		StatistikBulan: statistik,
		Bulan:          9,
		Tahun:          2024,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	response := domain.ToAnggaranDetailResponse(anggaran)

	assert.NotNil(t, response)
	assert.Equal(t, anggaran.KantongID, response.KantongID)
	assert.Equal(t, anggaran.NamaKantong, response.NamaKantong)
	assert.NotNil(t, response.DetailKantong)
	assert.Equal(t, kantong.ID, response.DetailKantong.ID)
	assert.Equal(t, len(statistik), len(response.StatistikBulan))
	assert.Equal(t, statistik[0].JumlahTransaksi, response.StatistikBulan[0].JumlahTransaksi)
}

func TestToAnggaranResponseList_Conversion(t *testing.T) {
	now := time.Now()
	items := []*domain.AnggaranItem{
		{
			KantongID:   "550e8400-e29b-41d4-a716-446655440001",
			NamaKantong: "Kantong Test 1",
			Rencana:     &[]float64{1000000}[0],
			Bulan:       9,
			Tahun:       2024,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			KantongID:   "550e8400-e29b-41d4-a716-446655440002",
			NamaKantong: "Kantong Test 2",
			Rencana:     nil,
			Bulan:       9,
			Tahun:       2024,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	responses := domain.ToAnggaranResponseList(items)

	assert.NotNil(t, responses)
	assert.Equal(t, len(items), len(responses))
	assert.Equal(t, items[0].KantongID, responses[0].KantongID)
	assert.Equal(t, items[1].KantongID, responses[1].KantongID)
	assert.Equal(t, items[0].Rencana, responses[0].Rencana)
	assert.Nil(t, responses[1].Rencana)
}

func TestStatistikHarian_Structure(t *testing.T) {
	now := time.Now()
	statistik := &domain.StatistikHarian{
		Tanggal:           now,
		JumlahTransaksi:   5,
		TotalPengeluaran:  250000,
		AkumulasiTerpakai: 750000,
	}

	assert.NotNil(t, statistik)
	assert.Equal(t, now, statistik.Tanggal)
	assert.Equal(t, 5, statistik.JumlahTransaksi)
	assert.Equal(t, float64(250000), statistik.TotalPengeluaran)
	assert.Equal(t, float64(750000), statistik.AkumulasiTerpakai)
}
