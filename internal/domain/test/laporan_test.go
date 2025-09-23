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

func TestStatistikKantongPeriode_Struct(t *testing.T) {
	periode := domain.PeriodeTanggal{
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-31",
	}

	dataKantong := []domain.DataKantongPeriode{
		{
			KantongID:        "550e8400-e29b-41d4-a716-446655440001",
			KantongNama:      "Kantong Belanja",
			TotalPengeluaran: 1500000.0,
		},
		{
			KantongID:        "550e8400-e29b-41d4-a716-446655440002",
			KantongNama:      "Transport",
			TotalPengeluaran: 800000.0,
		},
	}

	statistik := domain.StatistikKantongPeriode{
		Periode:          periode,
		DataKantong:      dataKantong,
		TotalPengeluaran: 2300000.0,
	}

	assert.Equal(t, periode, statistik.Periode)
	assert.Len(t, statistik.DataKantong, 2)
	assert.Equal(t, "Kantong Belanja", statistik.DataKantong[0].KantongNama)
	assert.Equal(t, 1500000.0, statistik.DataKantong[0].TotalPengeluaran)
	assert.Equal(t, 2300000.0, statistik.TotalPengeluaran)
}

func TestPengeluaranKantongDetail_Struct(t *testing.T) {
	periode := domain.PeriodeTanggal{
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-31",
	}

	dataKantong := []domain.DataKantongDetail{
		{
			KantongID:           "550e8400-e29b-41d4-a716-446655440001",
			KantongNama:         "Kantong Belanja",
			TotalPengeluaran:    1500000.0,
			PersentaseDariSaldo: 15.0,
			JumlahTransaksi:     25,
			RataRataPengeluaran: 60000.0,
			SaldoKantong:        10000000.0,
		},
	}

	detail := domain.PengeluaranKantongDetail{
		Periode:                periode,
		DataKantong:            dataKantong,
		TotalPengeluaran:       1500000.0,
		TotalSaldoSemuaKantong: 10000000.0,
	}

	assert.Equal(t, periode, detail.Periode)
	assert.Len(t, detail.DataKantong, 1)
	assert.Equal(t, "Kantong Belanja", detail.DataKantong[0].KantongNama)
	assert.Equal(t, 15.0, detail.DataKantong[0].PersentaseDariSaldo)
	assert.Equal(t, 60000.0, detail.DataKantong[0].RataRataPengeluaran)
	assert.Equal(t, 10000000.0, detail.TotalSaldoSemuaKantong)
}

func TestStatistikKantongPeriodeRequest_Struct(t *testing.T) {
	tanggalMulai := "2024-01-01"
	tanggalSelesai := "2024-01-31"

	req := domain.StatistikKantongPeriodeRequest{
		TanggalMulai:   &tanggalMulai,
		TanggalSelesai: &tanggalSelesai,
	}

	assert.Equal(t, "2024-01-01", *req.TanggalMulai)
	assert.Equal(t, "2024-01-31", *req.TanggalSelesai)
}

func TestPengeluaranKantongDetailRequest_Struct(t *testing.T) {
	tanggalMulai := "2024-01-01"
	tanggalSelesai := "2024-01-31"

	req := domain.PengeluaranKantongDetailRequest{
		TanggalMulai:   &tanggalMulai,
		TanggalSelesai: &tanggalSelesai,
	}

	assert.Equal(t, "2024-01-01", *req.TanggalMulai)
	assert.Equal(t, "2024-01-31", *req.TanggalSelesai)
}

func TestStatistikKantongPeriodeResponse_Struct(t *testing.T) {
	periode := domain.PeriodeTanggal{
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-31",
	}

	dataKantong := []domain.DataKantongPeriode{
		{
			KantongID:        "550e8400-e29b-41d4-a716-446655440001",
			KantongNama:      "Kantong Belanja",
			TotalPengeluaran: 1500000.0,
		},
	}

	statistik := domain.StatistikKantongPeriode{
		Periode:          periode,
		DataKantong:      dataKantong,
		TotalPengeluaran: 1500000.0,
	}

	response := domain.StatistikKantongPeriodeResponse{
		Success:   true,
		Message:   "Statistik kantong periode berhasil diambil",
		Code:      200,
		Data:      statistik,
		Timestamp: time.Now(),
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Statistik kantong periode berhasil diambil", response.Message)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, statistik, response.Data)
	assert.NotZero(t, response.Timestamp)
}

func TestPengeluaranKantongDetailResponse_Struct(t *testing.T) {
	periode := domain.PeriodeTanggal{
		TanggalMulai:   "2024-01-01",
		TanggalSelesai: "2024-01-31",
	}

	dataKantong := []domain.DataKantongDetail{
		{
			KantongID:           "550e8400-e29b-41d4-a716-446655440001",
			KantongNama:         "Kantong Belanja",
			TotalPengeluaran:    1500000.0,
			PersentaseDariSaldo: 15.0,
			JumlahTransaksi:     25,
			RataRataPengeluaran: 60000.0,
			SaldoKantong:        10000000.0,
		},
	}

	detail := domain.PengeluaranKantongDetail{
		Periode:                periode,
		DataKantong:            dataKantong,
		TotalPengeluaran:       1500000.0,
		TotalSaldoSemuaKantong: 10000000.0,
	}

	response := domain.PengeluaranKantongDetailResponse{
		Success:   true,
		Message:   "Detail pengeluaran kantong berhasil diambil",
		Code:      200,
		Data:      detail,
		Timestamp: time.Now(),
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Detail pengeluaran kantong berhasil diambil", response.Message)
	assert.Equal(t, 200, response.Code)
	assert.Equal(t, detail, response.Data)
	assert.NotZero(t, response.Timestamp)
}
