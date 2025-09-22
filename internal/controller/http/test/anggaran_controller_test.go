package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"fiber-boiler-plate/internal/controller/http"
	"fiber-boiler-plate/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAnggaranUsecase struct {
	mock.Mock
}

func (m *MockAnggaranUsecase) GetAnggaranList(userID uint, req *domain.AnggaranListRequest) ([]*domain.AnggaranResponse, *domain.PaginationMeta, error) {
	args := m.Called(userID, req)
	return args.Get(0).([]*domain.AnggaranResponse), args.Get(1).(*domain.PaginationMeta), args.Error(2)
}

func (m *MockAnggaranUsecase) GetAnggaranDetail(kantongID string, userID uint, bulan, tahun *int) (*domain.AnggaranDetailResponse, error) {
	args := m.Called(kantongID, userID, bulan, tahun)
	return args.Get(0).(*domain.AnggaranDetailResponse), args.Error(1)
}

func (m *MockAnggaranUsecase) CreatePenyesuaianAnggaran(userID uint, req *domain.PenyesuaianAnggaranRequest) (*domain.AnggaranResponse, error) {
	args := m.Called(userID, req)
	return args.Get(0).(*domain.AnggaranResponse), args.Error(1)
}

func (m *MockAnggaranUsecase) CreateAnggaranForNewKantong(kantong *domain.Kantong) error {
	args := m.Called(kantong)
	return args.Error(0)
}

func (m *MockAnggaranUsecase) UpdateAnggaranAfterTransaction(kantongID string, userID uint) error {
	args := m.Called(kantongID, userID)
	return args.Error(0)
}

func setupAnggaranTest() (*fiber.App, *MockAnggaranUsecase) {
	app := fiber.New()
	mockUsecase := &MockAnggaranUsecase{}
	controller := http.NewAnggaranController(mockUsecase)

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user_id", uint(1))
		return c.Next()
	})

	app.Get("/anggaran", controller.GetAnggaranList)
	app.Get("/anggaran/:kantong_id", controller.GetAnggaranDetail)
	app.Post("/anggaran/penyesuaian", controller.CreatePenyesuaianAnggaran)

	return app, mockUsecase
}

func TestGetAnggaranList_Success(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	mockResponses := []*domain.AnggaranResponse{
		{
			KantongID:   "550e8400-e29b-41d4-a716-446655440001",
			NamaKantong: "Kantong Test",
			Rencana:     &[]float64{1000000}[0],
			CarryIn:     0,
			Penyesuaian: 0,
			Terpakai:    450000,
			Sisa:        550000,
			Progres:     45.0,
			Bulan:       9,
			Tahun:       2024,
		},
	}

	mockMeta := &domain.PaginationMeta{
		CurrentPage:  1,
		TotalPages:   1,
		TotalRecords: 1,
		PerPage:      10,
	}

	mockUsecase.On("GetAnggaranList", uint(1), mock.AnythingOfType("*domain.AnggaranListRequest")).
		Return(mockResponses, mockMeta, nil)

	req := httptest.NewRequest("GET", "/anggaran?bulan=9&tahun=2024", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestGetAnggaranDetail_Success(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	mockResponse := &domain.AnggaranDetailResponse{
		KantongID:   "550e8400-e29b-41d4-a716-446655440001",
		NamaKantong: "Kantong Test",
		DetailKantong: &domain.KantongResponse{
			ID:       "550e8400-e29b-41d4-a716-446655440001",
			IDKartu:  "K4N7G1",
			Nama:     "Kantong Test",
			Kategori: "Pengeluaran",
			Saldo:    750000,
			Warna:    "Navy",
		},
		Rencana:        &[]float64{1000000}[0],
		CarryIn:        0,
		Penyesuaian:    0,
		Terpakai:       450000,
		Sisa:           550000,
		Progres:        45.0,
		StatistikBulan: []domain.StatistikHarian{},
		Bulan:          9,
		Tahun:          2024,
	}

	mockUsecase.On("GetAnggaranDetail", "550e8400-e29b-41d4-a716-446655440001", uint(1),
		mock.AnythingOfType("*int"), mock.AnythingOfType("*int")).
		Return(mockResponse, nil)

	req := httptest.NewRequest("GET", "/anggaran/550e8400-e29b-41d4-a716-446655440001?bulan=9&tahun=2024", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestGetAnggaranDetail_NotFound(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	mockUsecase.On("GetAnggaranDetail", "invalid-id", uint(1),
		mock.AnythingOfType("*int"), mock.AnythingOfType("*int")).
		Return((*domain.AnggaranDetailResponse)(nil), errors.New("kantong tidak ditemukan"))

	req := httptest.NewRequest("GET", "/anggaran/invalid-id", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestCreatePenyesuaianAnggaran_Success(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	reqBody := domain.PenyesuaianAnggaranRequest{
		KantongID: "550e8400-e29b-41d4-a716-446655440001",
		Jenis:     "kurangi",
		Jumlah:    50000,
		Bulan:     9,
		Tahun:     2024,
	}

	mockResponse := &domain.AnggaranResponse{
		KantongID:   "550e8400-e29b-41d4-a716-446655440001",
		NamaKantong: "Kantong Test",
		Rencana:     &[]float64{1000000}[0],
		CarryIn:     0,
		Penyesuaian: -50000,
		Terpakai:    450000,
		Sisa:        500000,
		Progres:     47.37,
		Bulan:       9,
		Tahun:       2024,
	}

	mockUsecase.On("CreatePenyesuaianAnggaran", uint(1), &reqBody).Return(mockResponse, nil)

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/anggaran/penyesuaian", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestCreatePenyesuaianAnggaran_ValidationError(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	reqBody := domain.PenyesuaianAnggaranRequest{
		KantongID: "",
		Jenis:     "invalid",
		Jumlah:    -100,
		Bulan:     13,
		Tahun:     2019,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/anggaran/penyesuaian", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
	mockUsecase.AssertNotCalled(t, "CreatePenyesuaianAnggaran")
}

func TestCreatePenyesuaianAnggaran_KantongNotFound(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	reqBody := domain.PenyesuaianAnggaranRequest{
		KantongID: "550e8400-e29b-41d4-a716-446655440001",
		Jenis:     "tambah",
		Jumlah:    50000,
		Bulan:     9,
		Tahun:     2024,
	}

	mockUsecase.On("CreatePenyesuaianAnggaran", uint(1), &reqBody).
		Return((*domain.AnggaranResponse)(nil), errors.New("kantong tidak ditemukan"))

	jsonBody, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/anggaran/penyesuaian", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestGetAnggaranList_WithSearchAndSort(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	mockResponses := []*domain.AnggaranResponse{
		{
			KantongID:   "550e8400-e29b-41d4-a716-446655440001",
			NamaKantong: "Kantong Belanja",
			Rencana:     &[]float64{1000000}[0],
			CarryIn:     0,
			Penyesuaian: 0,
			Terpakai:    450000,
			Sisa:        550000,
			Progres:     45.0,
			Bulan:       9,
			Tahun:       2024,
		},
	}

	mockMeta := &domain.PaginationMeta{
		CurrentPage:  1,
		TotalPages:   1,
		TotalRecords: 1,
		PerPage:      10,
	}

	mockUsecase.On("GetAnggaranList", uint(1), mock.MatchedBy(func(req *domain.AnggaranListRequest) bool {
		return req.Search != nil && *req.Search == "Belanja" &&
			req.SortBy == "sisa" &&
			req.SortDirection == "desc"
	})).Return(mockResponses, mockMeta, nil)

	req := httptest.NewRequest("GET", "/anggaran?search=Belanja&sort_by=sisa&sort_direction=desc", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	mockUsecase.AssertExpectations(t)
}

func TestGetAnggaranDetail_InvalidBulan(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	req := httptest.NewRequest("GET", "/anggaran/550e8400-e29b-41d4-a716-446655440001?bulan=13", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
	mockUsecase.AssertNotCalled(t, "GetAnggaranDetail")
}

func TestGetAnggaranDetail_InvalidTahun(t *testing.T) {
	app, mockUsecase := setupAnggaranTest()

	req := httptest.NewRequest("GET", "/anggaran/550e8400-e29b-41d4-a716-446655440001?tahun=2019", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
	mockUsecase.AssertNotCalled(t, "GetAnggaranDetail")
}
