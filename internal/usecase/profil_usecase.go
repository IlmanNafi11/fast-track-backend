package usecase

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type ProfilUsecase interface {
	GetProfil(userID uint) (*domain.ProfilResponse, error)
	UpdateProfil(userID uint, req domain.UpdateProfilRequest) (*domain.ProfilResponse, error)
}

type profilUsecase struct {
	userRepo  repo.UserRepository
	redisRepo repo.RedisRepository
}

func NewProfilUsecase(
	userRepo repo.UserRepository,
	redisRepo repo.RedisRepository,
) ProfilUsecase {
	return &profilUsecase{
		userRepo:  userRepo,
		redisRepo: redisRepo,
	}
}

func (uc *profilUsecase) GetProfil(userID uint) (*domain.ProfilResponse, error) {
	cacheKey := fmt.Sprintf("profil:user:%d", userID)

	var cachedProfil domain.ProfilResponse
	err := uc.redisRepo.GetJSON(cacheKey, &cachedProfil)
	if err == nil {
		return &cachedProfil, nil
	}

	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profil tidak ditemukan")
		}
		return nil, err
	}

	profil := &domain.ProfilResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	uc.redisRepo.SetJSON(cacheKey, profil, 30*time.Minute)

	return profil, nil
}

func (uc *profilUsecase) UpdateProfil(userID uint, req domain.UpdateProfilRequest) (*domain.ProfilResponse, error) {
	user, err := uc.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("profil tidak ditemukan")
		}
		return nil, err
	}

	user.Name = req.Name

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	profil := &domain.ProfilResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	cacheKey := fmt.Sprintf("profil:user:%d", userID)
	uc.redisRepo.Delete(cacheKey)
	uc.redisRepo.SetJSON(cacheKey, profil, 30*time.Minute)

	return profil, nil
}
