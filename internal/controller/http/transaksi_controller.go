package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TransaksiController struct {
	transaksiUsecase usecase.TransaksiUsecase
}

func NewTransaksiController(transaksiUsecase usecase.TransaksiUsecase) *TransaksiController {
	return &TransaksiController{
		transaksiUsecase: transaksiUsecase,
	}
}

func (ctrl *TransaksiController) GetTransaksiList(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	req := &domain.TransaksiListRequest{
		Page:          1,
		PerPage:       10,
		SortBy:        "tanggal",
		SortDirection: "desc",
	}

	if search := c.Query("search"); search != "" {
		req.Search = &search
	}
	if jenis := c.Query("jenis"); jenis != "" {
		req.Jenis = &jenis
	}
	if kantongNama := c.Query("kantong_nama"); kantongNama != "" {
		req.KantongNama = &kantongNama
	}
	if tanggalMulai := c.Query("tanggal_mulai"); tanggalMulai != "" {
		req.TanggalMulai = &tanggalMulai
	}
	if tanggalSelesai := c.Query("tanggal_selesai"); tanggalSelesai != "" {
		req.TanggalSelesai = &tanggalSelesai
	}
	if sortBy := c.Query("sort_by"); sortBy != "" {
		req.SortBy = sortBy
	}
	if sortDirection := c.Query("sort_direction"); sortDirection != "" {
		req.SortDirection = sortDirection
	}

	if page, err := strconv.Atoi(c.Query("page", "1")); err == nil && page > 0 {
		req.Page = page
	}
	if perPage, err := strconv.Atoi(c.Query("per_page", "10")); err == nil && perPage > 0 && perPage <= 100 {
		req.PerPage = perPage
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	result, err := ctrl.transaksiUsecase.GetTransaksiList(userID, req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return c.Status(result.Code).JSON(result)
}

func (ctrl *TransaksiController) GetTransaksiDetail(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID transaksi diperlukan", nil)
	}

	result, err := ctrl.transaksiUsecase.GetTransaksiDetail(id, userID)
	if err != nil {
		if err.Error() == "transaksi tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return c.Status(result.Code).JSON(result)
}

func (ctrl *TransaksiController) CreateTransaksi(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	var req domain.CreateTransaksiRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	result, err := ctrl.transaksiUsecase.CreateTransaksi(userID, &req)
	if err != nil {
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "format tanggal tidak valid" {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		if err.Error() == "saldo tidak mencukupi" {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

func (ctrl *TransaksiController) UpdateTransaksi(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID transaksi diperlukan", nil)
	}

	var req domain.UpdateTransaksiRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	result, err := ctrl.transaksiUsecase.UpdateTransaksi(id, userID, &req)
	if err != nil {
		if err.Error() == "transaksi tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "format tanggal tidak valid" {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		if err.Error() == "saldo tidak mencukupi" {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return c.Status(result.Code).JSON(result)
}

func (ctrl *TransaksiController) PatchTransaksi(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID transaksi diperlukan", nil)
	}

	var req domain.PatchTransaksiRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	result, err := ctrl.transaksiUsecase.PatchTransaksi(id, userID, &req)
	if err != nil {
		if err.Error() == "transaksi tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		if err.Error() == "format tanggal tidak valid" {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		if err.Error() == "saldo tidak mencukupi" {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return c.Status(result.Code).JSON(result)
}

func (ctrl *TransaksiController) DeleteTransaksi(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID transaksi diperlukan", nil)
	}

	err = ctrl.transaksiUsecase.DeleteTransaksi(id, userID)
	if err != nil {
		if err.Error() == "transaksi tidak ditemukan" {
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Transaksi berhasil dihapus", nil)
}
