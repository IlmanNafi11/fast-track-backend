package config

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunSeeder(db *gorm.DB, cfg *Config) {
	helper.Info("Menjalankan database seeder...")

	if cfg.Database.SeedUsers {
		seedUsers(db)
	} else {
		helper.Info("Seeder users dinonaktifkan melalui konfigurasi")
	}

	helper.Info("Database seeder selesai")
}

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "user@example.com").Count(&count)

	if count > 0 {
		helper.Info("User seed sudah ada, melewati seeding user")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("user1234"), bcrypt.DefaultCost)
	if err != nil {
		helper.Fatal("Gagal hash password", err)
	}

	user := domain.User{
		Email:    "user@example.com",
		Password: string(hashedPassword),
		Name:     "user example",
		IsActive: true,
	}

	if err := db.Create(&user).Error; err != nil {
		helper.Fatal("Gagal membuat user seed", err, logrus.Fields{
			"email": user.Email,
		})
	}

	helper.Info("User seed berhasil dibuat", logrus.Fields{
		"email": user.Email,
	})
}
