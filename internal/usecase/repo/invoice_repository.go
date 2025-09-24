package repo

import (
	"database/sql"
	"fiber-boiler-plate/internal/domain"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	GetAll(req *domain.InvoiceListRequest) ([]*domain.Invoice, int, error)
	GetByID(id string) (*domain.Invoice, error)
	UpdateStatus(id string, status string, keterangan *string) error
	GetStatistics(req *domain.InvoiceStatisticsRequest) (*domain.InvoiceStatistics, error)
}

type invoiceRepository struct {
	db    *gorm.DB
	redis RedisRepository
}

func NewInvoiceRepository(db *gorm.DB, redis RedisRepository) InvoiceRepository {
	return &invoiceRepository{
		db:    db,
		redis: redis,
	}
}

func (r *invoiceRepository) GetAll(req *domain.InvoiceListRequest) ([]*domain.Invoice, int, error) {
	cacheKey := fmt.Sprintf("invoice:list:page:%d:per_page:%d:sort:%s:%s",
		req.Page, req.PerPage, req.SortBy, req.SortDirection)

	if req.Search != nil && *req.Search != "" {
		cacheKey += fmt.Sprintf(":search:%s", *req.Search)
	}
	if req.Status != nil {
		cacheKey += fmt.Sprintf(":status:%s", *req.Status)
	}
	if req.NamaUser != nil {
		cacheKey += fmt.Sprintf(":nama_user:%s", *req.NamaUser)
	}
	if req.TanggalMulai != nil {
		cacheKey += fmt.Sprintf(":tanggal_mulai:%s", *req.TanggalMulai)
	}
	if req.TanggalSelesai != nil {
		cacheKey += fmt.Sprintf(":tanggal_selesai:%s", *req.TanggalSelesai)
	}

	var invoices []*domain.Invoice
	var total int64

	query := r.db.Model(&domain.Invoice{}).
		Preload("User").
		Preload("SubscriptionPlan")

	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		query = query.Where("invoices.id ILIKE ?", searchTerm)
	}

	if req.Status != nil && *req.Status != "" {
		query = query.Where("invoices.status = ?", *req.Status)
	}

	if req.NamaUser != nil && *req.NamaUser != "" {
		namaTerm := "%" + *req.NamaUser + "%"
		query = query.Joins("JOIN users ON users.id = invoices.user_id").
			Where("users.name ILIKE ?", namaTerm)
	}

	if req.TanggalMulai != nil && *req.TanggalMulai != "" {
		query = query.Where("DATE(invoices.dibayar_pada) >= ?", *req.TanggalMulai)
	}

	if req.TanggalSelesai != nil && *req.TanggalSelesai != "" {
		query = query.Where("DATE(invoices.dibayar_pada) <= ?", *req.TanggalSelesai)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sortColumn := req.SortBy
	if sortColumn == "nama_user" {
		sortColumn = "users.name"
		query = query.Joins("LEFT JOIN users ON users.id = invoices.user_id")
	} else if sortColumn == "invoice_id" {
		sortColumn = "invoices.id"
	} else if sortColumn == "dibayar_pada" {
		sortColumn = "invoices.dibayar_pada"
	} else {
		sortColumn = "invoices." + sortColumn
	}

	orderClause := sortColumn + " " + strings.ToUpper(req.SortDirection)

	if err := query.Order(orderClause).
		Limit(req.PerPage).
		Offset(req.GetOffset()).
		Find(&invoices).Error; err != nil {
		return nil, 0, err
	}

	r.redis.SetJSON(cacheKey, invoices, 5*time.Minute)

	return invoices, int(total), nil
}

func (r *invoiceRepository) GetByID(id string) (*domain.Invoice, error) {
	cacheKey := fmt.Sprintf("invoice:detail:%s", id)

	var invoice *domain.Invoice
	if err := r.redis.GetJSON(cacheKey, &invoice); err == nil && invoice != nil {
		return invoice, nil
	}

	invoice = &domain.Invoice{}
	if err := r.db.
		Preload("User").
		Preload("SubscriptionPlan").
		Where("id = ?", id).
		First(invoice).Error; err != nil {
		return nil, err
	}

	r.redis.SetJSON(cacheKey, invoice, 10*time.Minute)

	return invoice, nil
}

func (r *invoiceRepository) UpdateStatus(id string, status string, keterangan *string) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	if status == "sukses" {
		now := time.Now()
		updates["dibayar_pada"] = &now
	}

	if keterangan != nil {
		updates["keterangan"] = *keterangan
	}

	if err := r.db.Model(&domain.Invoice{}).
		Where("id = ?", id).
		Updates(updates).Error; err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("invoice:detail:%s", id)
	r.redis.Delete(cacheKey)

	listCachePattern := "invoice:list:*"
	keys, _ := r.redis.GetKeys(listCachePattern)
	for _, key := range keys {
		r.redis.Delete(key)
	}

	statsCachePattern := "invoice:stats:*"
	statsKeys, _ := r.redis.GetKeys(statsCachePattern)
	for _, key := range statsKeys {
		r.redis.Delete(key)
	}

	return nil
}

func (r *invoiceRepository) GetStatistics(req *domain.InvoiceStatisticsRequest) (*domain.InvoiceStatistics, error) {
	cacheKey := "invoice:stats:all"
	if req.Bulan != nil {
		cacheKey += fmt.Sprintf(":bulan:%d", *req.Bulan)
	}
	if req.Tahun != nil {
		cacheKey += fmt.Sprintf(":tahun:%d", *req.Tahun)
	}

	var stats *domain.InvoiceStatistics
	if err := r.redis.GetJSON(cacheKey, &stats); err == nil && stats != nil {
		return stats, nil
	}

	stats = &domain.InvoiceStatistics{}

	query := r.db.Model(&domain.Invoice{})

	if req.Bulan != nil && req.Tahun != nil {
		query = query.Where("MONTH(created_at) = ? AND YEAR(created_at) = ?", *req.Bulan, *req.Tahun)
	} else if req.Tahun != nil {
		query = query.Where("YEAR(created_at) = ?", *req.Tahun)
	}

	if err := query.Count(&stats.TotalInvoice).Error; err != nil {
		return nil, err
	}

	var totalSukses, totalGagal, totalPending int64
	query.Where("status = ?", "sukses").Count(&totalSukses)
	query.Where("status = ?", "gagal").Count(&totalGagal)
	query.Where("status = ?", "pending").Count(&totalPending)

	stats.TotalSukses = totalSukses
	stats.TotalGagal = totalGagal
	stats.TotalPending = totalPending

	var totalPendapatan sql.NullFloat64
	r.db.Model(&domain.Invoice{}).
		Where("status = ?", "sukses").
		Select("SUM(jumlah)").
		Scan(&totalPendapatan)

	if totalPendapatan.Valid {
		stats.TotalPendapatan = totalPendapatan.Float64
	}

	if stats.TotalSukses > 0 {
		stats.RataRataPembayaran = stats.TotalPendapatan / float64(stats.TotalSukses)
	}

	r.getMonthlyStats(stats, req)
	r.getTopSubscriptionPlans(stats, req)

	r.redis.SetJSON(cacheKey, stats, 15*time.Minute)

	return stats, nil
}

func (r *invoiceRepository) getMonthlyStats(stats *domain.InvoiceStatistics, req *domain.InvoiceStatisticsRequest) {
	query := `
		SELECT 
			MONTH(created_at) as month,
			MONTHNAME(created_at) as month_name,
			COUNT(*) as total_invoice,
			SUM(CASE WHEN status = 'sukses' THEN 1 ELSE 0 END) as total_sukses,
			SUM(CASE WHEN status = 'sukses' THEN jumlah ELSE 0 END) as total_pendapatan
		FROM invoices 
		WHERE 1=1
	`

	var args []interface{}
	if req.Tahun != nil {
		query += " AND YEAR(created_at) = ?"
		args = append(args, *req.Tahun)
	}

	query += " GROUP BY MONTH(created_at), MONTHNAME(created_at) ORDER BY MONTH(created_at)"

	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	stats.InvoiceBulanan = []domain.InvoiceStatsBulanan{}
	for rows.Next() {
		var month int
		var monthName string
		var totalInvoice, totalSukses int64
		var totalPendapatan float64

		rows.Scan(&month, &monthName, &totalInvoice, &totalSukses, &totalPendapatan)

		stats.InvoiceBulanan = append(stats.InvoiceBulanan, domain.InvoiceStatsBulanan{
			Bulan:           monthName,
			TotalInvoice:    totalInvoice,
			TotalSukses:     totalSukses,
			TotalPendapatan: totalPendapatan,
		})
	}
}

func (r *invoiceRepository) getTopSubscriptionPlans(stats *domain.InvoiceStatistics, req *domain.InvoiceStatisticsRequest) {
	query := `
		SELECT 
			sp.nama as subscription_plan_nama,
			COUNT(*) as jumlah_invoice,
			SUM(CASE WHEN i.status = 'sukses' THEN i.jumlah ELSE 0 END) as total_pendapatan
		FROM invoices i 
		LEFT JOIN subscription_plans sp ON sp.id = i.subscription_plan_id 
		WHERE sp.nama IS NOT NULL
	`

	var args []interface{}
	if req.Tahun != nil {
		query += " AND YEAR(i.created_at) = ?"
		args = append(args, *req.Tahun)
	}
	if req.Bulan != nil {
		query += " AND MONTH(i.created_at) = ?"
		args = append(args, *req.Bulan)
	}

	query += " GROUP BY sp.nama ORDER BY total_pendapatan DESC LIMIT 10"

	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return
	}
	defer rows.Close()

	stats.TopSubscriptionPlans = []domain.TopSubscriptionPlan{}
	for rows.Next() {
		var subscriptionPlanNama string
		var jumlahInvoice int64
		var totalPendapatan float64

		rows.Scan(&subscriptionPlanNama, &jumlahInvoice, &totalPendapatan)

		stats.TopSubscriptionPlans = append(stats.TopSubscriptionPlans, domain.TopSubscriptionPlan{
			SubscriptionPlanNama: subscriptionPlanNama,
			JumlahInvoice:        jumlahInvoice,
			TotalPendapatan:      totalPendapatan,
		})
	}
}
