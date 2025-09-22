package http

import (
	"strconv"

	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type AnggaranController struct {
	anggaranUsecase usecase.AnggaranUsecase
}

func NewAnggaranController(anggaranUsecase usecase.AnggaranUsecase) *AnggaranController {
	return &AnggaranController{
		anggaranUsecase: anggaranUsecase,
	}
}

func (ctrl *AnggaranController) GetAnggaranList(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	req := domain.NewAnggaranListRequest()

	if search := c.Query("search"); search != "" {
		req.Search = &search
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		req.SortBy = sortBy
	}

	if sortDirection := c.Query("sort_direction"); sortDirection != "" {
		req.SortDirection = sortDirection
	}

	if page := c.Query("page"); page != "" {
		if pageInt, err := strconv.Atoi(page); err == nil && pageInt > 0 {
			req.Page = pageInt
		}
	}

	if perPage := c.Query("per_page"); perPage != "" {
		if perPageInt, err := strconv.Atoi(perPage); err == nil && perPageInt > 0 && perPageInt <= 100 {
			req.PerPage = perPageInt
		}
	}

	if bulan := c.Query("bulan"); bulan != "" {
		if bulanInt, err := strconv.Atoi(bulan); err == nil && bulanInt >= 1 && bulanInt <= 12 {
			req.Bulan = &bulanInt
		}
	}

	if tahun := c.Query("tahun"); tahun != "" {
		if tahunInt, err := strconv.Atoi(tahun); err == nil && tahunInt >= 2020 {
			req.Tahun = &tahunInt
		}
	}

	if validationErrors := helper.ValidateStruct(*req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	responses, meta, err := ctrl.anggaranUsecase.GetAnggaranList(userID, req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendPaginatedResponse(c, fiber.StatusOK, "Daftar anggaran berhasil diambil", responses, *meta)
}

func (ctrl *AnggaranController) GetAnggaranDetail(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	kantongID := c.Params("kantong_id")

	if kantongID == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID kantong tidak valid", nil)
	}

	var bulan, tahun *int

	if bulanStr := c.Query("bulan"); bulanStr != "" {
		if bulanInt, err := strconv.Atoi(bulanStr); err == nil && bulanInt >= 1 && bulanInt <= 12 {
			bulan = &bulanInt
		} else {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Parameter bulan tidak valid", nil)
		}
	}

	if tahunStr := c.Query("tahun"); tahunStr != "" {
		if tahunInt, err := strconv.Atoi(tahunStr); err == nil && tahunInt >= 2020 {
			tahun = &tahunInt
		} else {
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Parameter tahun tidak valid", nil)
		}
	}

	response, err := ctrl.anggaranUsecase.GetAnggaranDetail(kantongID, userID, bulan, tahun)
	if err != nil {
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendNotFoundResponse(c, "Kantong tidak ditemukan")
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Detail anggaran berhasil diambil", response)
}

func (ctrl *AnggaranController) CreatePenyesuaianAnggaran(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req domain.PenyesuaianAnggaranRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	response, err := ctrl.anggaranUsecase.CreatePenyesuaianAnggaran(userID, &req)
	if err != nil {
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendNotFoundResponse(c, "Kantong tidak ditemukan")
		}
		if err.Error() == "tidak memiliki akses ke kantong ini" {
			return helper.SendForbiddenResponse(c)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Penyesuaian anggaran berhasil dibuat", response)
}
