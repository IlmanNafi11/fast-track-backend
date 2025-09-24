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
	userID := c.Locals("user_id").(uint)

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

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}

func (ctrl *LaporanController) GetStatistikTahunan(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

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

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}

func (ctrl *LaporanController) GetStatistikKantongBulanan(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

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

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}

func (ctrl *LaporanController) GetTopKantongPengeluaran(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

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

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}

func (ctrl *LaporanController) GetStatistikKantongPeriode(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	req := &domain.StatistikKantongPeriodeRequest{}

	if tanggalMulai := c.Query("tanggal_mulai"); tanggalMulai != "" {
		req.TanggalMulai = &tanggalMulai
	}
	if tanggalSelesai := c.Query("tanggal_selesai"); tanggalSelesai != "" {
		req.TanggalSelesai = &tanggalSelesai
	}

	if err := helper.ValidateStruct(req); err != nil {
		return helper.SendValidationErrorResponse(c, err)
	}

	response, err := ctrl.laporanUsecase.GetStatistikKantongPeriode(userID, req)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}

func (ctrl *LaporanController) GetPengeluaranKantongDetail(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	req := &domain.PengeluaranKantongDetailRequest{}

	if tanggalMulai := c.Query("tanggal_mulai"); tanggalMulai != "" {
		req.TanggalMulai = &tanggalMulai
	}
	if tanggalSelesai := c.Query("tanggal_selesai"); tanggalSelesai != "" {
		req.TanggalSelesai = &tanggalSelesai
	}

	if err := helper.ValidateStruct(req); err != nil {
		return helper.SendValidationErrorResponse(c, err)
	}

	response, err := ctrl.laporanUsecase.GetPengeluaranKantongDetail(userID, req)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}

func (ctrl *LaporanController) GetTrenBulanan(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	req := &domain.TrenBulananRequest{}

	if tahun := c.Query("tahun"); tahun != "" {
		if tahunInt, err := strconv.Atoi(tahun); err == nil {
			req.Tahun = &tahunInt
		} else {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format tahun tidak valid", nil)
		}
	}

	if err := helper.ValidateStruct(req); err != nil {
		return helper.SendValidationErrorResponse(c, err)
	}

	response, err := ctrl.laporanUsecase.GetTrenBulanan(userID, req)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}

func (ctrl *LaporanController) GetPerbandinganKantong(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	response, err := ctrl.laporanUsecase.GetPerbandinganKantong(userID)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}

func (ctrl *LaporanController) GetDetailPerbandinganKantong(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	response, err := ctrl.laporanUsecase.GetDetailPerbandinganKantong(userID)
	if err != nil {
		return helper.SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", err)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, response.Message, response.Data)
}
