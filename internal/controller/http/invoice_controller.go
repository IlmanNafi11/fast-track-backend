package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type InvoiceController struct {
	invoiceUsecase usecase.InvoiceUsecase
}

func NewInvoiceController(invoiceUsecase usecase.InvoiceUsecase) *InvoiceController {
	return &InvoiceController{
		invoiceUsecase: invoiceUsecase,
	}
}

func (ctrl *InvoiceController) GetAll(c *fiber.Ctx) error {
	var req domain.InvoiceListRequest

	if err := c.QueryParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format query parameter tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	invoices, meta, err := ctrl.invoiceUsecase.GetAll(&req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendPaginatedResponse(c, fiber.StatusOK, "Daftar invoice berhasil diambil", invoices, *meta)
}

func (ctrl *InvoiceController) GetByID(c *fiber.Ctx) error {
	id := c.Params("invoice_id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID invoice wajib diisi", nil)
	}

	invoice, err := ctrl.invoiceUsecase.GetByID(id)
	if err != nil {
		switch err.Error() {
		case "invoice tidak ditemukan":
			return helper.SendNotFoundResponse(c, err.Error())
		default:
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Detail invoice berhasil diambil", invoice)
}

func (ctrl *InvoiceController) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("invoice_id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID invoice wajib diisi", nil)
	}

	var req domain.UpdateInvoiceStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format data tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	invoice, err := ctrl.invoiceUsecase.UpdateStatus(id, &req)
	if err != nil {
		switch err.Error() {
		case "invoice tidak ditemukan":
			return helper.SendNotFoundResponse(c, err.Error())
		default:
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Status invoice berhasil diperbarui", invoice)
}

func (ctrl *InvoiceController) GetStatistics(c *fiber.Ctx) error {
	var req domain.InvoiceStatisticsRequest

	if err := c.QueryParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format query parameter tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	stats, err := ctrl.invoiceUsecase.GetStatistics(&req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Statistik invoice berhasil diambil", stats)
}
