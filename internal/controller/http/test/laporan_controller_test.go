package http_test

import (
	"errors"
	"fiber-boiler-plate/internal/controller/http"
	"fiber-boiler-plate/internal/domain"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLaporanUsecase struct {
	mock.Mock
}

func (m *MockLaporanUsecase) GetRingkasanLaporan(userID uint, req *domain.RingkasanLaporanRequest) (*domain.RingkasanLaporanResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.RingkasanLaporanResponse), args.Error(1)
}

func (m *MockLaporanUsecase) GetStatistikTahunan(userID uint, req *domain.StatistikTahunanRequest) (*domain.StatistikTahunanResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StatistikTahunanResponse), args.Error(1)
}

func (m *MockLaporanUsecase) GetStatistikKantongBulanan(userID uint, req *domain.StatistikKantongBulananRequest) (*domain.StatistikKantongBulananResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StatistikKantongBulananResponse), args.Error(1)
}

func (m *MockLaporanUsecase) GetTopKantongPengeluaran(userID uint, req *domain.TopKantongPengeluaranRequest) (*domain.TopKantongPengeluaranResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TopKantongPengeluaranResponse), args.Error(1)
}

func (m *MockLaporanUsecase) GetStatistikKantongPeriode(userID uint, req *domain.StatistikKantongPeriodeRequest) (*domain.StatistikKantongPeriodeResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.StatistikKantongPeriodeResponse), args.Error(1)
}

func (m *MockLaporanUsecase) GetPengeluaranKantongDetail(userID uint, req *domain.PengeluaranKantongDetailRequest) (*domain.PengeluaranKantongDetailResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PengeluaranKantongDetailResponse), args.Error(1)
}

func (m *MockLaporanUsecase) GetTrenBulanan(userID uint, req *domain.TrenBulananRequest) (*domain.TrenBulananResponse, error) {
	args := m.Called(userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TrenBulananResponse), args.Error(1)
}

func (m *MockLaporanUsecase) GetPerbandinganKantong(userID uint) (*domain.PerbandinganKantongResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PerbandinganKantongResponse), args.Error(1)
}

func (m *MockLaporanUsecase) GetDetailPerbandinganKantong(userID uint) (*domain.DetailPerbandinganKantongResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DetailPerbandinganKantongResponse), args.Error(1)
}

func setupLaporanController() (*fiber.App, *MockLaporanUsecase) {
	app := fiber.New()
	mockUsecase := new(MockLaporanUsecase)
	controller := http.NewLaporanController(mockUsecase)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return c.Next()
	})

	app.Get("/laporan/tren/bulanan", controller.GetTrenBulanan)
	app.Get("/laporan/perbandingan/kantong", controller.GetPerbandinganKantong)
	app.Get("/laporan/perbandingan/kantong/detail", controller.GetDetailPerbandinganKantong)

	return app, mockUsecase
}

func TestLaporanController_GetTrenBulanan_Success(t *testing.T) {
	app, mockUsecase := setupLaporanController()

	expectedResponse := &domain.TrenBulananResponse{
		Success: true,
		Message: "Data tren bulanan berhasil diambil",
		Code:    200,
		Data: domain.TrenBulanan{
			Tahun: 2024,
			DataTren: []domain.DataBulanan{
				{
					Bulan:            1,
					NamaBulan:        "Januari",
					TotalPemasukan:   5000000,
					TotalPengeluaran: 3500000,
				},
			},
			TotalPemasukanTahun:   62400000,
			TotalPengeluaranTahun: 45600000,
		},
		Timestamp: time.Now(),
	}

	mockUsecase.On("GetTrenBulanan", uint(1), mock.AnythingOfType("*domain.TrenBulananRequest")).Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/laporan/tren/bulanan?tahun=2024", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}

func TestLaporanController_GetTrenBulanan_InvalidTahun(t *testing.T) {
	app, mockUsecase := setupLaporanController()

	req := httptest.NewRequest("GET", "/laporan/tren/bulanan?tahun=invalid", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockUsecase.AssertNotCalled(t, "GetTrenBulanan")
}

func TestLaporanController_GetTrenBulanan_UsecaseError(t *testing.T) {
	app, mockUsecase := setupLaporanController()

	mockUsecase.On("GetTrenBulanan", uint(1), mock.AnythingOfType("*domain.TrenBulananRequest")).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/laporan/tren/bulanan", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}

func TestLaporanController_GetPerbandinganKantong_Success(t *testing.T) {
	app, mockUsecase := setupLaporanController()

	expectedResponse := &domain.PerbandinganKantongResponse{
		Success: true,
		Message: "Perbandingan pengeluaran per kantong berhasil diambil",
		Code:    200,
		Data: domain.PerbandinganKantong{
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
		},
		Timestamp: time.Now(),
	}

	mockUsecase.On("GetPerbandinganKantong", uint(1)).Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/laporan/perbandingan/kantong", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}

func TestLaporanController_GetPerbandinganKantong_UsecaseError(t *testing.T) {
	app, mockUsecase := setupLaporanController()

	mockUsecase.On("GetPerbandinganKantong", uint(1)).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/laporan/perbandingan/kantong", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}

func TestLaporanController_GetDetailPerbandinganKantong_Success(t *testing.T) {
	app, mockUsecase := setupLaporanController()

	expectedResponse := &domain.DetailPerbandinganKantongResponse{
		Success: true,
		Message: "Detail perbandingan pengeluaran per kantong berhasil diambil",
		Code:    200,
		Data: domain.DetailPerbandinganKantong{
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
		},
		Timestamp: time.Now(),
	}

	mockUsecase.On("GetDetailPerbandinganKantong", uint(1)).Return(expectedResponse, nil)

	req := httptest.NewRequest("GET", "/laporan/perbandingan/kantong/detail", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}

func TestLaporanController_GetDetailPerbandinganKantong_UsecaseError(t *testing.T) {
	app, mockUsecase := setupLaporanController()

	mockUsecase.On("GetDetailPerbandinganKantong", uint(1)).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/laporan/perbandingan/kantong/detail", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	mockUsecase.AssertExpectations(t)
}
