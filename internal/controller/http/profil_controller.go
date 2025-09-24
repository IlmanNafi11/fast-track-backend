package http

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type ProfilController struct {
	profilUsecase usecase.ProfilUsecase
}

func NewProfilController(profilUsecase usecase.ProfilUsecase) *ProfilController {
	return &ProfilController{
		profilUsecase: profilUsecase,
	}
}

func (ctrl *ProfilController) GetProfil(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendUnauthorizedResponse(c)
	}

	profil, err := ctrl.profilUsecase.GetProfil(userID)
	if err != nil {
		if err.Error() == "profil tidak ditemukan" {
			return helper.SendNotFoundResponse(c, err.Error())
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Data profil berhasil diambil", profil)
}

func (ctrl *ProfilController) UpdateProfil(c *fiber.Ctx) error {
	userID, err := helper.GetUserIDFromToken(c)
	if err != nil {
		return helper.SendUnauthorizedResponse(c)
	}

	var req domain.UpdateProfilRequest
	if err := c.BodyParser(&req); err != nil {
		return helper.SendErrorResponse(c, fiber.StatusBadRequest, "Format request tidak valid", nil)
	}

	if validationErrors := helper.ValidateStruct(req); len(validationErrors) > 0 {
		return helper.SendValidationErrorResponse(c, validationErrors)
	}

	profil, err := ctrl.profilUsecase.UpdateProfil(userID, req)
	if err != nil {
		if err.Error() == "profil tidak ditemukan" {
			return helper.SendNotFoundResponse(c, err.Error())
		}
		return helper.SendInternalServerErrorResponse(c)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Profil berhasil diperbarui", profil)
}
