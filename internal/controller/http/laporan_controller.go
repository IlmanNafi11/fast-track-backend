package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LaporanController struct {
	laporanUsecase usecase.LaporanUsecase
}

func NewLaporanController(laporanUsecase usecase.LaporanUsecase) *LaporanController {
	return &LaporanController{
		laporanUsecase: laporanUsecase,
	}
}

func (ctrl *LaporanController) GetRingkasanLaporan(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	req := &domain.RingkasanLaporanRequest{}

	if tanggalMulai := c.Query("tanggal_mulai"); tanggalMulai != "" {
		req.TanggalMulai = &tanggalMulai
	}
	if tanggalSelesai := c.Query("tanggal_selesai"); tanggalSelesai != "" {
		req.TanggalSelesai = &tanggalSelesai
	}

	if err := helper.ValidateStruct(req); err != nil {
		return helper.SendValidationErrorResponse(c, err)
	}

	response, err := ctrl.laporanUsecase.GetRingkasanLaporan(userID, req)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (ctrl *LaporanController) GetStatistikTahunan(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	req := &domain.StatistikTahunanRequest{}

	if tahunStr := c.Query("tahun"); tahunStr != "" {
		tahun, err := strconv.Atoi(tahunStr)
		if err != nil {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format tahun tidak valid", nil)
		}
		req.Tahun = &tahun
	}

	if err := helper.ValidateStruct(req); err != nil {
		return helper.SendValidationErrorResponse(c, err)
	}

	response, err := ctrl.laporanUsecase.GetStatistikTahunan(userID, req)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (ctrl *LaporanController) GetStatistikKantongBulanan(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	req := &domain.StatistikKantongBulananRequest{}

	if bulanStr := c.Query("bulan"); bulanStr != "" {
		bulan, err := strconv.Atoi(bulanStr)
		if err != nil {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format bulan tidak valid", nil)
		}
		req.Bulan = &bulan
	}

	if tahunStr := c.Query("tahun"); tahunStr != "" {
		tahun, err := strconv.Atoi(tahunStr)
		if err != nil {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format tahun tidak valid", nil)
		}
		req.Tahun = &tahun
	}

	if err := helper.ValidateStruct(req); err != nil {
		return helper.SendValidationErrorResponse(c, err)
	}

	response, err := ctrl.laporanUsecase.GetStatistikKantongBulanan(userID, req)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func (ctrl *LaporanController) GetTopKantongPengeluaran(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusUnauthorized, "Token tidak valid", nil)
	}

	req := &domain.TopKantongPengeluaranRequest{}

	if bulanStr := c.Query("bulan"); bulanStr != "" {
		bulan, err := strconv.Atoi(bulanStr)
		if err != nil {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format bulan tidak valid", nil)
		}
		req.Bulan = &bulan
	}

	if tahunStr := c.Query("tahun"); tahunStr != "" {
		tahun, err := strconv.Atoi(tahunStr)
		if err != nil {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format tahun tidak valid", nil)
		}
		req.Tahun = &tahun
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format limit tidak valid", nil)
		}
		req.Limit = &limit
	}

	if err := helper.ValidateStruct(req); err != nil {
		return helper.SendValidationErrorResponse(c, err)
	}

	response, err := ctrl.laporanUsecase.GetTopKantongPengeluaran(userID, req)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
