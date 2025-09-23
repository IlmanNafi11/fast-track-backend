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
	anggaranRepo := repo.NewAnggaranRepository(db, redisRepo)
	laporanRepo := repo.NewLaporanRepository(db)
	subscriptionPlanRepo := repo.NewSubscriptionPlanRepository(db, redisRepo)
	permissionRepo := repo.NewPermissionRepository(db, redisRepo)
	roleRepo := repo.NewRoleRepository(db, redisRepo)

	authUsecase := usecase.NewAuthUsecase(userRepo, refreshTokenRepo, resetTokenRepo, cfg)
	authController := http.NewAuthController(authUsecase)

	kantongUsecase := usecase.NewKantongUsecase(kantongRepo, userRepo)
	kantongController := http.NewKantongController(kantongUsecase)

	transaksiUsecase := usecase.NewTransaksiUsecase(transaksiRepo, kantongRepo, redisRepo)
	transaksiController := http.NewTransaksiController(transaksiUsecase)

	anggaranUsecase := usecase.NewAnggaranUsecase(anggaranRepo, kantongRepo, transaksiRepo, redisRepo)
	anggaranController := http.NewAnggaranController(anggaranUsecase)

	laporanUsecase := usecase.NewLaporanUsecase(laporanRepo, redisRepo)
	laporanController := http.NewLaporanController(laporanUsecase)

	subscriptionPlanUsecase := usecase.NewSubscriptionPlanUsecase(subscriptionPlanRepo)
	subscriptionPlanController := http.NewSubscriptionPlanController(subscriptionPlanUsecase)

	permissionUsecase := usecase.NewPermissionUsecase(permissionRepo)
	permissionController := http.NewPermissionController(permissionUsecase)

	roleUsecase := usecase.NewRoleUsecase(roleRepo, permissionRepo)
	roleController := http.NewRoleController(roleUsecase)

	kantongUsecase.SetAnggaranUsecase(anggaranUsecase)
	transaksiUsecase.SetAnggaranUsecase(anggaranUsecase)

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
	kantong.Get("/:id", kantongController.GetKantongByID)
	kantong.Post("/", kantongController.CreateKantong)
	kantong.Put("/:id", kantongController.UpdateKantong)
	kantong.Patch("/:id", kantongController.PatchKantong)
	kantong.Delete("/:id", kantongController.DeleteKantong)
	kantong.Post("/transfer", kantongController.TransferKantong)

	transaksi := api.Group("/transaksi", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	transaksi.Get("/", transaksiController.GetTransaksiList)
	transaksi.Get("/:id", transaksiController.GetTransaksiDetail)
	transaksi.Post("/", transaksiController.CreateTransaksi)
	transaksi.Put("/:id", transaksiController.UpdateTransaksi)
	transaksi.Patch("/:id", transaksiController.PatchTransaksi)
	transaksi.Delete("/:id", transaksiController.DeleteTransaksi)

	anggaran := api.Group("/anggaran", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	anggaran.Get("/", anggaranController.GetAnggaranList)
	anggaran.Get("/:kantong_id", anggaranController.GetAnggaranDetail)
	anggaran.Post("/penyesuaian", anggaranController.CreatePenyesuaianAnggaran)

	laporan := api.Group("/laporan", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	laporan.Get("/ringkasan", laporanController.GetRingkasanLaporan)
	laporan.Get("/statistik/tahunan", laporanController.GetStatistikTahunan)
	laporan.Get("/statistik/kantong-bulanan", laporanController.GetStatistikKantongBulanan)
	laporan.Get("/statistik/top-kantong", laporanController.GetTopKantongPengeluaran)
	laporan.Get("/statistik/kantong-periode", laporanController.GetStatistikKantongPeriode)
	laporan.Get("/pengeluaran-kantong-detail", laporanController.GetPengeluaranKantongDetail)
	laporan.Get("/tren/bulanan", laporanController.GetTrenBulanan)
	laporan.Get("/perbandingan/kantong", laporanController.GetPerbandinganKantong)
	laporan.Get("/perbandingan/kantong/detail", laporanController.GetDetailPerbandinganKantong)

	subscriptionPlan := api.Group("/subscription-plans", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	subscriptionPlan.Get("/", subscriptionPlanController.GetAll)
	subscriptionPlan.Get("/:id", subscriptionPlanController.GetByID)
	subscriptionPlan.Post("/", subscriptionPlanController.Create)
	subscriptionPlan.Put("/:id", subscriptionPlanController.Update)
	subscriptionPlan.Patch("/:id", subscriptionPlanController.Patch)
	subscriptionPlan.Delete("/:id", subscriptionPlanController.Delete)

	permission := api.Group("/permission", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	permission.Get("/", permissionController.GetPermissionList)
	permission.Get("/:id", permissionController.GetPermissionByID)
	permission.Post("/", permissionController.CreatePermission)
	permission.Put("/:id", permissionController.UpdatePermission)
	permission.Delete("/:id", permissionController.DeletePermission)

	role := api.Group("/role", helper.JWTAuthMiddleware(cfg.JWT.Secret))
	role.Get("/", roleController.GetRoleList)
	role.Get("/:id", roleController.GetRoleByID)
	role.Post("/", roleController.CreateRole)
	role.Put("/:id", roleController.UpdateRole)
	role.Delete("/:id", roleController.DeleteRole)
	role.Get("/:id/permissions", roleController.GetRolePermissions)

	monitoring := api.Group("/monitoring")
	monitoring.Get("/health", healthController.ComprehensiveHealthCheck)
	monitoring.Get("/metrics", healthController.GetSystemMetrics)
	monitoring.Get("/status", healthController.GetApplicationStatus)

	app.Get("/health", healthController.BasicHealthCheck)

	return app
}
