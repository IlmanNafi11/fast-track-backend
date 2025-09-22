package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type KantongController struct {
	kantongUsecase usecase.KantongUsecase
}

func NewKantongController(kantongUsecase usecase.KantongUsecase) *KantongController {
	return &KantongController{
		kantongUsecase: kantongUsecase,
	}
}

func (c *KantongController) GetKantongList(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)

	req := domain.NewKantongListRequest()

	if search := ctx.Query("search"); search != "" {
		req.Search = &search
	}

	if sortBy := ctx.Query("sort_by"); sortBy != "" {
		req.SortBy = sortBy
	}

	if sortDirection := ctx.Query("sort_direction"); sortDirection != "" {
		req.SortDirection = sortDirection
	}

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	}

	if perPageStr := ctx.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage > 0 {
			req.PerPage = perPage
		}
	}

	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(ctx, validationErrors)
	}

	kantongs, meta, err := c.kantongUsecase.GetKantongList(userID, req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(ctx)
	}

	return helper.SendPaginatedResponse(ctx, 200, "Daftar kantong berhasil diambil", kantongs, *meta)
}

func (c *KantongController) GetKantongByID(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)
	id := ctx.Query("id")

	if id == "" {
		return helper.SendErrorResponse(ctx, 400, "ID kantong wajib diisi", nil)
	}

	kantong, err := c.kantongUsecase.GetKantongByID(id, userID)
	if err != nil {
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendNotFoundResponse(ctx, err.Error())
		}
		return helper.SendInternalServerErrorResponse(ctx)
	}

	return helper.SendSuccessResponse(ctx, 200, "Detail kantong berhasil diambil", kantong)
}

func (c *KantongController) CreateKantong(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)

	var req domain.CreateKantongRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(ctx, 400, "Format data tidak valid", err.Error())
	}

	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(ctx, validationErrors)
	}

	kantong, err := c.kantongUsecase.CreateKantong(&req, userID)
	if err != nil {
		if err.Error() == "nama kantong sudah ada" {
			return helper.SendErrorResponse(ctx, 409, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(ctx)
	}

	return helper.SendSuccessResponse(ctx, 201, "Kantong berhasil dibuat", kantong)
}

func (c *KantongController) UpdateKantong(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)
	id := ctx.Query("id")

	if id == "" {
		return helper.SendErrorResponse(ctx, 400, "ID kantong wajib diisi", nil)
	}

	var req domain.UpdateKantongRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(ctx, 400, "Format data tidak valid", err.Error())
	}

	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(ctx, validationErrors)
	}

	kantong, err := c.kantongUsecase.UpdateKantong(id, &req, userID)
	if err != nil {
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendNotFoundResponse(ctx, err.Error())
		}
		if err.Error() == "nama kantong sudah ada" {
			return helper.SendErrorResponse(ctx, 409, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(ctx)
	}

	return helper.SendSuccessResponse(ctx, 200, "Kantong berhasil diperbarui", kantong)
}

func (c *KantongController) PatchKantong(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)
	id := ctx.Query("id")

	if id == "" {
		return helper.SendErrorResponse(ctx, 400, "ID kantong wajib diisi", nil)
	}

	var req domain.PatchKantongRequest
	if err := ctx.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(ctx, 400, "Format data tidak valid", err.Error())
	}

	validationErrors := helper.ValidateStruct(req)
	if len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(ctx, validationErrors)
	}

	kantong, err := c.kantongUsecase.PatchKantong(id, &req, userID)
	if err != nil {
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendNotFoundResponse(ctx, err.Error())
		}
		if err.Error() == "nama kantong sudah ada" {
			return helper.SendErrorResponse(ctx, 409, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(ctx)
	}

	return helper.SendSuccessResponse(ctx, 200, "Kantong berhasil diperbarui", kantong)
}

func (c *KantongController) DeleteKantong(ctx *fiber.Ctx) error {
	userID := ctx.Locals("user_id").(uint)
	id := ctx.Query("id")

	if id == "" {
		return helper.SendErrorResponse(ctx, 400, "ID kantong wajib diisi", nil)
	}

	err := c.kantongUsecase.DeleteKantong(id, userID)
	if err != nil {
		if err.Error() == "kantong tidak ditemukan" {
			return helper.SendNotFoundResponse(ctx, err.Error())
		}
		if err.Error() == "kantong tidak dapat dihapus karena masih memiliki saldo" {
			return helper.SendErrorResponse(ctx, 409, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(ctx)
	}

	return helper.SendSuccessResponse(ctx, 200, "Kantong berhasil dihapus", nil)
}
