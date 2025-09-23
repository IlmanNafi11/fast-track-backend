package domain_test

import (
	"fiber-boiler-plate/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRingkasanLaporan_Struct(t *testing.T) {
	periode := domain.PeriodeTanggal{
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-31",
	}

	ringkasan := domain.RingkasanLaporan{
		TotalPemasukan:            5000000.0,
		TotalPengeluaran:          3500000.0,
		TotalSaldo:                12500000.0,
		RataRataPengeluaranHarian: 112903.23,
		Periode:                   periode,
	}

	assert.Equal(t, 5000000.0, ringkasan.TotalPemasukan)
	assert.Equal(t, 3500000.0, ringkasan.TotalPengeluaran)
	assert.Equal(t, 12500000.0, ringkasan.TotalSaldo)
	assert.Equal(t, 112903.23, ringkasan.RataRataPengeluaranHarian)
	assert.Equal(t, "2024-01-01", ringkasan.Periode.TanggalMulai)
	assert.Equal(t, "2024-01-31", ringkasan.Periode.TanggalSelesai)
}

func TestStatistikTahunan_Struct(t *testing.T) {
	dataBulanan := []domain.DataBulanan{
		{
			Bulan:            1,
			NamaBulan:        "Januari",
			TotalPemasukan:   5000000.0,
			TotalPengeluaran: 3500000.0,
		},
		{
			Bulan:            2,
			NamaBulan:        "Februari",
			TotalPemasukan:   5200000.0,
			TotalPengeluaran: 3800000.0,
		},
	}

	statistik := domain.StatistikTahunan{
		Tahun:                 2024,
		DataBulanan:           dataBulanan,
		TotalPemasukanTahun:   60000000.0,
		TotalPengeluaranTahun: 42000000.0,
	}

	assert.Equal(t, 2024, statistik.Tahun)
	assert.Len(t, statistik.DataBulanan, 2)
	assert.Equal(t, 1, statistik.DataBulanan[0].Bulan)
	assert.Equal(t, "Januari", statistik.DataBulanan[0].NamaBulan)
	assert.Equal(t, 60000000.0, statistik.TotalPemasukanTahun)
	assert.Equal(t, 42000000.0, statistik.TotalPengeluaranTahun)
}

func TestDataKantongBulanan_Struct(t *testing.T) {
	dataKantong := domain.DataKantongBulanan{
		KantongID:        "550e8400-e29b-41d4-a716-446655440001",
		KantongNama:      "Kantong Belanja",
		Kategori:         "Pengeluaran",
		TotalPengeluaran: 1500000.0,
		JumlahTransaksi:  25,
		Persentase:       42.86,
	}

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440001", dataKantong.KantongID)
	assert.Equal(t, "Kantong Belanja", dataKantong.KantongNama)
	assert.Equal(t, "Pengeluaran", dataKantong.Kategori)
	assert.Equal(t, 1500000.0, dataKantong.TotalPengeluaran)
	assert.Equal(t, 25, dataKantong.JumlahTransaksi)
	assert.Equal(t, 42.86, dataKantong.Persentase)
}

func TestRingkasanLaporanResponse_Struct(t *testing.T) {
	periode := domain.PeriodeTanggal{
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-31",
	}

	ringkasan := domain.RingkasanLaporan{
		TotalPemasukan:            5000000.0,
		TotalPengeluaran:          3500000.0,
		TotalSaldo:                12500000.0,
		RataRataPengeluaranHarian: 112903.23,
		Periode:                   periode,
	}

	response := domain.RingkasanLaporanResponse{
		Success:   true,
		Message:   "Ringkasan laporan berhasil diambil",
		Code:      200,
		Data:      ringkasan,
		Timestamp: time.Now(),
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Ringkasan laporan berhasil diambil", response.Message)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, ringkasan.TotalPemasukan, response.Data.TotalPemasukan)
	assert.NotZero(t, response.Timestamp)
}
