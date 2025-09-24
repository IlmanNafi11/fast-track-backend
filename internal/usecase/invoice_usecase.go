package usecase

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"fmt"

	"gorm.io/gorm"
)

type InvoiceUsecase interface {
	GetAll(req *domain.InvoiceListRequest) ([]*domain.InvoiceResponse, *domain.PaginationMeta, error)
	GetByID(id string) (*domain.InvoiceDetailResponse, error)
	UpdateStatus(id string, req *domain.UpdateInvoiceStatusRequest) (*domain.InvoiceDetailResponse, error)
	GetStatistics(req *domain.InvoiceStatisticsRequest) (*domain.InvoiceStatistics, error)
}

type invoiceUsecase struct {
	invoiceRepo repo.InvoiceRepository
}

func NewInvoiceUsecase(invoiceRepo repo.InvoiceRepository) InvoiceUsecase {
	return &invoiceUsecase{
		invoiceRepo: invoiceRepo,
	}
}

func (uc *invoiceUsecase) GetAll(req *domain.InvoiceListRequest) ([]*domain.InvoiceResponse, *domain.PaginationMeta, error) {
	req.SetDefaults()

	invoices, total, err := uc.invoiceRepo.GetAll(req)
	if err != nil {
		return nil, nil, fmt.Errorf("gagal mengambil daftar invoice: %w", err)
	}

	totalPages := (total + req.PerPage - 1) / req.PerPage
	meta := &domain.PaginationMeta{
		CurrentPage:  req.Page,
		TotalPages:   totalPages,
		TotalRecords: total,
		PerPage:      req.PerPage,
	}

	response := make([]*domain.InvoiceResponse, len(invoices))
	for i, invoice := range invoices {
		response[i] = uc.mapToResponse(invoice)
	}

	return response, meta, nil
}

func (uc *invoiceUsecase) GetByID(id string) (*domain.InvoiceDetailResponse, error) {
	invoice, err := uc.invoiceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil detail invoice: %w", err)
	}

	return uc.mapToDetailResponse(invoice), nil
}

func (uc *invoiceUsecase) UpdateStatus(id string, req *domain.UpdateInvoiceStatusRequest) (*domain.InvoiceDetailResponse, error) {
	_, err := uc.invoiceRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invoice tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal mengambil data invoice: %w", err)
	}

	if err := uc.invoiceRepo.UpdateStatus(id, req.Status, req.Keterangan); err != nil {
		return nil, fmt.Errorf("gagal memperbarui status invoice: %w", err)
	}

	updatedInvoice, err := uc.invoiceRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data invoice yang diperbarui: %w", err)
	}

	return uc.mapToDetailResponse(updatedInvoice), nil
}

func (uc *invoiceUsecase) GetStatistics(req *domain.InvoiceStatisticsRequest) (*domain.InvoiceStatistics, error) {
	stats, err := uc.invoiceRepo.GetStatistics(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil statistik invoice: %w", err)
	}

	return stats, nil
}

func (uc *invoiceUsecase) mapToResponse(invoice *domain.Invoice) *domain.InvoiceResponse {
	response := &domain.InvoiceResponse{
		InvoiceID:          invoice.ID,
		UserID:             invoice.UserID,
		Jumlah:             invoice.Jumlah,
		Status:             invoice.Status,
		DibayarPada:        invoice.DibayarPada,
		MetodePembayaran:   invoice.MetodePembayaran,
		Keterangan:         invoice.Keterangan,
		SubscriptionPlanID: invoice.SubscriptionPlanID,
		CreatedAt:          invoice.CreatedAt,
		UpdatedAt:          invoice.UpdatedAt,
	}

	if invoice.User.ID != 0 {
		response.NamaUser = invoice.User.Name
		response.UserEmail = invoice.User.Email
	}

	if invoice.SubscriptionPlan != nil && invoice.SubscriptionPlan.ID != (domain.SubscriptionPlan{}).ID {
		response.SubscriptionPlanNama = &invoice.SubscriptionPlan.Nama
	}

	return response
}

func (uc *invoiceUsecase) mapToDetailResponse(invoice *domain.Invoice) *domain.InvoiceDetailResponse {
	response := &domain.InvoiceDetailResponse{
		InvoiceID:          invoice.ID,
		UserID:             invoice.UserID,
		Jumlah:             invoice.Jumlah,
		Status:             invoice.Status,
		DibayarPada:        invoice.DibayarPada,
		MetodePembayaran:   invoice.MetodePembayaran,
		Keterangan:         invoice.Keterangan,
		SubscriptionPlanID: invoice.SubscriptionPlanID,
		CreatedAt:          invoice.CreatedAt,
		UpdatedAt:          invoice.UpdatedAt,
	}

	if invoice.User.ID != 0 {
		response.NamaUser = invoice.User.Name
		response.UserEmail = invoice.User.Email
	}

	if invoice.SubscriptionPlan != nil && invoice.SubscriptionPlan.ID != (domain.SubscriptionPlan{}).ID {
		response.SubscriptionPlanNama = &invoice.SubscriptionPlan.Nama
	}

	return response
}
