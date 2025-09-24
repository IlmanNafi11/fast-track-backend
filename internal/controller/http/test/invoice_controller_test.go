package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fiber-boiler-plate/internal/controller/http"
	"fiber-boiler-plate/internal/domain"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockInvoiceUsecase struct {
	mock.Mock
}

func (m *MockInvoiceUsecase) GetAll(req *domain.InvoiceListRequest) ([]*domain.InvoiceResponse, *domain.PaginationMeta, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).([]*domain.InvoiceResponse), args.Get(1).(*domain.PaginationMeta), args.Error(2)
}

func (m *MockInvoiceUsecase) GetByID(id string) (*domain.InvoiceDetailResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InvoiceDetailResponse), args.Error(1)
}

func (m *MockInvoiceUsecase) UpdateStatus(id string, req *domain.UpdateInvoiceStatusRequest) (*domain.InvoiceDetailResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InvoiceDetailResponse), args.Error(1)
}

func (m *MockInvoiceUsecase) GetStatistics(req *domain.InvoiceStatisticsRequest) (*domain.InvoiceStatistics, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InvoiceStatistics), args.Error(1)
}

func TestNewInvoiceController(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	assert.NotNil(t, controller)
}

func TestInvoiceController_GetAll_Success(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice", controller.GetAll)

	now := time.Now()
	userID := uint(1)
	planID := uuid.New()

	mockInvoices := []*domain.InvoiceResponse{
		{
			InvoiceID:            "INV-2024-001",
			NamaUser:             "John Doe",
			UserID:               userID,
			UserEmail:            "john@example.com",
			Jumlah:               99000,
			Status:               "sukses",
			DibayarPada:          &now,
			MetodePembayaran:     stringPtr("Bank Transfer"),
			Keterangan:           stringPtr("Pembayaran PRO"),
			SubscriptionPlanID:   &planID,
			SubscriptionPlanNama: stringPtr("PRO Monthly"),
			CreatedAt:            now,
			UpdatedAt:            now,
		},
	}

	mockMeta := &domain.PaginationMeta{
		CurrentPage:  1,
		TotalPages:   1,
		TotalRecords: 1,
		PerPage:      10,
	}

	mockUsecase.On("GetAll", mock.MatchedBy(func(req *domain.InvoiceListRequest) bool {
		return req.Page == 1 && req.PerPage == 10
	})).Return(mockInvoices, mockMeta, nil)

	req := httptest.NewRequest("GET", "/invoice?page=1&per_page=10", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Daftar invoice berhasil diambil", response["message"])
	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])
	assert.NotNil(t, response["meta"])

	mockUsecase.AssertExpectations(t)
}

func TestInvoiceController_GetAll_ValidationError(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice", controller.GetAll)

	req := httptest.NewRequest("GET", "/invoice?page=0&per_page=101", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Data validasi tidak valid", response["message"])
	assert.Equal(t, float64(400), response["code"])

	mockUsecase.AssertNotCalled(t, "GetAll")
}

func TestInvoiceController_GetAll_UsecaseError(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice", controller.GetAll)

	mockUsecase.On("GetAll", mock.AnythingOfType("*domain.InvoiceListRequest")).Return(([]*domain.InvoiceResponse)(nil), (*domain.PaginationMeta)(nil), errors.New("database error"))

	req := httptest.NewRequest("GET", "/invoice?page=1&per_page=10", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Terjadi kesalahan pada server", response["message"])
	assert.Equal(t, float64(500), response["code"])

	mockUsecase.AssertExpectations(t)
}

func TestInvoiceController_GetByID_Success(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice/:invoice_id", controller.GetByID)

	now := time.Now()
	userID := uint(1)
	planID := uuid.New()
	invoiceID := "INV-2024-001"

	mockInvoice := &domain.InvoiceDetailResponse{
		InvoiceID:            invoiceID,
		NamaUser:             "John Doe",
		UserID:               userID,
		UserEmail:            "john@example.com",
		Jumlah:               99000,
		Status:               "sukses",
		DibayarPada:          &now,
		MetodePembayaran:     stringPtr("Bank Transfer"),
		Keterangan:           stringPtr("Pembayaran PRO"),
		SubscriptionPlanID:   &planID,
		SubscriptionPlanNama: stringPtr("PRO Monthly"),
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	mockUsecase.On("GetByID", invoiceID).Return(mockInvoice, nil)

	req := httptest.NewRequest("GET", "/invoice/"+invoiceID, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Detail invoice berhasil diambil", response["message"])
	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])

	mockUsecase.AssertExpectations(t)
}

func TestInvoiceController_GetByID_EmptyID(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice/:invoice_id", controller.GetByID)

	req := httptest.NewRequest("GET", "/invoice/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockUsecase.AssertNotCalled(t, "GetByID")
}

func TestInvoiceController_GetByID_NotFound(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice/:invoice_id", controller.GetByID)

	invoiceID := "INV-2024-999"

	mockUsecase.On("GetByID", invoiceID).Return(nil, errors.New("invoice tidak ditemukan"))

	req := httptest.NewRequest("GET", "/invoice/"+invoiceID, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "invoice tidak ditemukan", response["message"])
	assert.Equal(t, float64(404), response["code"])

	mockUsecase.AssertExpectations(t)
}

func TestInvoiceController_GetByID_InternalError(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice/:invoice_id", controller.GetByID)

	invoiceID := "INV-2024-001"

	mockUsecase.On("GetByID", invoiceID).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/invoice/"+invoiceID, nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Terjadi kesalahan pada server", response["message"])
	assert.Equal(t, float64(500), response["code"])

	mockUsecase.AssertExpectations(t)
}

func TestInvoiceController_UpdateStatus_Success(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Patch("/invoice/:invoice_id", controller.UpdateStatus)

	now := time.Now()
	userID := uint(1)
	planID := uuid.New()
	invoiceID := "INV-2024-001"
	keterangan := "Pembayaran telah dikonfirmasi"

	updateReq := domain.UpdateInvoiceStatusRequest{
		Status:     "sukses",
		Keterangan: &keterangan,
	}

	mockInvoice := &domain.InvoiceDetailResponse{
		InvoiceID:            invoiceID,
		NamaUser:             "John Doe",
		UserID:               userID,
		UserEmail:            "john@example.com",
		Jumlah:               99000,
		Status:               "sukses",
		DibayarPada:          &now,
		MetodePembayaran:     stringPtr("Bank Transfer"),
		Keterangan:           &keterangan,
		SubscriptionPlanID:   &planID,
		SubscriptionPlanNama: stringPtr("PRO Monthly"),
		CreatedAt:            now,
		UpdatedAt:            now,
	}

	mockUsecase.On("UpdateStatus", invoiceID, &updateReq).Return(mockInvoice, nil)

	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PATCH", "/invoice/"+invoiceID, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Status invoice berhasil diperbarui", response["message"])
	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])

	mockUsecase.AssertExpectations(t)
}

func TestInvoiceController_UpdateStatus_ValidationError(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Patch("/invoice/:invoice_id", controller.UpdateStatus)

	invoiceID := "INV-2024-001"

	invalidReq := map[string]interface{}{
		"status": "invalid_status",
	}

	reqBody, _ := json.Marshal(invalidReq)
	req := httptest.NewRequest("PATCH", "/invoice/"+invoiceID, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Data validasi tidak valid", response["message"])
	assert.Equal(t, float64(400), response["code"])

	mockUsecase.AssertNotCalled(t, "UpdateStatus")
}

func TestInvoiceController_UpdateStatus_InvalidJSON(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Patch("/invoice/:invoice_id", controller.UpdateStatus)

	invoiceID := "INV-2024-001"

	req := httptest.NewRequest("PATCH", "/invoice/"+invoiceID, bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Format data tidak valid", response["message"])
	assert.Equal(t, float64(400), response["code"])

	mockUsecase.AssertNotCalled(t, "UpdateStatus")
}

func TestInvoiceController_UpdateStatus_NotFound(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Patch("/invoice/:invoice_id", controller.UpdateStatus)

	invoiceID := "INV-2024-999"

	updateReq := domain.UpdateInvoiceStatusRequest{
		Status: "sukses",
	}

	mockUsecase.On("UpdateStatus", invoiceID, &updateReq).Return(nil, errors.New("invoice tidak ditemukan"))

	reqBody, _ := json.Marshal(updateReq)
	req := httptest.NewRequest("PATCH", "/invoice/"+invoiceID, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "invoice tidak ditemukan", response["message"])
	assert.Equal(t, float64(404), response["code"])

	mockUsecase.AssertExpectations(t)
}

func TestInvoiceController_GetStatistics_Success(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice/statistics", controller.GetStatistics)

	mockStats := &domain.InvoiceStatistics{
		TotalInvoice:       100,
		TotalSukses:        85,
		TotalGagal:         10,
		TotalPending:       5,
		TotalPendapatan:    8500000,
		RataRataPembayaran: 100000,
		InvoiceBulanan: []domain.InvoiceStatsBulanan{
			{
				Bulan:           "Januari",
				TotalInvoice:    100,
				TotalSukses:     85,
				TotalPendapatan: 8500000,
			},
		},
		TopSubscriptionPlans: []domain.TopSubscriptionPlan{
			{
				SubscriptionPlanNama: "PRO Monthly",
				JumlahInvoice:        50,
				TotalPendapatan:      5000000,
			},
		},
	}

	mockUsecase.On("GetStatistics", mock.AnythingOfType("*domain.InvoiceStatisticsRequest")).Return(mockStats, nil)

	req := httptest.NewRequest("GET", "/invoice/statistics?bulan=1&tahun=2024", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.Equal(t, "Statistik invoice berhasil diambil", response["message"])
	assert.Equal(t, float64(200), response["code"])
	assert.NotNil(t, response["data"])

	mockUsecase.AssertExpectations(t)
}

func TestInvoiceController_GetStatistics_ValidationError(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice/statistics", controller.GetStatistics)

	req := httptest.NewRequest("GET", "/invoice/statistics?bulan=13&tahun=2019", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Data validasi tidak valid", response["message"])
	assert.Equal(t, float64(400), response["code"])

	mockUsecase.AssertNotCalled(t, "GetStatistics")
}

func TestInvoiceController_GetStatistics_UsecaseError(t *testing.T) {
	mockUsecase := new(MockInvoiceUsecase)
	controller := http.NewInvoiceController(mockUsecase)

	app := fiber.New()
	app.Get("/invoice/statistics", controller.GetStatistics)

	mockUsecase.On("GetStatistics", mock.AnythingOfType("*domain.InvoiceStatisticsRequest")).Return(nil, errors.New("database error"))

	req := httptest.NewRequest("GET", "/invoice/statistics", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, false, response["success"])
	assert.Equal(t, "Terjadi kesalahan pada server", response["message"])
	assert.Equal(t, float64(500), response["code"])

	mockUsecase.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
