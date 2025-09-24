package usecase_test

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockInvoiceRepository struct {
	mock.Mock
}

func (m *MockInvoiceRepository) GetAll(req *domain.InvoiceListRequest) ([]*domain.Invoice, int, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Invoice), args.Int(1), args.Error(2)
}

func (m *MockInvoiceRepository) GetByID(id string) (*domain.Invoice, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Invoice), args.Error(1)
}

func (m *MockInvoiceRepository) UpdateStatus(id string, status string, keterangan *string) error {
	args := m.Called(id, status, keterangan)
	return args.Error(0)
}

func (m *MockInvoiceRepository) GetStatistics(req *domain.InvoiceStatisticsRequest) (*domain.InvoiceStatistics, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.InvoiceStatistics), args.Error(1)
}

func TestNewInvoiceUsecase(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	assert.NotNil(t, invoiceUsecase)
}

func TestInvoiceUsecase_GetAll_Success(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	now := time.Now()
	userID := uint(1)
	planID := uuid.New()

	mockInvoices := []*domain.Invoice{
		{
			ID:                 "INV-2024-001",
			UserID:             userID,
			User:               domain.User{ID: userID, Name: "John Doe", Email: "john@example.com"},
			Jumlah:             99000,
			Status:             "sukses",
			DibayarPada:        &now,
			MetodePembayaran:   stringPtr("Bank Transfer"),
			Keterangan:         stringPtr("Pembayaran PRO"),
			SubscriptionPlanID: &planID,
			SubscriptionPlan:   &domain.SubscriptionPlan{ID: planID, Nama: "PRO Monthly"},
			CreatedAt:          now,
			UpdatedAt:          now,
		},
	}

	req := &domain.InvoiceListRequest{
		Page:    1,
		PerPage: 10,
	}

	mockRepo.On("GetAll", req).Return(mockInvoices, 1, nil)

	result, meta, err := invoiceUsecase.GetAll(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, meta)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, "INV-2024-001", result[0].InvoiceID)
	assert.Equal(t, "John Doe", result[0].NamaUser)
	assert.Equal(t, "john@example.com", result[0].UserEmail)
	assert.Equal(t, 1, meta.CurrentPage)
	assert.Equal(t, 1, meta.TotalPages)
	assert.Equal(t, 1, meta.TotalRecords)
	assert.Equal(t, 10, meta.PerPage)

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_GetAll_RepositoryError(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	req := &domain.InvoiceListRequest{
		Page:    1,
		PerPage: 10,
	}

	mockRepo.On("GetAll", req).Return(nil, 0, errors.New("database error"))

	result, meta, err := invoiceUsecase.GetAll(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Nil(t, meta)
	assert.Contains(t, err.Error(), "gagal mengambil daftar invoice")

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_GetByID_Success(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	now := time.Now()
	userID := uint(1)
	planID := uuid.New()
	invoiceID := "INV-2024-001"

	mockInvoice := &domain.Invoice{
		ID:                 invoiceID,
		UserID:             userID,
		User:               domain.User{ID: userID, Name: "John Doe", Email: "john@example.com"},
		Jumlah:             99000,
		Status:             "sukses",
		DibayarPada:        &now,
		MetodePembayaran:   stringPtr("Bank Transfer"),
		Keterangan:         stringPtr("Pembayaran PRO"),
		SubscriptionPlanID: &planID,
		SubscriptionPlan:   &domain.SubscriptionPlan{ID: planID, Nama: "PRO Monthly"},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	mockRepo.On("GetByID", invoiceID).Return(mockInvoice, nil)

	result, err := invoiceUsecase.GetByID(invoiceID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, invoiceID, result.InvoiceID)
	assert.Equal(t, "John Doe", result.NamaUser)
	assert.Equal(t, "john@example.com", result.UserEmail)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, float64(99000), result.Jumlah)
	assert.Equal(t, "sukses", result.Status)
	assert.Equal(t, "PRO Monthly", *result.SubscriptionPlanNama)

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_GetByID_NotFound(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	invoiceID := "INV-2024-999"

	mockRepo.On("GetByID", invoiceID).Return(nil, gorm.ErrRecordNotFound)

	result, err := invoiceUsecase.GetByID(invoiceID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invoice tidak ditemukan", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_GetByID_RepositoryError(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	invoiceID := "INV-2024-001"

	mockRepo.On("GetByID", invoiceID).Return(nil, errors.New("database error"))

	result, err := invoiceUsecase.GetByID(invoiceID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "gagal mengambil detail invoice")

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_UpdateStatus_Success(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	now := time.Now()
	userID := uint(1)
	planID := uuid.New()
	invoiceID := "INV-2024-001"
	keterangan := "Pembayaran telah dikonfirmasi"

	mockInvoice := &domain.Invoice{
		ID:                 invoiceID,
		UserID:             userID,
		User:               domain.User{ID: userID, Name: "John Doe", Email: "john@example.com"},
		Jumlah:             99000,
		Status:             "pending",
		DibayarPada:        nil,
		MetodePembayaran:   stringPtr("Bank Transfer"),
		Keterangan:         stringPtr("Pembayaran PRO"),
		SubscriptionPlanID: &planID,
		SubscriptionPlan:   &domain.SubscriptionPlan{ID: planID, Nama: "PRO Monthly"},
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	updatedInvoice := &domain.Invoice{
		ID:                 invoiceID,
		UserID:             userID,
		User:               domain.User{ID: userID, Name: "John Doe", Email: "john@example.com"},
		Jumlah:             99000,
		Status:             "sukses",
		DibayarPada:        &now,
		MetodePembayaran:   stringPtr("Bank Transfer"),
		Keterangan:         &keterangan,
		SubscriptionPlanID: &planID,
		SubscriptionPlan:   &domain.SubscriptionPlan{ID: planID, Nama: "PRO Monthly"},
		CreatedAt:          now,
		UpdatedAt:          time.Now(),
	}

	req := &domain.UpdateInvoiceStatusRequest{
		Status:     "sukses",
		Keterangan: &keterangan,
	}

	mockRepo.On("GetByID", invoiceID).Return(mockInvoice, nil).Once()
	mockRepo.On("UpdateStatus", invoiceID, req.Status, req.Keterangan).Return(nil)
	mockRepo.On("GetByID", invoiceID).Return(updatedInvoice, nil).Once()

	result, err := invoiceUsecase.UpdateStatus(invoiceID, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, invoiceID, result.InvoiceID)
	assert.Equal(t, "sukses", result.Status)
	assert.Equal(t, keterangan, *result.Keterangan)

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_UpdateStatus_InvoiceNotFound(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	invoiceID := "INV-2024-999"
	req := &domain.UpdateInvoiceStatusRequest{
		Status: "sukses",
	}

	mockRepo.On("GetByID", invoiceID).Return(nil, gorm.ErrRecordNotFound)

	result, err := invoiceUsecase.UpdateStatus(invoiceID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invoice tidak ditemukan", err.Error())

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_UpdateStatus_UpdateError(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	now := time.Now()
	userID := uint(1)
	invoiceID := "INV-2024-001"

	mockInvoice := &domain.Invoice{
		ID:        invoiceID,
		UserID:    userID,
		User:      domain.User{ID: userID, Name: "John Doe", Email: "john@example.com"},
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	req := &domain.UpdateInvoiceStatusRequest{
		Status: "sukses",
	}

	mockRepo.On("GetByID", invoiceID).Return(mockInvoice, nil)
	mockRepo.On("UpdateStatus", invoiceID, req.Status, req.Keterangan).Return(errors.New("update error"))

	result, err := invoiceUsecase.UpdateStatus(invoiceID, req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "gagal memperbarui status invoice")

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_GetStatistics_Success(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	bulan := 1
	tahun := 2024
	req := &domain.InvoiceStatisticsRequest{
		Bulan: &bulan,
		Tahun: &tahun,
	}

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

	mockRepo.On("GetStatistics", req).Return(mockStats, nil)

	result, err := invoiceUsecase.GetStatistics(req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(100), result.TotalInvoice)
	assert.Equal(t, int64(85), result.TotalSukses)
	assert.Equal(t, int64(10), result.TotalGagal)
	assert.Equal(t, int64(5), result.TotalPending)
	assert.Equal(t, float64(8500000), result.TotalPendapatan)
	assert.Equal(t, float64(100000), result.RataRataPembayaran)
	assert.Equal(t, 1, len(result.InvoiceBulanan))
	assert.Equal(t, 1, len(result.TopSubscriptionPlans))

	mockRepo.AssertExpectations(t)
}

func TestInvoiceUsecase_GetStatistics_RepositoryError(t *testing.T) {
	mockRepo := new(MockInvoiceRepository)
	invoiceUsecase := usecase.NewInvoiceUsecase(mockRepo)

	req := &domain.InvoiceStatisticsRequest{}

	mockRepo.On("GetStatistics", req).Return(nil, errors.New("database error"))

	result, err := invoiceUsecase.GetStatistics(req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "gagal mengambil statistik invoice")

	mockRepo.AssertExpectations(t)
}

func stringPtr(s string) *string {
	return &s
}
