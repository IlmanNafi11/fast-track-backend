package repo

import (
	"fmt"
	"math"
	"time"

	"fiber-boiler-plate/internal/domain"

	"gorm.io/gorm"
)

type anggaranRepository struct {
	db    *gorm.DB
	redis RedisRepository
}

func NewAnggaranRepository(db *gorm.DB, redis RedisRepository) AnggaranRepository {
	return &anggaranRepository{
		db:    db,
		redis: redis,
	}
}

func (r *anggaranRepository) GetByUserID(userID uint, req *domain.AnggaranListRequest) ([]*domain.AnggaranItem, int, error) {
	cacheKey := fmt.Sprintf("anggaran:list:%d:%d:%d", userID, *req.Bulan, *req.Tahun)

	var anggarans []domain.Anggaran
	query := r.db.Preload("Kantong").Where("user_id = ? AND bulan = ? AND tahun = ?", userID, *req.Bulan, *req.Tahun)

	if req.Search != nil && *req.Search != "" {
		query = query.Joins("LEFT JOIN kantongs ON anggarans.kantong_id = kantongs.id").
			Where("kantongs.nama ILIKE ?", "%"+*req.Search+"%")
	}

	var total int64
	query.Model(&domain.Anggaran{}).Count(&total)

	orderClause := r.buildOrderClause(req.SortBy, req.SortDirection)
	offset := (req.Page - 1) * req.PerPage

	if err := query.Order(orderClause).Offset(offset).Limit(req.PerPage).Find(&anggarans).Error; err != nil {
		return nil, 0, err
	}

	items := make([]*domain.AnggaranItem, len(anggarans))
	for i, anggaran := range anggarans {
		item := r.toAnggaranItem(&anggaran)
		item, err := r.calculateAnggaranValues(item, userID)
		if err != nil {
			return nil, 0, err
		}
		items[i] = item
	}

	r.redis.SetJSON(cacheKey, items, 5*time.Minute)

	return items, int(total), nil
}

func (r *anggaranRepository) GetByKantongID(kantongID string, userID uint, bulan, tahun int) (*domain.AnggaranItem, error) {
	cacheKey := fmt.Sprintf("anggaran:detail:%s:%d:%d:%d", kantongID, userID, bulan, tahun)

	var cachedItem domain.AnggaranItem
	if err := r.redis.GetJSON(cacheKey, &cachedItem); err == nil {
		return &cachedItem, nil
	}

	var anggaran domain.Anggaran
	err := r.db.Preload("Kantong").
		Where("kantong_id = ? AND user_id = ? AND bulan = ? AND tahun = ?", kantongID, userID, bulan, tahun).
		First(&anggaran).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return r.createDefaultAnggaran(kantongID, userID, bulan, tahun)
		}
		return nil, err
	}

	item := r.toAnggaranItem(&anggaran)
	item, err = r.calculateAnggaranValues(item, userID)
	if err != nil {
		return nil, err
	}

	statistik, err := r.GetStatistikBulan(kantongID, userID, bulan, tahun)
	if err != nil {
		return nil, err
	}
	item.StatistikBulan = statistik

	r.redis.SetJSON(cacheKey, item, 5*time.Minute)

	return item, nil
}

func (r *anggaranRepository) CreateOrUpdate(anggaran *domain.AnggaranItem) error {
	var existing domain.Anggaran
	err := r.db.Where("kantong_id = ? AND user_id = ? AND bulan = ? AND tahun = ?",
		anggaran.KantongID, anggaran.DetailKantong.UserID, anggaran.Bulan, anggaran.Tahun).
		First(&existing).Error

	dbAnggaran := &domain.Anggaran{
		KantongID:   anggaran.KantongID,
		UserID:      anggaran.DetailKantong.UserID,
		Bulan:       anggaran.Bulan,
		Tahun:       anggaran.Tahun,
		Rencana:     anggaran.Rencana,
		CarryIn:     anggaran.CarryIn,
		Penyesuaian: anggaran.Penyesuaian,
		Terpakai:    anggaran.Terpakai,
		Sisa:        anggaran.Sisa,
		Progres:     anggaran.Progres,
	}

	if err == gorm.ErrRecordNotFound {
		err = r.db.Create(dbAnggaran).Error
	} else {
		existing.Rencana = dbAnggaran.Rencana
		existing.CarryIn = dbAnggaran.CarryIn
		existing.Penyesuaian = dbAnggaran.Penyesuaian
		existing.Terpakai = dbAnggaran.Terpakai
		existing.Sisa = dbAnggaran.Sisa
		existing.Progres = dbAnggaran.Progres
		err = r.db.Save(&existing).Error
	}

	if err == nil {
		r.clearAnggaranCache(anggaran.KantongID, anggaran.DetailKantong.UserID, anggaran.Bulan, anggaran.Tahun)
	}

	return err
}

func (r *anggaranRepository) CreatePenyesuaian(userID uint, req *domain.PenyesuaianAnggaranRequest) (*domain.AnggaranItem, error) {
	_, err := r.GetByKantongID(req.KantongID, userID, req.Bulan, req.Tahun)
	if err != nil {
		return nil, err
	}

	var anggaranDB domain.Anggaran
	err = r.db.Where("kantong_id = ? AND user_id = ? AND bulan = ? AND tahun = ?",
		req.KantongID, userID, req.Bulan, req.Tahun).First(&anggaranDB).Error
	if err != nil {
		return nil, err
	}

	penyesuaian := &domain.PenyesuaianAnggaran{
		AnggaranID: anggaranDB.ID,
		Jenis:      req.Jenis,
		Jumlah:     req.Jumlah,
	}

	err = r.db.Create(penyesuaian).Error
	if err != nil {
		return nil, err
	}

	var totalPenyesuaian float64
	err = r.db.Model(&domain.PenyesuaianAnggaran{}).
		Select("SUM(CASE WHEN jenis = 'tambah' THEN jumlah WHEN jenis = 'kurangi' THEN -jumlah ELSE 0 END)").
		Where("anggaran_id = ?", anggaranDB.ID).Scan(&totalPenyesuaian).Error
	if err != nil {
		return nil, err
	}

	anggaranDB.Penyesuaian = totalPenyesuaian
	anggaranDB.Sisa = r.calculateSisa(anggaranDB.Rencana, anggaranDB.CarryIn, anggaranDB.Penyesuaian, anggaranDB.Terpakai)
	anggaranDB.Progres = r.calculateProgres(anggaranDB.Rencana, anggaranDB.Penyesuaian, anggaranDB.Terpakai)

	err = r.db.Save(&anggaranDB).Error
	if err != nil {
		return nil, err
	}

	r.clearAnggaranCache(req.KantongID, userID, req.Bulan, req.Tahun)

	return r.GetByKantongID(req.KantongID, userID, req.Bulan, req.Tahun)
}

func (r *anggaranRepository) GetStatistikBulan(kantongID string, userID uint, bulan, tahun int) ([]domain.StatistikHarian, error) {
	startDate := time.Date(tahun, time.Month(bulan), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)

	var results []struct {
		Tanggal          time.Time
		JumlahTransaksi  int64
		TotalPengeluaran float64
	}

	err := r.db.Table("transaksis").
		Select("DATE(created_at) as tanggal, COUNT(*) as jumlah_transaksi, SUM(nominal) as total_pengeluaran").
		Where("kantong_id = ? AND user_id = ? AND created_at >= ? AND created_at <= ?",
			kantongID, userID, startDate, endDate).
		Group("DATE(created_at)").
		Order("tanggal").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	statistik := make([]domain.StatistikHarian, len(results))
	var akumulasi float64

	for i, result := range results {
		akumulasi += result.TotalPengeluaran
		statistik[i] = domain.StatistikHarian{
			Tanggal:           result.Tanggal,
			JumlahTransaksi:   int(result.JumlahTransaksi),
			TotalPengeluaran:  result.TotalPengeluaran,
			AkumulasiTerpakai: akumulasi,
		}
	}

	return statistik, nil
}

func (r *anggaranRepository) RecalculateAnggaran(kantongID string, userID uint, bulan, tahun int) (*domain.AnggaranItem, error) {
	var anggaran domain.Anggaran
	err := r.db.Where("kantong_id = ? AND user_id = ? AND bulan = ? AND tahun = ?",
		kantongID, userID, bulan, tahun).First(&anggaran).Error
	if err != nil {
		return nil, err
	}

	var totalTransaksi float64
	err = r.db.Model(&domain.Transaksi{}).
		Where("kantong_id = ? AND user_id = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?",
			kantongID, userID, bulan, tahun).
		Select("COALESCE(SUM(nominal), 0)").Scan(&totalTransaksi).Error
	if err != nil {
		return nil, err
	}

	anggaran.Terpakai = totalTransaksi
	anggaran.Sisa = r.calculateSisa(anggaran.Rencana, anggaran.CarryIn, anggaran.Penyesuaian, anggaran.Terpakai)
	anggaran.Progres = r.calculateProgres(anggaran.Rencana, anggaran.Penyesuaian, anggaran.Terpakai)

	err = r.db.Save(&anggaran).Error
	if err != nil {
		return nil, err
	}

	r.clearAnggaranCache(kantongID, userID, bulan, tahun)

	return r.GetByKantongID(kantongID, userID, bulan, tahun)
}

func (r *anggaranRepository) CreateAnggaranForKantong(kantong *domain.Kantong) error {
	now := time.Now()

	anggaran := &domain.Anggaran{
		KantongID:   kantong.ID,
		UserID:      kantong.UserID,
		Bulan:       int(now.Month()),
		Tahun:       now.Year(),
		Rencana:     kantong.Limit,
		CarryIn:     0,
		Penyesuaian: 0,
		Terpakai:    0,
		Sisa:        0,
		Progres:     0,
	}

	if anggaran.Rencana != nil {
		anggaran.Sisa = *anggaran.Rencana
	}

	return r.db.Create(anggaran).Error
}

func (r *anggaranRepository) UpdateAnggaranAfterTransaksi(kantongID string, userID uint) error {
	now := time.Now()
	return r.RecalculateAnggaranByMonth(kantongID, userID, int(now.Month()), now.Year())
}

func (r *anggaranRepository) RecalculateAnggaranByMonth(kantongID string, userID uint, bulan, tahun int) error {
	_, err := r.RecalculateAnggaran(kantongID, userID, bulan, tahun)
	return err
}

func (r *anggaranRepository) createDefaultAnggaran(kantongID string, userID uint, bulan, tahun int) (*domain.AnggaranItem, error) {
	var kantong domain.Kantong
	err := r.db.Where("id = ? AND user_id = ?", kantongID, userID).First(&kantong).Error
	if err != nil {
		return nil, err
	}

	anggaran := &domain.Anggaran{
		KantongID:   kantongID,
		UserID:      userID,
		Bulan:       bulan,
		Tahun:       tahun,
		Rencana:     kantong.Limit,
		CarryIn:     0,
		Penyesuaian: 0,
		Terpakai:    0,
		Sisa:        0,
		Progres:     0,
	}

	if anggaran.Rencana != nil {
		anggaran.Sisa = *anggaran.Rencana
	}

	err = r.db.Create(anggaran).Error
	if err != nil {
		return nil, err
	}

	return r.toAnggaranItem(anggaran), nil
}

func (r *anggaranRepository) toAnggaranItem(anggaran *domain.Anggaran) *domain.AnggaranItem {
	return &domain.AnggaranItem{
		KantongID:     anggaran.KantongID,
		NamaKantong:   anggaran.Kantong.Nama,
		Rencana:       anggaran.Rencana,
		CarryIn:       anggaran.CarryIn,
		Penyesuaian:   anggaran.Penyesuaian,
		Terpakai:      anggaran.Terpakai,
		Sisa:          anggaran.Sisa,
		Progres:       anggaran.Progres,
		DetailKantong: &anggaran.Kantong,
		Bulan:         anggaran.Bulan,
		Tahun:         anggaran.Tahun,
		CreatedAt:     anggaran.CreatedAt,
		UpdatedAt:     anggaran.UpdatedAt,
	}
}

func (r *anggaranRepository) calculateAnggaranValues(item *domain.AnggaranItem, userID uint) (*domain.AnggaranItem, error) {
	var totalTransaksi float64
	err := r.db.Model(&domain.Transaksi{}).
		Where("kantong_id = ? AND user_id = ? AND EXTRACT(MONTH FROM created_at) = ? AND EXTRACT(YEAR FROM created_at) = ?",
			item.KantongID, userID, item.Bulan, item.Tahun).
		Select("COALESCE(SUM(nominal), 0)").Scan(&totalTransaksi).Error
	if err != nil {
		return nil, err
	}

	item.Terpakai = totalTransaksi
	item.Sisa = r.calculateSisa(item.Rencana, item.CarryIn, item.Penyesuaian, item.Terpakai)
	item.Progres = r.calculateProgres(item.Rencana, item.Penyesuaian, item.Terpakai)

	return item, nil
}

func (r *anggaranRepository) calculateSisa(rencana *float64, carryIn, penyesuaian, terpakai float64) float64 {
	if rencana == nil {
		return carryIn + penyesuaian - terpakai
	}
	return *rencana + carryIn + penyesuaian - terpakai
}

func (r *anggaranRepository) calculateProgres(rencana *float64, penyesuaian, terpakai float64) float64 {
	if rencana == nil || *rencana+penyesuaian <= 0 {
		return 0
	}
	progres := (terpakai / (*rencana + penyesuaian)) * 100
	return math.Round(progres*100) / 100
}

func (r *anggaranRepository) buildOrderClause(sortBy, sortDirection string) string {
	validSortFields := map[string]string{
		"nama_kantong": "kantongs.nama",
		"rencana":      "anggarans.rencana",
		"carry_in":     "anggarans.carry_in",
		"penyesuaian":  "anggarans.penyesuaian",
		"sisa":         "anggarans.sisa",
	}

	field, exists := validSortFields[sortBy]
	if !exists {
		field = "kantongs.nama"
	}

	if sortDirection != "desc" {
		sortDirection = "asc"
	}

	if sortBy == "nama_kantong" {
		return fmt.Sprintf("kantongs.nama %s", sortDirection)
	}

	return fmt.Sprintf("%s %s", field, sortDirection)
}

func (r *anggaranRepository) clearAnggaranCache(kantongID string, userID uint, bulan, tahun int) {
	listCacheKey := fmt.Sprintf("anggaran:list:%d:%d:%d", userID, bulan, tahun)
	detailCacheKey := fmt.Sprintf("anggaran:detail:%s:%d:%d:%d", kantongID, userID, bulan, tahun)

	r.redis.Delete(listCacheKey)
	r.redis.Delete(detailCacheKey)
}
