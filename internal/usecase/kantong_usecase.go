package usecase

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"math"

	"gorm.io/gorm"
)

type KantongUsecase interface {
	GetKantongList(userID uint, req *domain.KantongListRequest) ([]*domain.KantongResponse, *domain.PaginationMeta, error)
	GetKantongByID(id string, userID uint) (*domain.KantongResponse, error)
	CreateKantong(req *domain.CreateKantongRequest, userID uint) (*domain.KantongResponse, error)
	UpdateKantong(id string, req *domain.UpdateKantongRequest, userID uint) (*domain.KantongResponse, error)
	PatchKantong(id string, req *domain.PatchKantongRequest, userID uint) (*domain.KantongResponse, error)
	DeleteKantong(id string, userID uint) error
}

type kantongUsecase struct {
	kantongRepo repo.KantongRepository
	userRepo    repo.UserRepository
}

func NewKantongUsecase(kantongRepo repo.KantongRepository, userRepo repo.UserRepository) KantongUsecase {
	return &kantongUsecase{
		kantongRepo: kantongRepo,
		userRepo:    userRepo,
	}
}

func (u *kantongUsecase) GetKantongList(userID uint, req *domain.KantongListRequest) ([]*domain.KantongResponse, *domain.PaginationMeta, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}
	if req.PerPage > 100 {
		req.PerPage = 100
	}
	if req.SortBy == "" {
		req.SortBy = "nama"
	}
	if req.SortDirection == "" {
		req.SortDirection = "asc"
	}

	kantongs, total, err := u.kantongRepo.GetByUserID(userID, req)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PerPage)))

	meta := &domain.PaginationMeta{
		CurrentPage:  req.Page,
		TotalPages:   totalPages,
		TotalRecords: total,
		PerPage:      req.PerPage,
	}

	return domain.ToKantongResponseList(kantongs), meta, nil
}

func (u *kantongUsecase) GetKantongByID(id string, userID uint) (*domain.KantongResponse, error) {
	kantong, err := u.kantongRepo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kantong tidak ditemukan")
		}
		return nil, err
	}

	return domain.ToKantongResponse(kantong), nil
}

func (u *kantongUsecase) CreateKantong(req *domain.CreateKantongRequest, userID uint) (*domain.KantongResponse, error) {
	if _, err := u.userRepo.GetByID(userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pengguna tidak ditemukan")
		}
		return nil, err
	}

	exists, err := u.kantongRepo.IsNameExistForUser(req.Nama, userID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("nama kantong sudah ada")
	}

	idKartu, err := u.kantongRepo.GenerateUniqueIDKartu()
	if err != nil {
		return nil, err
	}

	saldo := float64(0)
	if req.Saldo != nil {
		saldo = *req.Saldo
	}

	kantong := &domain.Kantong{
		IDKartu:   idKartu,
		UserID:    userID,
		Nama:      req.Nama,
		Kategori:  req.Kategori,
		Deskripsi: req.Deskripsi,
		Limit:     req.Limit,
		Saldo:     saldo,
		Warna:     req.Warna,
	}

	if err := u.kantongRepo.Create(kantong); err != nil {
		return nil, err
	}

	return domain.ToKantongResponse(kantong), nil
}

func (u *kantongUsecase) UpdateKantong(id string, req *domain.UpdateKantongRequest, userID uint) (*domain.KantongResponse, error) {
	kantong, err := u.kantongRepo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kantong tidak ditemukan")
		}
		return nil, err
	}

	if kantong.Nama != req.Nama {
		exists, err := u.kantongRepo.IsNameExistForUser(req.Nama, userID, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("nama kantong sudah ada")
		}
	}

	kantong.Nama = req.Nama
	kantong.Kategori = req.Kategori
	kantong.Deskripsi = req.Deskripsi
	kantong.Limit = req.Limit
	kantong.Saldo = req.Saldo
	kantong.Warna = req.Warna

	if err := u.kantongRepo.Update(kantong); err != nil {
		return nil, err
	}

	return domain.ToKantongResponse(kantong), nil
}

func (u *kantongUsecase) PatchKantong(id string, req *domain.PatchKantongRequest, userID uint) (*domain.KantongResponse, error) {
	kantong, err := u.kantongRepo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("kantong tidak ditemukan")
		}
		return nil, err
	}

	if req.Nama != nil {
		if kantong.Nama != *req.Nama {
			exists, err := u.kantongRepo.IsNameExistForUser(*req.Nama, userID, id)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, errors.New("nama kantong sudah ada")
			}
		}
		kantong.Nama = *req.Nama
	}

	if req.Kategori != nil {
		kantong.Kategori = *req.Kategori
	}

	if req.Deskripsi != nil {
		kantong.Deskripsi = req.Deskripsi
	}

	if req.Limit != nil {
		kantong.Limit = req.Limit
	}

	if req.Saldo != nil {
		kantong.Saldo = *req.Saldo
	}

	if req.Warna != nil {
		kantong.Warna = *req.Warna
	}

	if err := u.kantongRepo.Update(kantong); err != nil {
		return nil, err
	}

	return domain.ToKantongResponse(kantong), nil
}

func (u *kantongUsecase) DeleteKantong(id string, userID uint) error {
	kantong, err := u.kantongRepo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("kantong tidak ditemukan")
		}
		return err
	}

	if kantong.Saldo > 0 {
		return errors.New("kantong tidak dapat dihapus karena masih memiliki saldo")
	}

	return u.kantongRepo.Delete(id, userID)
}
