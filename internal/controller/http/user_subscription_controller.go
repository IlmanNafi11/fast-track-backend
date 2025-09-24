package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type UserSubscriptionController struct {
	userSubscriptionUsecase usecase.UserSubscriptionUsecase
}

func NewUserSubscriptionController(userSubscriptionUsecase usecase.UserSubscriptionUsecase) *UserSubscriptionController {
	return &UserSubscriptionController{
		userSubscriptionUsecase: userSubscriptionUsecase,
	}
}

func (ctrl *UserSubscriptionController) GetAll(c *fiber.Ctx) error {
	var req domain.UserSubscriptionListRequest

	if err := c.QueryParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format query parameter tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	subscriptions, meta, err := ctrl.userSubscriptionUsecase.GetAll(&req)
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendPaginatedResponse(c, fiber.StatusOK, "Daftar subscription pengguna berhasil diambil", subscriptions, *meta)
}

func (ctrl *UserSubscriptionController) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID subscription pengguna wajib diisi", nil)
	}

	subscription, err := ctrl.userSubscriptionUsecase.GetByID(id)
	if err != nil {
		switch err.Error() {
		case "subscription pengguna tidak ditemukan":
			return helper.SendNotFoundResponse(c, err.Error())
		default:
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Detail subscription pengguna berhasil diambil", subscription)
}

func (ctrl *UserSubscriptionController) UpdateStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID subscription pengguna wajib diisi", nil)
	}

	var req domain.UpdateUserSubscriptionRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format data tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	subscription, err := ctrl.userSubscriptionUsecase.UpdateStatus(id, &req)
	if err != nil {
		switch err.Error() {
		case "subscription pengguna tidak ditemukan":
			return helper.SendNotFoundResponse(c, err.Error())
		case "subscription yang sudah dibatalkan atau berakhir tidak dapat diubah":
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		default:
			if err.Error() == "subscription sudah dalam status pause" ||
				err.Error() == "subscription sudah dalam status active" {
				return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
			}
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Status subscription pengguna berhasil diupdate", subscription)
}

func (ctrl *UserSubscriptionController) UpdatePaymentMethod(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "ID subscription pengguna wajib diisi", nil)
	}

	var req domain.UpdatePaymentMethodRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format data tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	subscription, err := ctrl.userSubscriptionUsecase.UpdatePaymentMethod(id, &req)
	if err != nil {
		switch err.Error() {
		case "subscription pengguna tidak ditemukan":
			return helper.SendNotFoundResponse(c, err.Error())
		case "subscription yang sudah dibatalkan atau berakhir tidak dapat diubah":
			return helper.SendErrorResponse(c, fiber.StatusConflict, err.Error(), nil)
		default:
			return helper.SendInternalServerErrorResponse(c)
		}
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Metode pembayaran berhasil diupdate", subscription)
}

func (ctrl *UserSubscriptionController) GetStatistics(c *fiber.Ctx) error {
	stats, err := ctrl.userSubscriptionUsecase.GetStatistics()
	if err != nil {
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Statistik subscription pengguna berhasil diambil", stats)
}
