package config

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	helper.Info("Menjalankan database migration...")

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.RefreshToken{},
		&domain.PasswordResetToken{},
		&domain.Kantong{},
		&domain.Permission{},
		&domain.Role{},
		&domain.RolePermission{},
	); err != nil {
		helper.Fatal("Gagal melakukan auto migrate", err)
	}

	helper.Info("Database migration selesai")
}
