package main

import (
	"fiber-boiler-plate/config"
	"fiber-boiler-plate/internal/app"
	"fiber-boiler-plate/internal/helper"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.LoadConfig()

	helper.InitLogger(cfg.App.Env)

	db := config.ConnectDatabase(cfg)
	rdb := config.ConnectRedis(cfg)

	if cfg.Database.MigrateOnStart && cfg.Database.AutoMigrate {
		helper.Info("Menjalankan auto migration untuk environment", logrus.Fields{
			"environment": cfg.App.Env,
		})
		config.RunMigration(db)
	} else {
		helper.Info("Auto migration dinonaktifkan melalui konfigurasi")
	}

	if cfg.Database.RunSeeder {
		if cfg.App.Env == "production" {
			helper.Warn("PERINGATAN: Seeder tidak direkomendasikan untuk production environment!")
			helper.Warn("Melewati eksekusi seeder untuk keamanan production")
			helper.Info("Untuk menjalankan seeder di production, ubah APP_ENV ke nilai lain")
		} else {
			helper.Info("Menjalankan seeder untuk environment", logrus.Fields{
				"environment": cfg.App.Env,
			})
			config.RunSeeder(db, cfg)
		}
	} else {
		helper.Info("Seeder dinonaktifkan melalui konfigurasi DB_RUN_SEEDER=false")
	}

	server := app.NewServer(cfg, db, rdb)

	helper.Info("Server berjalan di port", logrus.Fields{
		"port": cfg.App.Port,
	})

	if err := server.Listen(":" + cfg.App.Port); err != nil {
		helper.Fatal("Gagal menjalankan server", err, logrus.Fields{
			"port": cfg.App.Port,
		})
	}
}
