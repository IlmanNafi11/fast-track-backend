package usecase

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/usecase/repo"
	"fmt"
	"time"
)

type LaporanUsecase interface {
	GetRingkasanLaporan(userID uint, req *domain.RingkasanLaporanRequest) (*domain.RingkasanLaporanResponse, error)
	GetStatistikTahunan(userID uint, req *domain.StatistikTahunanRequest) (*domain.StatistikTahunanResponse, error)
	GetStatistikKantongBulanan(userID uint, req *domain.StatistikKantongBulananRequest) (*domain.StatistikKantongBulananResponse, error)
	GetTopKantongPengeluaran(userID uint, req *domain.TopKantongPengeluaranRequest) (*domain.TopKantongPengeluaranResponse, error)
	GetStatistikKantongPeriode(userID uint, req *domain.StatistikKantongPeriodeRequest) (*domain.StatistikKantongPeriodeResponse, error)
	GetPengeluaranKantongDetail(userID uint, req *domain.PengeluaranKantongDetailRequest) (*domain.PengeluaranKantongDetailResponse, error)
}

type laporanUsecase struct {
	laporanRepo repo.LaporanRepository
	redisRepo   repo.RedisRepository
}

func NewLaporanUsecase(
	laporanRepo repo.LaporanRepository,
	redisRepo repo.RedisRepository,
) LaporanUsecase {
	return &laporanUsecase{
		laporanRepo: laporanRepo,
		redisRepo:   redisRepo,
	}
}

func (uc *laporanUsecase) GetRingkasanLaporan(userID uint, req *domain.RingkasanLaporanRequest) (*domain.RingkasanLaporanResponse, error) {
	tanggalMulai, tanggalSelesai := uc.getDefaultDateRange(req.TanggalMulai, req.TanggalSelesai)

	cacheKey := fmt.Sprintf("laporan:ringkasan:%d:%s:%s", userID, tanggalMulai.Format("2006-01-02"), tanggalSelesai.Format("2006-01-02"))

	var cachedResponse domain.RingkasanLaporanResponse
	if err := uc.redisRepo.GetJSON(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	ringkasan, err := uc.laporanRepo.GetRingkasanLaporan(userID, tanggalMulai, tanggalSelesai)
	if err != nil {
		return nil, err
	}

	response := &domain.RingkasanLaporanResponse{
		Success:   true,
		Message:   "Ringkasan laporan berhasil diambil",
		Code:      200,
		Data:      *ringkasan,
		Timestamp: time.Now(),
	}

	uc.redisRepo.SetJSON(cacheKey, response, 10*time.Minute)

	return response, nil
}

func (uc *laporanUsecase) GetStatistikTahunan(userID uint, req *domain.StatistikTahunanRequest) (*domain.StatistikTahunanResponse, error) {
	tahun := time.Now().Year()
	if req.Tahun != nil {
		tahun = *req.Tahun
	}

	cacheKey := fmt.Sprintf("laporan:statistik_tahunan:%d:%d", userID, tahun)

	var cachedResponse domain.StatistikTahunanResponse
	if err := uc.redisRepo.GetJSON(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	statistik, err := uc.laporanRepo.GetStatistikTahunan(userID, tahun)
	if err != nil {
		return nil, err
	}

	response := &domain.StatistikTahunanResponse{
		Success:   true,
		Message:   "Statistik tahunan berhasil diambil",
		Code:      200,
		Data:      *statistik,
		Timestamp: time.Now(),
	}

	uc.redisRepo.SetJSON(cacheKey, response, 30*time.Minute)

	return response, nil
}

func (uc *laporanUsecase) GetStatistikKantongBulanan(userID uint, req *domain.StatistikKantongBulananRequest) (*domain.StatistikKantongBulananResponse, error) {
	bulan, tahun := uc.getDefaultMonth(req.Bulan, req.Tahun)

	cacheKey := fmt.Sprintf("laporan:statistik_kantong_bulanan:%d:%d:%d", userID, bulan, tahun)

	var cachedResponse domain.StatistikKantongBulananResponse
	if err := uc.redisRepo.GetJSON(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	statistik, err := uc.laporanRepo.GetStatistikKantongBulanan(userID, bulan, tahun)
	if err != nil {
		return nil, err
	}

	response := &domain.StatistikKantongBulananResponse{
		Success:   true,
		Message:   "Statistik kantong bulanan berhasil diambil",
		Code:      200,
		Data:      *statistik,
		Timestamp: time.Now(),
	}

	uc.redisRepo.SetJSON(cacheKey, response, 15*time.Minute)

	return response, nil
}

func (uc *laporanUsecase) GetTopKantongPengeluaran(userID uint, req *domain.TopKantongPengeluaranRequest) (*domain.TopKantongPengeluaranResponse, error) {
	bulan, tahun := uc.getDefaultMonth(req.Bulan, req.Tahun)
	limit := 5
	if req.Limit != nil {
		limit = *req.Limit
	}

	cacheKey := fmt.Sprintf("laporan:top_kantong:%d:%d:%d:%d", userID, bulan, tahun, limit)

	var cachedResponse domain.TopKantongPengeluaranResponse
	if err := uc.redisRepo.GetJSON(cacheKey, &cachedResponse); err == nil {
		return &cachedResponse, nil
	}

	topKantong, err := uc.laporanRepo.GetTopKantongPengeluaran(userID, bulan, tahun, limit)
	if err != nil {
		return nil, err
	}

	response := &domain.TopKantongPengeluaranResponse{
		Success:   true,
		Message:   "Top kantong pengeluaran berhasil diambil",
		Code:      200,
		Data:      *topKantong,
		Timestamp: time.Now(),
	}

	uc.redisRepo.SetJSON(cacheKey, response, 15*time.Minute)

	return response, nil
}

func (uc *laporanUsecase) getDefaultDateRange(tanggalMulai, tanggalSelesai *string) (time.Time, time.Time) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	start := startOfMonth
	end := endOfMonth

	if tanggalMulai != nil {
		if parsed, err := time.Parse("2006-01-02", *tanggalMulai); err == nil {
			start = parsed
		}
	}

	if tanggalSelesai != nil {
		if parsed, err := time.Parse("2006-01-02", *tanggalSelesai); err == nil {
			end = parsed
		}
	}

	return start, end
}

func (uc *laporanUsecase) getDefaultMonth(bulan, tahun *int) (int, int) {
	now := time.Now()
	currentMonth := int(now.Month())
	currentYear := now.Year()

	month := currentMonth
	year := currentYear

	if bulan != nil {
		month = *bulan
	}

	if tahun != nil {
		year = *tahun
	}

	return month, year
}

func (uc *laporanUsecase) GetStatistikKantongPeriode(userID uint, req *domain.StatistikKantongPeriodeRequest) (*domain.StatistikKantongPeriodeResponse, error) {
	cacheKey := fmt.Sprintf("statistik_kantong_periode:%d:%v:%v", userID, req.TanggalMulai, req.TanggalSelesai)

	var response *domain.StatistikKantongPeriodeResponse
	err := uc.redisRepo.GetJSON(cacheKey, &response)
	if err == nil && response != nil {
		return response, nil
	}

	tanggalMulai, tanggalSelesai := uc.getDefaultDateRange(req.TanggalMulai, req.TanggalSelesai)

	data, err := uc.laporanRepo.GetStatistikKantongPeriode(userID, tanggalMulai, tanggalSelesai)
	if err != nil {
		return nil, err
	}

	response = &domain.StatistikKantongPeriodeResponse{
		Success:   true,
		Message:   "Statistik kantong periode berhasil diambil",
		Code:      200,
		Data:      *data,
		Timestamp: time.Now(),
	}

	uc.redisRepo.SetJSON(cacheKey, response, 15*time.Minute)

	return response, nil
}

func (uc *laporanUsecase) GetPengeluaranKantongDetail(userID uint, req *domain.PengeluaranKantongDetailRequest) (*domain.PengeluaranKantongDetailResponse, error) {
	cacheKey := fmt.Sprintf("pengeluaran_kantong_detail:%d:%v:%v", userID, req.TanggalMulai, req.TanggalSelesai)

	var response *domain.PengeluaranKantongDetailResponse
	err := uc.redisRepo.GetJSON(cacheKey, &response)
	if err == nil && response != nil {
		return response, nil
	}

	tanggalMulai, tanggalSelesai := uc.getDefaultDateRange(req.TanggalMulai, req.TanggalSelesai)

	data, err := uc.laporanRepo.GetPengeluaranKantongDetail(userID, tanggalMulai, tanggalSelesai)
	if err != nil {
		return nil, err
	}

	response = &domain.PengeluaranKantongDetailResponse{
		Success:   true,
		Message:   "Detail pengeluaran kantong berhasil diambil",
		Code:      200,
		Data:      *data,
		Timestamp: time.Now(),
	}

	uc.redisRepo.SetJSON(cacheKey, response, 20*time.Minute)

	return response, nil
}
