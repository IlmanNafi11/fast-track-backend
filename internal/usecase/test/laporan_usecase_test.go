package usecase_test

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRedisRepository struct {
	mock.Mock
}

func (m *MockRedisRepository) Set(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockRedisRepository) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisRepository) GetJSON(key string, dest interface{}) error {
	args := m.Called(key, dest)
	return args.Error(0)
}

func (m *MockRedisRepository) SetJSON(key string, value interface{}, ttl time.Duration) error {
	args := m.Called(key, value, ttl)
	return args.Error(0)
}

func (m *MockRedisRepository) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockRedisRepository) Exists(key string) (bool, error) {
	args := m.Called(key)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisRepository) Increment(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisRepository) Decrement(key string) (int64, error) {
	args := m.Called(key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisRepository) SetExpire(key string, ttl time.Duration) error {
	args := m.Called(key, ttl)
	return args.Error(0)
}

func (m *MockRedisRepository) GetTTL(key string) (time.Duration, error) {
	args := m.Called(key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisRepository) GetKeys(pattern string) ([]string, error) {
	args := m.Called(pattern)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisRepository) FlushAll() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisRepository) Ping() error {
	args := m.Called()
	return args.Error(0)
}

type MockLaporanRepository struct {
	mock.Mock
}

func (m *MockLaporanRepository) GetRingkasanLaporan(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.RingkasanLaporan, error) {
	args := m.Called(userID, tanggalMulai, tanggalSelesai)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RingkasanLaporan), args.Error(1)
}

func (m *MockLaporanRepository) GetStatistikTahunan(userID uint, tahun int) (*domain.StatistikTahunan, error) {
	args := m.Called(userID, tahun)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StatistikTahunan), args.Error(1)
}

func (m *MockLaporanRepository) GetStatistikKantongBulanan(userID uint, bulan, tahun int) (*domain.StatistikKantongBulanan, error) {
	args := m.Called(userID, bulan, tahun)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StatistikKantongBulanan), args.Error(1)
}

func (m *MockLaporanRepository) GetTopKantongPengeluaran(userID uint, bulan, tahun, limit int) (*domain.TopKantongPengeluaran, error) {
	args := m.Called(userID, bulan, tahun, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TopKantongPengeluaran), args.Error(1)
}

func (m *MockLaporanRepository) GetStatistikKantongPeriode(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.StatistikKantongPeriode, error) {
	args := m.Called(userID, tanggalMulai, tanggalSelesai)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StatistikKantongPeriode), args.Error(1)
}

func (m *MockLaporanRepository) GetPengeluaranKantongDetail(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.PengeluaranKantongDetail, error) {
	args := m.Called(userID, tanggalMulai, tanggalSelesai)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PengeluaranKantongDetail), args.Error(1)
}

func (m *MockLaporanRepository) GetTrenBulanan(userID uint, tahun int) (*domain.TrenBulanan, error) {
	args := m.Called(userID, tahun)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TrenBulanan), args.Error(1)
}

func (m *MockLaporanRepository) GetPerbandinganKantong(userID uint, bulanIni, tahunIni, bulanLalu, tahunLalu int) (*domain.PerbandinganKantong, error) {
	args := m.Called(userID, bulanIni, tahunIni, bulanLalu, tahunLalu)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PerbandinganKantong), args.Error(1)
}

func (m *MockLaporanRepository) GetDetailPerbandinganKantong(userID uint, bulanIni, tahunIni, bulanLalu, tahunLalu int) (*domain.DetailPerbandinganKantong, error) {
	args := m.Called(userID, bulanIni, tahunIni, bulanLalu, tahunLalu)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DetailPerbandinganKantong), args.Error(1)
}

func TestLaporanUsecase_GetRingkasanLaporan_Success(t *testing.T) {
	mockLaporanRepo := new(MockLaporanRepository)
	mockRedisRepo := new(MockRedisRepository)
	laporanUsecase := usecase.NewLaporanUsecase(mockLaporanRepo, mockRedisRepo)

	userID := uint(1)
	req := &domain.RingkasanLaporanRequest{}

	expectedData := &domain.RingkasanLaporan{
		TotalPemasukan:            5000000,
		TotalPengeluaran:          3500000,
		TotalSaldo:                12500000,
		RataRataPengeluaranHarian: 112903.23,
		Periode: domain.PeriodeTanggal{
			TanggalMulai:   "2024-01-01",
			TanggalSelesai: "2024-01-31",
		},
	}

	mockRedisRepo.On("GetJSON", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.RingkasanLaporanResponse")).Return(assert.AnError)
	mockLaporanRepo.On("GetRingkasanLaporan", userID, mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).Return(expectedData, nil)
	mockRedisRepo.On("SetJSON", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.RingkasanLaporanResponse"), mock.AnythingOfType("time.Duration")).Return(nil)

	result, err := laporanUsecase.GetRingkasanLaporan(userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "Ringkasan laporan berhasil diambil", result.Message)
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, expectedData.TotalPemasukan, result.Data.TotalPemasukan)
	assert.Equal(t, expectedData.TotalPengeluaran, result.Data.TotalPengeluaran)
	mockLaporanRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

func TestLaporanUsecase_GetStatistikTahunan_Success(t *testing.T) {
	mockLaporanRepo := new(MockLaporanRepository)
	mockRedisRepo := new(MockRedisRepository)
	laporanUsecase := usecase.NewLaporanUsecase(mockLaporanRepo, mockRedisRepo)

	userID := uint(1)
	tahun := 2024
	req := &domain.StatistikTahunanRequest{Tahun: &tahun}

	expectedData := &domain.StatistikTahunan{
		Tahun: 2024,
		DataBulanan: []domain.DataBulanan{
			{
				Bulan:            1,
				NamaBulan:        "Januari",
				TotalPemasukan:   5000000,
				TotalPengeluaran: 3500000,
			},
		},
		TotalPemasukanTahun:   60000000,
		TotalPengeluaranTahun: 42000000,
	}

	mockRedisRepo.On("GetJSON", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.StatistikTahunanResponse")).Return(assert.AnError)
	mockLaporanRepo.On("GetStatistikTahunan", userID, tahun).Return(expectedData, nil)
	mockRedisRepo.On("SetJSON", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.StatistikTahunanResponse"), mock.AnythingOfType("time.Duration")).Return(nil)

	result, err := laporanUsecase.GetStatistikTahunan(userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "Statistik tahunan berhasil diambil", result.Message)
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, expectedData.Tahun, result.Data.Tahun)
	assert.Equal(t, expectedData.TotalPemasukanTahun, result.Data.TotalPemasukanTahun)
	mockLaporanRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

func TestLaporanUsecase_GetStatistikKantongPeriode_Success(t *testing.T) {
	mockLaporanRepo := new(MockLaporanRepository)
	mockRedisRepo := new(MockRedisRepository)
	laporanUsecase := usecase.NewLaporanUsecase(mockLaporanRepo, mockRedisRepo)

	userID := uint(1)
	tanggalMulai := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tanggalSelesai := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	tanggalMulaiStr := "2024-01-01"
	tanggalSelesaiStr := "2024-01-31"
	req := &domain.StatistikKantongPeriodeRequest{
		TanggalMulai:   &tanggalMulaiStr,
		TanggalSelesai: &tanggalSelesaiStr,
	}

	expectedData := &domain.StatistikKantongPeriode{
		Periode: domain.PeriodeTanggal{
			TanggalMulai:   "2024-01-01",
			TanggalSelesai: "2024-01-31",
		},
		DataKantong: []domain.DataKantongPeriode{
			{
				KantongID:        "1",
				KantongNama:      "Wallet Utama",
				TotalPengeluaran: 3500000,
			},
		},
		TotalPengeluaran: 3500000,
	}

	mockRedisRepo.On("GetJSON", mock.AnythingOfType("string"), mock.Anything).Return(assert.AnError)
	mockLaporanRepo.On("GetStatistikKantongPeriode", userID, tanggalMulai, tanggalSelesai).Return(expectedData, nil)
	mockRedisRepo.On("SetJSON", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	result, err := laporanUsecase.GetStatistikKantongPeriode(userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "Statistik kantong periode berhasil diambil", result.Message)
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, expectedData.TotalPengeluaran, result.Data.TotalPengeluaran)
	assert.Len(t, result.Data.DataKantong, 1)
	assert.Equal(t, "Wallet Utama", result.Data.DataKantong[0].KantongNama)
	mockLaporanRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

func TestLaporanUsecase_GetPengeluaranKantongDetail_Success(t *testing.T) {
	mockLaporanRepo := new(MockLaporanRepository)
	mockRedisRepo := new(MockRedisRepository)
	laporanUsecase := usecase.NewLaporanUsecase(mockLaporanRepo, mockRedisRepo)

	userID := uint(1)
	tanggalMulai := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	tanggalSelesai := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	tanggalMulaiStr := "2024-01-01"
	tanggalSelesaiStr := "2024-01-31"
	req := &domain.PengeluaranKantongDetailRequest{
		TanggalMulai:   &tanggalMulaiStr,
		TanggalSelesai: &tanggalSelesaiStr,
	}

	expectedData := &domain.PengeluaranKantongDetail{
		Periode: domain.PeriodeTanggal{
			TanggalMulai:   "2024-01-01",
			TanggalSelesai: "2024-01-31",
		},
		DataKantong: []domain.DataKantongDetail{
			{
				KantongID:           "1",
				KantongNama:         "Wallet Utama",
				TotalPengeluaran:    3500000,
				PersentaseDariSaldo: 75.5,
				JumlahTransaksi:     25,
				RataRataPengeluaran: 140000,
				SaldoKantong:        2000000,
			},
		},
		TotalPengeluaran:       3500000,
		TotalSaldoSemuaKantong: 2000000,
	}

	mockRedisRepo.On("GetJSON", mock.AnythingOfType("string"), mock.Anything).Return(assert.AnError)
	mockLaporanRepo.On("GetPengeluaranKantongDetail", userID, tanggalMulai, tanggalSelesai).Return(expectedData, nil)
	mockRedisRepo.On("SetJSON", mock.AnythingOfType("string"), mock.Anything, mock.AnythingOfType("time.Duration")).Return(nil)

	result, err := laporanUsecase.GetPengeluaranKantongDetail(userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "Detail pengeluaran kantong berhasil diambil", result.Message)
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, expectedData.TotalPengeluaran, result.Data.TotalPengeluaran)
	assert.Equal(t, expectedData.TotalSaldoSemuaKantong, result.Data.TotalSaldoSemuaKantong)
	assert.Len(t, result.Data.DataKantong, 1)
	assert.Equal(t, "Wallet Utama", result.Data.DataKantong[0].KantongNama)
	mockLaporanRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}


func TestLaporanUsecase_GetTrenBulanan_Success(t *testing.T) {
	mockLaporanRepo := new(MockLaporanRepository)
	mockRedisRepo := new(MockRedisRepository)
	laporanUsecase := usecase.NewLaporanUsecase(mockLaporanRepo, mockRedisRepo)

	userID := uint(1)
	tahun := 2024

	expectedData := &domain.TrenBulanan{
		Tahun: 2024,
		DataTren: []domain.DataBulanan{
			{
				Bulan:            1,
				NamaBulan:        "Januari",
				TotalPemasukan:   5000000,
				TotalPengeluaran: 3500000,
			},
			{
				Bulan:            2,
				NamaBulan:        "Februari",
				TotalPemasukan:   5200000,
				TotalPengeluaran: 3800000,
			},
		},
		TotalPemasukanTahun:   62400000,
		TotalPengeluaranTahun: 45600000,
	}

	cacheKey := "tren_bulanan:1:2024"

	mockRedisRepo.On("GetJSON", cacheKey, mock.AnythingOfType("**domain.TrenBulananResponse")).Return(assert.AnError)

	mockLaporanRepo.On("GetTrenBulanan", userID, tahun).Return(expectedData, nil)

	mockRedisRepo.On("SetJSON", cacheKey, mock.AnythingOfType("*domain.TrenBulananResponse"), 10*time.Minute).Return(nil)

	req := &domain.TrenBulananRequest{
		Tahun: &tahun,
	}

	result, err := laporanUsecase.GetTrenBulanan(userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "Data tren bulanan berhasil diambil", result.Message)
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, expectedData.Tahun, result.Data.Tahun)
	assert.Equal(t, expectedData.TotalPemasukanTahun, result.Data.TotalPemasukanTahun)
	assert.Equal(t, expectedData.TotalPengeluaranTahun, result.Data.TotalPengeluaranTahun)
	assert.Len(t, result.Data.DataTren, 2)
	assert.Equal(t, "Januari", result.Data.DataTren[0].NamaBulan)
	mockLaporanRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

func TestLaporanUsecase_GetPerbandinganKantong_Success(t *testing.T) {
	mockLaporanRepo := new(MockLaporanRepository)
	mockRedisRepo := new(MockRedisRepository)
	laporanUsecase := usecase.NewLaporanUsecase(mockLaporanRepo, mockRedisRepo)

	userID := uint(1)

	expectedData := &domain.PerbandinganKantong{
		BulanIni: domain.PeriodeBulan{
			Bulan:     12,
			Tahun:     2024,
			NamaBulan: "Desember",
		},
		BulanSebelumnya: domain.PeriodeBulan{
			Bulan:     11,
			Tahun:     2024,
			NamaBulan: "November",
		},
		DataKantong: []domain.DataPerbandinganKantong{
			{
				KantongID:       "550e8400-e29b-41d4-a716-446655440001",
				KantongNama:     "Kantong Belanja",
				JumlahBulanIni:  1500000,
				JumlahBulanLalu: 1200000,
			},
		},
		TotalBulanIni:  2300000,
		TotalBulanLalu: 1900000,
	}

	mockLaporanRepo.On("GetPerbandinganKantong", userID, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(expectedData, nil)

	mockRedisRepo.On("GetJSON", mock.AnythingOfType("string"), mock.AnythingOfType("**domain.PerbandinganKantongResponse")).Return(assert.AnError)
	mockRedisRepo.On("SetJSON", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.PerbandinganKantongResponse"), 15*time.Minute).Return(nil)

	result, err := laporanUsecase.GetPerbandinganKantong(userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "Perbandingan pengeluaran per kantong berhasil diambil", result.Message)
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, expectedData.BulanIni.Bulan, result.Data.BulanIni.Bulan)
	assert.Equal(t, expectedData.BulanSebelumnya.Bulan, result.Data.BulanSebelumnya.Bulan)
	assert.Equal(t, expectedData.TotalBulanIni, result.Data.TotalBulanIni)
	assert.Equal(t, expectedData.TotalBulanLalu, result.Data.TotalBulanLalu)
	assert.Len(t, result.Data.DataKantong, 1)
	assert.Equal(t, "Kantong Belanja", result.Data.DataKantong[0].KantongNama)
	mockLaporanRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}

func TestLaporanUsecase_GetDetailPerbandinganKantong_Success(t *testing.T) {
	mockLaporanRepo := new(MockLaporanRepository)
	mockRedisRepo := new(MockRedisRepository)
	laporanUsecase := usecase.NewLaporanUsecase(mockLaporanRepo, mockRedisRepo)

	userID := uint(1)

	expectedData := &domain.DetailPerbandinganKantong{
		BulanIni: domain.PeriodeBulan{
			Bulan:     12,
			Tahun:     2024,
			NamaBulan: "Desember",
		},
		BulanSebelumnya: domain.PeriodeBulan{
			Bulan:     11,
			Tahun:     2024,
			NamaBulan: "November",
		},
		DataKantong: []domain.DataDetailPerbandinganKantong{
			{
				KantongID:           "550e8400-e29b-41d4-a716-446655440001",
				KantongNama:         "Kantong Belanja",
				JumlahBulanIni:      1500000,
				JumlahBulanLalu:     1200000,
				RataRataPengeluaran: 100000,
				Persentase:          25.0,
				Trend:               "naik",
			},
		},
		TotalBulanIni:   2300000,
		TotalBulanLalu:  1900000,
		RataRataTotal:   100000,
		PersentaseTotal: 21.05,
		TrendTotal:      "naik",
	}

	mockLaporanRepo.On("GetDetailPerbandinganKantong", userID, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(expectedData, nil)

	mockRedisRepo.On("GetJSON", mock.AnythingOfType("string"), mock.AnythingOfType("**domain.DetailPerbandinganKantongResponse")).Return(assert.AnError)
	mockRedisRepo.On("SetJSON", mock.AnythingOfType("string"), mock.AnythingOfType("*domain.DetailPerbandinganKantongResponse"), 18*time.Minute).Return(nil)

	result, err := laporanUsecase.GetDetailPerbandinganKantong(userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Success)
	assert.Equal(t, "Detail perbandingan pengeluaran per kantong berhasil diambil", result.Message)
	assert.Equal(t, 200, result.Code)
	assert.Equal(t, expectedData.BulanIni.Bulan, result.Data.BulanIni.Bulan)
	assert.Equal(t, expectedData.BulanSebelumnya.Bulan, result.Data.BulanSebelumnya.Bulan)
	assert.Equal(t, expectedData.TotalBulanIni, result.Data.TotalBulanIni)
	assert.Equal(t, expectedData.TotalBulanLalu, result.Data.TotalBulanLalu)
	assert.Equal(t, expectedData.RataRataTotal, result.Data.RataRataTotal)
	assert.Equal(t, expectedData.PersentaseTotal, result.Data.PersentaseTotal)
	assert.Equal(t, expectedData.TrendTotal, result.Data.TrendTotal)
	assert.Len(t, result.Data.DataKantong, 1)
	assert.Equal(t, "Kantong Belanja", result.Data.DataKantong[0].KantongNama)
	assert.Equal(t, 25.0, result.Data.DataKantong[0].Persentase)
	assert.Equal(t, "naik", result.Data.DataKantong[0].Trend)
	mockLaporanRepo.AssertExpectations(t)
	mockRedisRepo.AssertExpectations(t)
}
