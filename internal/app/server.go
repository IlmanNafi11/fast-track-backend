package app

import (
	"fiber-boiler-plate/config"
	"fiber-boiler-plate/internal/controller/http"
	"fiber-boiler-plate/internal/helper"
	"fiber-boiler-plate/internal/usecase"
	"fiber-boiler-plate/internal/usecase/repo"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func NewServer(cfg *config.Config, db *gorm.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return helper.SendInternalServerErrorResponse(c)
		},
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Refresh-Token",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	userRepo := repo.NewUserRepository(db)
	refreshTokenRepo := repo.NewRefreshTokenRepository(db)
	resetTokenRepo := repo.NewPasswordResetTokenRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, refreshTokenRepo, resetTokenRepo, cfg)
	authController := http.NewAuthController(authUsecase)

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/refresh", authController.RefreshToken)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Post("/reset-password/confirm", authController.ConfirmResetPassword)

	protected := auth.Group("/", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	protected.Post("logout", authController.Logout)

	app.Get("/health", func(c *fiber.Ctx) error {
		return helper.SendSuccessResponse(c, fiber.StatusOK, "Server berjalan dengan baik", map[string]string{
			"status": "healthy",
			"app":    cfg.App.Name,
		})
	})

	return app
}
