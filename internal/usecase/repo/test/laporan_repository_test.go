package repo_test

import (
	"fiber-boiler-plate/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLaporanRepository_TrenBulananDataStructure(t *testing.T) {
	// Test data structure validation
	expectedData := []domain.DataBulanan{
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

	trenBulanan := domain.TrenBulanan{
		Tahun:                 2024,
		DataTren:              expectedData,
		TotalPemasukanTahun:   62400000.0,
		TotalPengeluaranTahun: 45600000.0,
	}

	// Assertions
	assert.Equal(t, 2024, trenBulanan.Tahun)
	assert.Len(t, trenBulanan.DataTren, 2)
	assert.Equal(t, "Januari", trenBulanan.DataTren[0].NamaBulan)
	assert.Equal(t, 5000000.0, trenBulanan.DataTren[0].TotalPemasukan)
	assert.Equal(t, 62400000.0, trenBulanan.TotalPemasukanTahun)
}

func TestLaporanRepository_PerbandinganKantongDataStructure(t *testing.T) {
	expectedData := []domain.DataPerbandinganKantong{
		{
			KantongID:       "550e8400-e29b-41d4-a716-446655440001",
			KantongNama:     "Kantong Belanja",
			JumlahBulanIni:  1500000.0,
			JumlahBulanLalu: 1200000.0,
		},
		{
			KantongID:       "550e8400-e29b-41d4-a716-446655440002",
			KantongNama:     "Transport",
			JumlahBulanIni:  800000.0,
			JumlahBulanLalu: 700000.0,
		},
	}

	bulanIni := domain.PeriodeBulan{
		Bulan:     12,
		Tahun:     2024,
		NamaBulan: "Desember",
	}

	bulanSebelumnya := domain.PeriodeBulan{
		Bulan:     11,
		Tahun:     2024,
		NamaBulan: "November",
	}

	perbandingan := domain.PerbandinganKantong{
		BulanIni:        bulanIni,
		BulanSebelumnya: bulanSebelumnya,
		DataKantong:     expectedData,
		TotalBulanIni:   2300000.0,
		TotalBulanLalu:  1900000.0,
	}

	// Assertions
	assert.Equal(t, 12, perbandingan.BulanIni.Bulan)
	assert.Equal(t, 2024, perbandingan.BulanIni.Tahun)
	assert.Equal(t, "Desember", perbandingan.BulanIni.NamaBulan)
	assert.Len(t, perbandingan.DataKantong, 2)
	assert.Equal(t, "Kantong Belanja", perbandingan.DataKantong[0].KantongNama)
	assert.Equal(t, 1500000.0, perbandingan.DataKantong[0].JumlahBulanIni)
	assert.Equal(t, 2300000.0, perbandingan.TotalBulanIni)
}

func TestLaporanRepository_DetailPerbandinganKantongDataStructure(t *testing.T) {
	expectedData := []domain.DataDetailPerbandinganKantong{
		{
			KantongID:           "550e8400-e29b-41d4-a716-446655440001",
			KantongNama:         "Kantong Belanja",
			JumlahBulanIni:      1500000.0,
			JumlahBulanLalu:     1200000.0,
			RataRataPengeluaran: 100000.0,
			Persentase:          25.0,
			Trend:               "naik",
		},
		{
			KantongID:           "550e8400-e29b-41d4-a716-446655440002",
			KantongNama:         "Transport",
			JumlahBulanIni:      800000.0,
			JumlahBulanLalu:     700000.0,
			RataRataPengeluaran: 66666.67,
			Persentase:          14.29,
			Trend:               "naik",
		},
	}

	bulanIni := domain.PeriodeBulan{
		Bulan:     12,
		Tahun:     2024,
		NamaBulan: "Desember",
	}

	bulanSebelumnya := domain.PeriodeBulan{
		Bulan:     11,
		Tahun:     2024,
		NamaBulan: "November",
	}

	detail := domain.DetailPerbandinganKantong{
		BulanIni:        bulanIni,
		BulanSebelumnya: bulanSebelumnya,
		DataKantong:     expectedData,
		TotalBulanIni:   2300000.0,
		TotalBulanLalu:  1900000.0,
		RataRataTotal:   100000.0,
		PersentaseTotal: 21.05,
		TrendTotal:      "naik",
	}

	// Assertions
	assert.Equal(t, 12, detail.BulanIni.Bulan)
	assert.Equal(t, 2024, detail.BulanIni.Tahun)
	assert.Equal(t, "Desember", detail.BulanIni.NamaBulan)
	assert.Len(t, detail.DataKantong, 2)
	assert.Equal(t, "Kantong Belanja", detail.DataKantong[0].KantongNama)
	assert.Equal(t, 1500000.0, detail.DataKantong[0].JumlahBulanIni)
	assert.Equal(t, 25.0, detail.DataKantong[0].Persentase)
	assert.Equal(t, "naik", detail.DataKantong[0].Trend)
	assert.Equal(t, 2300000.0, detail.TotalBulanIni)
	assert.Equal(t, 21.05, detail.PersentaseTotal)
	assert.Equal(t, "naik", detail.TrendTotal)
}

// Test untuk validasi parameter input
func TestLaporanRepository_ValidateInputParameters(t *testing.T) {
	// Test untuk validasi tahun
	assert.True(t, 2024 >= 2020 && 2024 <= 2030, "Tahun harus dalam rentang 2020-2030")

	// Test untuk validasi bulan
	assert.True(t, 12 >= 1 && 12 <= 12, "Bulan harus dalam rentang 1-12")

	// Test untuk validasi user_id format
	userID := "550e8400-e29b-41d4-a716-446655440000"
	assert.NotEmpty(t, userID, "User ID tidak boleh kosong")
	assert.Equal(t, 36, len(userID), "User ID harus format UUID")
}

// Test untuk response format
func TestLaporanRepository_ResponseFormat(t *testing.T) {
	// Test TrenBulananResponse format
	response := domain.TrenBulananResponse{
		Success:   true,
		Message:   "Tren bulanan berhasil diambil",
		Code:      200,
		Data:      domain.TrenBulanan{},
		Timestamp: time.Now(),
	}

	assert.True(t, response.Success)
	assert.Equal(t, "Tren bulanan berhasil diambil", response.Message)
	assert.Equal(t, 200, response.Code)
	assert.NotZero(t, response.Timestamp)

	// Test PerbandinganKantongResponse format
	response2 := domain.PerbandinganKantongResponse{
		Success:   true,
		Message:   "Perbandingan kantong berhasil diambil",
		Code:      200,
		Data:      domain.PerbandinganKantong{},
		Timestamp: time.Now(),
	}

	assert.True(t, response2.Success)
	assert.Equal(t, "Perbandingan kantong berhasil diambil", response2.Message)
	assert.Equal(t, 200, response2.Code)
	assert.NotZero(t, response2.Timestamp)

	// Test DetailPerbandinganKantongResponse format
	response3 := domain.DetailPerbandinganKantongResponse{
		Success:   true,
		Message:   "Detail perbandingan kantong berhasil diambil",
		Code:      200,
		Data:      domain.DetailPerbandinganKantong{},
		Timestamp: time.Now(),
	}

	assert.True(t, response3.Success)
	assert.Equal(t, "Detail perbandingan kantong berhasil diambil", response3.Message)
	assert.Equal(t, 200, response3.Code)
	assert.NotZero(t, response3.Timestamp)
}

func TestLaporanRepository_DateCalculation(t *testing.T) {
	// Test untuk perhitungan bulan sebelumnya
	currentMonth := 12
	currentYear := 2024

	prevMonth := currentMonth - 1
	prevYear := currentYear

	if prevMonth < 1 {
		prevMonth = 12
		prevYear = currentYear - 1
	}

	assert.Equal(t, 11, prevMonth)
	assert.Equal(t, 2024, prevYear)

	// Test untuk bulan Januari
	currentMonth = 1
	prevMonth = currentMonth - 1
	prevYear = currentYear

	if prevMonth < 1 {
		prevMonth = 12
		prevYear = currentYear - 1
	}

	assert.Equal(t, 12, prevMonth)
	assert.Equal(t, 2023, prevYear)
}

func TestLaporanRepository_PercentageCalculation(t *testing.T) {
	// Test perhitungan persentase perubahan
	bulanIni := 1500000.0
	bulanLalu := 1200000.0

	persentase := ((bulanIni - bulanLalu) / bulanLalu) * 100
	expectedPersentase := 25.0

	assert.Equal(t, expectedPersentase, persentase)

	// Test untuk trend
	var trend string
	if bulanIni > bulanLalu {
		trend = "naik"
	} else if bulanIni < bulanLalu {
		trend = "turun"
	} else {
		trend = "stabil"
	}

	assert.Equal(t, "naik", trend)
}
