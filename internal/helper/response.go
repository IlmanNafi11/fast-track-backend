package helper

import (
	"fiber-boiler-plate/internal/domain"
	"time"

	"github.com/gofiber/fiber/v2"
)

func SendSuccessResponse(c *fiber.Ctx, code int, message string, data interface{}) error {
	response := domain.SuccessResponse{
		BaseResponse: domain.BaseResponse{
			Success: true,
			Message: message,
			Code:    code,
		},
		Data:      data,
		Timestamp: time.Now(),
	}
	return c.Status(code).JSON(response)
}

func SendErrorResponse(c *fiber.Ctx, code int, message string, errors interface{}) error {
	response := domain.ErrorResponse{
		BaseResponse: domain.BaseResponse{
			Success: false,
			Message: message,
			Code:    code,
		},
		Errors:    errors,
		Timestamp: time.Now(),
	}
	return c.Status(code).JSON(response)
}

func SendPaginatedResponse(c *fiber.Ctx, code int, message string, data interface{}, meta domain.PaginationMeta) error {
	response := domain.PaginatedResponse{
		SuccessResponse: domain.SuccessResponse{
			BaseResponse: domain.BaseResponse{
				Success: true,
				Message: message,
				Code:    code,
			},
			Data:      data,
			Timestamp: time.Now(),
		},
		Meta: meta,
	}
	return c.Status(code).JSON(response)
}

func SendValidationErrorResponse(c *fiber.Ctx, validationErrors []domain.ValidationError) error {
	return SendErrorResponse(c, fiber.StatusBadRequest, "Data validasi tidak valid", validationErrors)
}

func SendInternalServerErrorResponse(c *fiber.Ctx) error {
	return SendErrorResponse(c, fiber.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
}

func SendUnauthorizedResponse(c *fiber.Ctx) error {
	return SendErrorResponse(c, fiber.StatusUnauthorized, "Tidak memiliki akses", nil)
}

func SendForbiddenResponse(c *fiber.Ctx) error {
	return SendErrorResponse(c, fiber.StatusForbidden, "Akses ditolak", nil)
}

func SendNotFoundResponse(c *fiber.Ctx, message string) error {
	if message == "" {
		message = "Data tidak ditemukan"
	}
	return SendErrorResponse(c, fiber.StatusNotFound, message, nil)
}
