package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type SubscriptionPlanController struct {
	subscriptionPlanUsecase usecase.SubscriptionPlanUsecase
}

func NewSubscriptionPlanController(subscriptionPlanUsecase usecase.SubscriptionPlanUsecase) *SubscriptionPlanController {
	return &SubscriptionPlanController{
		subscriptionPlanUsecase: subscriptionPlanUsecase,
	}
}

func (ctrl *SubscriptionPlanController) GetAll(c *fiber.Ctx) error {
	var req domain.SubscriptionPlanListRequest

	if err := c.QueryParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format query parameter tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	plans, meta, err := ctrl.subscriptionPlanUsecase.GetAll(&req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendPaginatedResponse(c, fiber.StatusOK, "Daftar subscription plan berhasil diambil", plans, *meta)
}

func (ctrl *SubscriptionPlanController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID subscription plan wajib diisi", nil)
	}

	plan, err := ctrl.subscriptionPlanUsecase.GetByID(id)
	if err != nil {
		switch err.Error() {
		case "format ID tidak valid":
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		case "subscription plan tidak ditemukan":
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		default:
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Detail subscription plan berhasil diambil", plan)
}

func (ctrl *SubscriptionPlanController) Create(c *fiber.Ctx) error {
	var req domain.CreateSubscriptionPlanRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	plan, err := ctrl.subscriptionPlanUsecase.Create(&req)
	if err != nil {
		if err.Error() == "subscription plan dengan nama tersebut sudah ada" {
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusCreated, "Subscription plan berhasil dibuat", plan)
}

func (ctrl *SubscriptionPlanController) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID subscription plan wajib diisi", nil)
	}

	var req domain.UpdateSubscriptionPlanRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	plan, err := ctrl.subscriptionPlanUsecase.Update(id, &req)
	if err != nil {
		switch err.Error() {
		case "format ID tidak valid":
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		case "subscription plan tidak ditemukan":
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		case "subscription plan dengan nama tersebut sudah ada":
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Subscription plan berhasil diupdate", plan)
}

func (ctrl *SubscriptionPlanController) Patch(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID subscription plan wajib diisi", nil)
	}

	var req domain.PatchSubscriptionPlanRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	plan, err := ctrl.subscriptionPlanUsecase.Patch(id, &req)
	if err != nil {
		switch err.Error() {
		case "format ID tidak valid":
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		case "subscription plan tidak ditemukan":
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		case "subscription plan dengan nama tersebut sudah ada":
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Subscription plan berhasil diupdate", plan)
}

func (ctrl *SubscriptionPlanController) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID subscription plan wajib diisi", nil)
	}

	err := ctrl.subscriptionPlanUsecase.Delete(id)
	if err != nil {
		switch err.Error() {
		case "format ID tidak valid":
			return helper.SendErrorResponse(c, fiber.StatusBadRequest, err.Error(), nil)
		case "subscription plan tidak ditemukan":
			return helper.SendErrorResponse(c, fiber.StatusNotFound, err.Error(), nil)
		case "subscription plan tidak dapat dihapus karena sedang digunakan oleh pengguna aktif":
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Subscription plan berhasil dihapus", nil)
}
