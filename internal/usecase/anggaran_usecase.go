package usecase

import (
	"errors"
	"time"

	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
)

type AnggaranUsecase interface {
	GetAnggaranList(userID uint, req *domain.AnggaranListRequest) ([]*domain.AnggaranResponse, *domain.PaginationMeta, error)
	GetAnggaranDetail(kantongID string, userID uint, bulan, tahun *int) (*domain.AnggaranDetailResponse, error)
	CreatePenyesuaianAnggaran(userID uint, req *domain.PenyesuaianAnggaranRequest) (*domain.AnggaranResponse, error)
	CreateAnggaranForNewKantong(kantong *domain.Kantong) error
	UpdateAnggaranAfterTransaction(kantongID string, userID uint) error
}

type anggaranUsecase struct {
	anggaranRepo  repo.AnggaranRepository
	kantongRepo   repo.KantongRepository
	transaksiRepo repo.TransaksiRepository
	redisRepo     repo.RedisRepository
}

func NewAnggaranUsecase(
	anggaranRepo repo.AnggaranRepository,
	kantongRepo repo.KantongRepository,
	transaksiRepo repo.TransaksiRepository,
	redisRepo repo.RedisRepository,
) AnggaranUsecase {
	return &anggaranUsecase{
		anggaranRepo:  anggaranRepo,
		kantongRepo:   kantongRepo,
		transaksiRepo: transaksiRepo,
		redisRepo:     redisRepo,
	}
}

func (uc *anggaranUsecase) GetAnggaranList(userID uint, req *domain.AnggaranListRequest) ([]*domain.AnggaranResponse, *domain.PaginationMeta, error) {
	if req.Bulan == nil || req.Tahun == nil {
		now := time.Now()
		if req.Bulan == nil {
			bulan := int(now.Month())
			req.Bulan = &bulan
		}
		if req.Tahun == nil {
			tahun := now.Year()
			req.Tahun = &tahun
		}
	}

	items, total, err := uc.anggaranRepo.GetByUserID(userID, req)
	if err != nil {
		return nil, nil, err
	}

	responses := domain.ToAnggaranResponseList(items)

	totalPages := (total + req.PerPage - 1) / req.PerPage
	meta := &domain.PaginationMeta{
		CurrentPage:  req.Page,
		TotalPages:   totalPages,
		TotalRecords: total,
		PerPage:      req.PerPage,
	}

	return responses, meta, nil
}

func (uc *anggaranUsecase) GetAnggaranDetail(kantongID string, userID uint, bulan, tahun *int) (*domain.AnggaranDetailResponse, error) {
	now := time.Now()
	if bulan == nil {
		defaultBulan := int(now.Month())
		bulan = &defaultBulan
	}
	if tahun == nil {
		defaultTahun := now.Year()
		tahun = &defaultTahun
	}

	kantong, err := uc.kantongRepo.GetByID(kantongID, userID)
	if err != nil {
		return nil, errors.New("kantong tidak ditemukan")
	}

	item, err := uc.anggaranRepo.GetByKantongID(kantongID, userID, *bulan, *tahun)
	if err != nil {
		return nil, err
	}

	item.DetailKantong = kantong

	statistik, err := uc.anggaranRepo.GetStatistikBulan(kantongID, userID, *bulan, *tahun)
	if err != nil {
		return nil, err
	}
	item.StatistikBulan = statistik

	return domain.ToAnggaranDetailResponse(item), nil
}

func (uc *anggaranUsecase) CreatePenyesuaianAnggaran(userID uint, req *domain.PenyesuaianAnggaranRequest) (*domain.AnggaranResponse, error) {
	kantong, err := uc.kantongRepo.GetByID(req.KantongID, userID)
	if err != nil {
		return nil, errors.New("kantong tidak ditemukan")
	}

	if kantong.UserID != userID {
		return nil, errors.New("tidak memiliki akses ke kantong ini")
	}

	item, err := uc.anggaranRepo.CreatePenyesuaian(userID, req)
	if err != nil {
		return nil, err
	}

	return domain.ToAnggaranResponse(item), nil
}

func (uc *anggaranUsecase) CreateAnggaranForNewKantong(kantong *domain.Kantong) error {
	return uc.anggaranRepo.CreateAnggaranForKantong(kantong)
}

func (uc *anggaranUsecase) UpdateAnggaranAfterTransaction(kantongID string, userID uint) error {
	return uc.anggaranRepo.UpdateAnggaranAfterTransaksi(kantongID, userID)
}
