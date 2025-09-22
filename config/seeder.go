package config

import (
	"fiber-boiler-plate/internal/domain"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RunSeeder(db *gorm.DB, cfg *Config) {
	log.Println("Menjalankan database seeder...")

	if cfg.Database.SeedUsers {
		seedUsers(db)
	} else {
		log.Println("Seeder users dinonaktifkan melalui konfigurasi")
	}

	log.Println("Database seeder selesai")
}

func seedUsers(db *gorm.DB) {
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "user@example.com").Count(&count)

	if count > 0 {
		log.Println("User seed sudah ada, melewati seeding user")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("user1234"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Gagal hash password:", err)
	}

	user := domain.User{
		Email:    "user@example.com",
		Password: string(hashedPassword),
		Name:     "user example",
		IsActive: true,
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatal("Gagal membuat user seed:", err)
	}

	log.Println("User seed berhasil dibuat dengan email: user@example.com")
}
