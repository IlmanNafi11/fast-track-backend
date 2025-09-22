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
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewServer(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *fiber.App {
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
	redisRepo := repo.NewRedisRepository(rdb)
	kantongRepo := repo.NewKantongRepository(db, redisRepo)
	transaksiRepo := repo.NewTransaksiRepository(db)

	authUsecase := usecase.NewAuthUsecase(userRepo, refreshTokenRepo, resetTokenRepo, cfg)
	authController := http.NewAuthController(authUsecase)

	kantongUsecase := usecase.NewKantongUsecase(kantongRepo, userRepo)
	kantongController := http.NewKantongController(kantongUsecase)

	transaksiUsecase := usecase.NewTransaksiUsecase(transaksiRepo, kantongRepo, redisRepo)
	transaksiController := http.NewTransaksiController(transaksiUsecase)

	healthUsecase := usecase.NewHealthUsecase(db, rdb, cfg)
	healthController := http.NewHealthController(healthUsecase)

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", authController.Register)
	auth.Post("/login", authController.Login)
	auth.Post("/refresh", authController.RefreshToken)
	auth.Post("/reset-password", authController.ResetPassword)
	auth.Post("/reset-password/confirm", authController.ConfirmResetPassword)

	protected := auth.Group("/", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	protected.Post("logout", authController.Logout)

	kantong := api.Group("/kantong", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	kantong.Get("/", kantongController.GetKantongList)
	kantong.Get("/detail", kantongController.GetKantongByID)
	kantong.Post("/", kantongController.CreateKantong)
	kantong.Put("/", kantongController.UpdateKantong)
	kantong.Patch("/", kantongController.PatchKantong)
	kantong.Delete("/", kantongController.DeleteKantong)

	transaksi := api.Group("/transaksi", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	transaksi.Get("/", transaksiController.GetTransaksiList)
	transaksi.Get("/detail/:id", transaksiController.GetTransaksiDetail)
	transaksi.Post("/", transaksiController.CreateTransaksi)
	transaksi.Put("/update/:id", transaksiController.UpdateTransaksi)
	transaksi.Patch("/patch/:id", transaksiController.PatchTransaksi)
	transaksi.Delete("/delete/:id", transaksiController.DeleteTransaksi)

	monitoring := api.Group("/monitoring")
	monitoring.Get("/health", healthController.ComprehensiveHealthCheck)
	monitoring.Get("/metrics", healthController.GetSystemMetrics)
	monitoring.Get("/status", healthController.GetApplicationStatus)

	app.Get("/health", healthController.BasicHealthCheck)

	return app
}
