package usecase

import (
	"encoding/json"
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TransaksiUsecase interface {
	GetTransaksiList(userID uint, req *domain.TransaksiListRequest) (*domain.TransaksiListResponse, error)
	GetTransaksiDetail(id string, userID uint) (*domain.TransaksiDetailResponse, error)
	CreateTransaksi(userID uint, req *domain.CreateTransaksiRequest) (*domain.TransaksiDetailResponse, error)
	UpdateTransaksi(id string, userID uint, req *domain.UpdateTransaksiRequest) (*domain.TransaksiDetailResponse, error)
	PatchTransaksi(id string, userID uint, req *domain.PatchTransaksiRequest) (*domain.TransaksiDetailResponse, error)
	DeleteTransaksi(id string, userID uint) error
	SetAnggaranUsecase(anggaranUsecase AnggaranUsecase)
}

type transaksiUsecase struct {
	transaksiRepo   repo.TransaksiRepository
	kantongRepo     repo.KantongRepository
	redisRepo       repo.RedisRepository
	anggaranUsecase AnggaranUsecase
}

func NewTransaksiUsecase(
	transaksiRepo repo.TransaksiRepository,
	kantongRepo repo.KantongRepository,
	redisRepo repo.RedisRepository,
) TransaksiUsecase {
	return &transaksiUsecase{
		transaksiRepo:   transaksiRepo,
		kantongRepo:     kantongRepo,
		redisRepo:       redisRepo,
		anggaranUsecase: nil,
	}
}

func (uc *transaksiUsecase) SetAnggaranUsecase(anggaranUsecase AnggaranUsecase) {
	uc.anggaranUsecase = anggaranUsecase
}

func (uc *transaksiUsecase) GetTransaksiList(userID uint, req *domain.TransaksiListRequest) (*domain.TransaksiListResponse, error) {
	cacheKey := uc.generateListCacheKey(userID, req)

	var cachedResponse domain.TransaksiListResponse
	if err := uc.redisRepo.GetJSON(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	transaksiList, total, err := uc.transaksiRepo.GetByUserID(userID, req)
	if err != nil {
		return nil, err
	}

	totalPages := (total + req.PerPage - 1) / req.PerPage
	if totalPages == 0 {
		totalPages = 1
	}

	response := &domain.TransaksiListResponse{
		Success: true,
		Message: "Daftar transaksi berhasil diambil",
		Code:    200,
		Data:    make([]domain.TransaksiResponse, 0),
		Meta: domain.PaginationMeta{
			CurrentPage:  req.Page,
			TotalPages:   totalPages,
			TotalRecords: total,
			PerPage:      req.PerPage,
		},
		Timestamp: time.Now(),
	}

	for _, transaksi := range transaksiList {
		response.Data = append(response.Data, *transaksi)
	}

	uc.redisRepo.SetJSON(cacheKey, response, 5*time.Minute)

	return response, nil
}

func (uc *transaksiUsecase) GetTransaksiDetail(id string, userID uint) (*domain.TransaksiDetailResponse, error) {
	cacheKey := uc.generateDetailCacheKey(id, userID)

	var cachedResponse domain.TransaksiDetailResponse
	if err := uc.redisRepo.GetJSON(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	transaksi, err := uc.transaksiRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	response := &domain.TransaksiDetailResponse{
		Success:   true,
		Message:   "Detail transaksi berhasil diambil",
		Code:      200,
		Data:      *transaksi,
		Timestamp: time.Now(),
	}

	uc.redisRepo.SetJSON(cacheKey, response, 10*time.Minute)

	return response, nil
}

func (uc *transaksiUsecase) CreateTransaksi(userID uint, req *domain.CreateTransaksiRequest) (*domain.TransaksiDetailResponse, error) {
	_, err := uc.kantongRepo.GetByID(req.KantongID, userID)
	if err != nil {
		return nil, errors.New("kantong tidak ditemukan")
	}

	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid")
	}

	transaksi := &domain.Transaksi{
		ID:        uuid.New().String(),
		UserID:    userID,
		KantongID: req.KantongID,
		Tanggal:   tanggal,
		Jenis:     req.Jenis,
		Jumlah:    req.Jumlah,
		Catatan:   req.Catatan,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.transaksiRepo.Create(transaksi); err != nil {
		return nil, err
	}

	if uc.anggaranUsecase != nil {
		uc.anggaranUsecase.UpdateAnggaranAfterTransaction(req.KantongID, userID)
	}

	uc.invalidateUserCache(userID)

	return uc.GetTransaksiDetail(transaksi.ID, userID)
}

func (uc *transaksiUsecase) UpdateTransaksi(id string, userID uint, req *domain.UpdateTransaksiRequest) (*domain.TransaksiDetailResponse, error) {
	_, err := uc.transaksiRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	_, err = uc.kantongRepo.GetByID(req.KantongID, userID)
	if err != nil {
		return nil, errors.New("kantong tidak ditemukan")
	}

	tanggal, err := time.Parse("2006-01-02", req.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal tidak valid")
	}

	transaksi := &domain.Transaksi{
		ID:        id,
		UserID:    userID,
		KantongID: req.KantongID,
		Tanggal:   tanggal,
		Jenis:     req.Jenis,
		Jumlah:    req.Jumlah,
		Catatan:   req.Catatan,
		UpdatedAt: time.Now(),
	}

	if err := uc.transaksiRepo.Update(transaksi); err != nil {
		return nil, err
	}

	uc.invalidateUserCache(userID)

	return uc.GetTransaksiDetail(id, userID)
}

func (uc *transaksiUsecase) PatchTransaksi(id string, userID uint, req *domain.PatchTransaksiRequest) (*domain.TransaksiDetailResponse, error) {
	existingTransaksi, err := uc.transaksiRepo.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	var updateReq domain.UpdateTransaksiRequest
	updateReq.Tanggal = existingTransaksi.Tanggal
	updateReq.Jenis = existingTransaksi.Jenis
	updateReq.Jumlah = existingTransaksi.Jumlah
	updateReq.KantongID = existingTransaksi.KantongID
	updateReq.Catatan = existingTransaksi.Catatan

	if req.Tanggal != nil {
		updateReq.Tanggal = *req.Tanggal
	}
	if req.Jenis != nil {
		updateReq.Jenis = *req.Jenis
	}
	if req.Jumlah != nil {
		updateReq.Jumlah = *req.Jumlah
	}
	if req.KantongID != nil {
		updateReq.KantongID = *req.KantongID
	}
	if req.Catatan != nil {
		updateReq.Catatan = req.Catatan
	}

	return uc.UpdateTransaksi(id, userID, &updateReq)
}

func (uc *transaksiUsecase) DeleteTransaksi(id string, userID uint) error {
	_, err := uc.transaksiRepo.GetByID(id, userID)
	if err != nil {
		return err
	}

	if err := uc.transaksiRepo.Delete(id, userID); err != nil {
		return err
	}

	uc.invalidateUserCache(userID)
	uc.redisRepo.Delete(uc.generateDetailCacheKey(id, userID))

	return nil
}

func (uc *transaksiUsecase) generateListCacheKey(userID uint, req *domain.TransaksiListRequest) string {
	params := make(map[string]interface{})

	if req.Search != nil {
		params["search"] = *req.Search
	}
	if req.Jenis != nil {
		params["jenis"] = *req.Jenis
	}
	if req.KantongNama != nil {
		params["kantong_nama"] = *req.KantongNama
	}
	if req.TanggalMulai != nil {
		params["tanggal_mulai"] = *req.TanggalMulai
	}
	if req.TanggalSelesai != nil {
		params["tanggal_selesai"] = *req.TanggalSelesai
	}

	params["sort_by"] = req.SortBy
	params["sort_direction"] = req.SortDirection
	params["page"] = req.Page
	params["per_page"] = req.PerPage

	paramsJSON, _ := json.Marshal(params)
	return fmt.Sprintf("transaksi_list:%d:%s", userID, string(paramsJSON))
}

func (uc *transaksiUsecase) generateDetailCacheKey(id string, userID uint) string {
	return fmt.Sprintf("transaksi_detail:%s:%d", id, userID)
}

func (uc *transaksiUsecase) invalidateUserCache(userID uint) {
	pattern := fmt.Sprintf("transaksi_list:%d:*", userID)

	keys, err := uc.redisRepo.Get(pattern)
	if err == nil {
		var keyList []string
		json.Unmarshal([]byte(keys), &keyList)
		for _, key := range keyList {
			uc.redisRepo.Delete(key)
		}
	}
}
