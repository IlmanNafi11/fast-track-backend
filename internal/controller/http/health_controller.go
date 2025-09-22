package http

import (
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	healthUsecase usecase.HealthUsecase
}

func NewHealthController(healthUsecase usecase.HealthUsecase) *HealthController {
	return &HealthController{
		healthUsecase: healthUsecase,
	}
}

func (ctrl *HealthController) BasicHealthCheck(c *fiber.Ctx) error {
	healthData := ctrl.healthUsecase.GetBasicHealth()

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Server berjalan dengan baik", healthData)
}

func (ctrl *HealthController) ComprehensiveHealthCheck(c *fiber.Ctx) error {
	healthData := ctrl.healthUsecase.GetComprehensiveHealth()

	if healthData.Status == "unhealthy" {
		return helper.SendErrorResponse(c, fiber.StatusServiceUnavailable, "Beberapa komponen sistem mengalami masalah", healthData)
	}

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Pemeriksaan kesehatan sistem berhasil", healthData)
}

func (ctrl *HealthController) GetSystemMetrics(c *fiber.Ctx) error {
	metricsData := ctrl.healthUsecase.GetSystemMetrics()

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Metrics sistem berhasil diambil", metricsData)
}

func (ctrl *HealthController) GetApplicationStatus(c *fiber.Ctx) error {
	statusData := ctrl.healthUsecase.GetApplicationStatus()

	return helper.SendSuccessResponse(c, fiber.StatusOK, "Status aplikasi berhasil diambil", statusData)
}
