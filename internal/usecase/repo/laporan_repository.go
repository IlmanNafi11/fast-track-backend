package repo

import (
	"fiber-boiler-plate/internal/domain"
	"time"

	"gorm.io/gorm"
)

type laporanRepository struct {
	db *gorm.DB
}

func NewLaporanRepository(db *gorm.DB) LaporanRepository {
	return &laporanRepository{db: db}
}

func (r *laporanRepository) GetRingkasanLaporan(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.RingkasanLaporan, error) {
	var totalPemasukan, totalPengeluaran float64

	err := r.db.Table("transaksis").
		Where("user_id = ? AND tanggal BETWEEN ? AND ? AND jenis = ?", userID, tanggalMulai, tanggalSelesai, "Pemasukan").
		Select("COALESCE(SUM(jumlah), 0)").
		Row().
		Scan(&totalPemasukan)
	if err != nil {
		return nil, err
	}

	err = r.db.Table("transaksis").
		Where("user_id = ? AND tanggal BETWEEN ? AND ? AND jenis = ?", userID, tanggalMulai, tanggalSelesai, "Pengeluaran").
		Select("COALESCE(SUM(jumlah), 0)").
		Row().
		Scan(&totalPengeluaran)
	if err != nil {
		return nil, err
	}

	var totalSaldo float64
	err = r.db.Table("kantongs").
		Where("user_id = ?", userID).
		Select("COALESCE(SUM(saldo), 0)").
		Row().
		Scan(&totalSaldo)
	if err != nil {
		return nil, err
	}

	days := int(tanggalSelesai.Sub(tanggalMulai).Hours()/24) + 1
	rataRataPengeluaranHarian := float64(0)
	if days > 0 {
		rataRataPengeluaranHarian = totalPengeluaran / float64(days)
	}

	return &domain.RingkasanLaporan{
		TotalPemasukan:            totalPemasukan,
		TotalPengeluaran:          totalPengeluaran,
		TotalSaldo:                totalSaldo,
		RataRataPengeluaranHarian: rataRataPengeluaranHarian,
		Periode: domain.PeriodeTanggal{
			TanggalMulai:   tanggalMulai.Format("2006-01-02"),
			TanggalSelesai: tanggalSelesai.Format("2006-01-02"),
		},
	}, nil
}

func (r *laporanRepository) GetStatistikTahunan(userID uint, tahun int) (*domain.StatistikTahunan, error) {
	var results []struct {
		Bulan            int     `json:"bulan"`
		TotalPemasukan   float64 `json:"total_pemasukan"`
		TotalPengeluaran float64 `json:"total_pengeluaran"`
	}

	query := `
		SELECT 
			EXTRACT(MONTH FROM tanggal) as bulan,
			COALESCE(SUM(CASE WHEN jenis = 'Pemasukan' THEN jumlah ELSE 0 END), 0) as total_pemasukan,
			COALESCE(SUM(CASE WHEN jenis = 'Pengeluaran' THEN jumlah ELSE 0 END), 0) as total_pengeluaran
		FROM transaksis 
		WHERE user_id = ? AND EXTRACT(YEAR FROM tanggal) = ?
		GROUP BY EXTRACT(MONTH FROM tanggal)
		ORDER BY bulan
	`

	err := r.db.Raw(query, userID, tahun).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	namaBulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	monthlyData := make(map[int]struct {
		TotalPemasukan   float64
		TotalPengeluaran float64
	})

	for _, result := range results {
		monthlyData[result.Bulan] = struct {
			TotalPemasukan   float64
			TotalPengeluaran float64
		}{
			TotalPemasukan:   result.TotalPemasukan,
			TotalPengeluaran: result.TotalPengeluaran,
		}
	}

	var dataBulanan []domain.DataBulanan
	var totalPemasukanTahun, totalPengeluaranTahun float64

	for bulan := 1; bulan <= 12; bulan++ {
		data := monthlyData[bulan]
		dataBulanan = append(dataBulanan, domain.DataBulanan{
			Bulan:            bulan,
			NamaBulan:        namaBulan[bulan],
			TotalPemasukan:   data.TotalPemasukan,
			TotalPengeluaran: data.TotalPengeluaran,
		})
		totalPemasukanTahun += data.TotalPemasukan
		totalPengeluaranTahun += data.TotalPengeluaran
	}

	return &domain.StatistikTahunan{
		Tahun:                 tahun,
		DataBulanan:           dataBulanan,
		TotalPemasukanTahun:   totalPemasukanTahun,
		TotalPengeluaranTahun: totalPengeluaranTahun,
	}, nil
}

func (r *laporanRepository) GetStatistikKantongBulanan(userID uint, bulan, tahun int) (*domain.StatistikKantongBulanan, error) {
	var results []struct {
		KantongID        string  `json:"kantong_id"`
		KantongNama      string  `json:"kantong_nama"`
		Kategori         string  `json:"kategori"`
		TotalPengeluaran float64 `json:"total_pengeluaran"`
		JumlahTransaksi  int     `json:"jumlah_transaksi"`
	}

	query := `
		SELECT 
			k.id as kantong_id,
			k.nama as kantong_nama,
			k.kategori,
			COALESCE(SUM(t.jumlah), 0) as total_pengeluaran,
			COUNT(t.id) as jumlah_transaksi
		FROM kantongs k
		LEFT JOIN transaksis t ON k.id = t.kantong_id 
			AND t.jenis = 'Pengeluaran' 
			AND EXTRACT(MONTH FROM t.tanggal) = ? 
			AND EXTRACT(YEAR FROM t.tanggal) = ?
		WHERE k.user_id = ?
		GROUP BY k.id, k.nama, k.kategori
		HAVING COALESCE(SUM(t.jumlah), 0) > 0
		ORDER BY total_pengeluaran DESC
	`

	err := r.db.Raw(query, bulan, tahun, userID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var totalPengeluaran float64
	var totalTransaksi int

	for _, result := range results {
		totalPengeluaran += result.TotalPengeluaran
		totalTransaksi += result.JumlahTransaksi
	}

	var dataKantong []domain.DataKantongBulanan
	for _, result := range results {
		persentase := float64(0)
		if totalPengeluaran > 0 {
			persentase = (result.TotalPengeluaran / totalPengeluaran) * 100
		}

		dataKantong = append(dataKantong, domain.DataKantongBulanan{
			KantongID:        result.KantongID,
			KantongNama:      result.KantongNama,
			Kategori:         result.Kategori,
			TotalPengeluaran: result.TotalPengeluaran,
			JumlahTransaksi:  result.JumlahTransaksi,
			Persentase:       persentase,
		})
	}

	namaBulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	return &domain.StatistikKantongBulanan{
		Periode: domain.PeriodeBulan{
			Bulan:     bulan,
			NamaBulan: namaBulan[bulan],
			Tahun:     tahun,
		},
		DataKantong:      dataKantong,
		TotalPengeluaran: totalPengeluaran,
		TotalTransaksi:   totalTransaksi,
	}, nil
}

func (r *laporanRepository) GetTopKantongPengeluaran(userID uint, bulan, tahun, limit int) (*domain.TopKantongPengeluaran, error) {
	var results []struct {
		KantongID        string  `json:"kantong_id"`
		KantongNama      string  `json:"kantong_nama"`
		Kategori         string  `json:"kategori"`
		TotalPengeluaran float64 `json:"total_pengeluaran"`
		JumlahTransaksi  int     `json:"jumlah_transaksi"`
	}

	query := `
		SELECT 
			k.id as kantong_id,
			k.nama as kantong_nama,
			k.kategori,
			COALESCE(SUM(t.jumlah), 0) as total_pengeluaran,
			COUNT(t.id) as jumlah_transaksi
		FROM kantongs k
		LEFT JOIN transaksis t ON k.id = t.kantong_id 
			AND t.jenis = 'Pengeluaran' 
			AND EXTRACT(MONTH FROM t.tanggal) = ? 
			AND EXTRACT(YEAR FROM t.tanggal) = ?
		WHERE k.user_id = ?
		GROUP BY k.id, k.nama, k.kategori
		HAVING COALESCE(SUM(t.jumlah), 0) > 0
		ORDER BY total_pengeluaran DESC
		LIMIT ?
	`

	err := r.db.Raw(query, bulan, tahun, userID, limit).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var totalPengeluaranSemua float64
	var totalTransaksiSemua int

	queryTotal := `
		SELECT 
			COALESCE(SUM(t.jumlah), 0) as total_pengeluaran,
			COUNT(t.id) as total_transaksi
		FROM transaksis t
		JOIN kantongs k ON t.kantong_id = k.id
		WHERE k.user_id = ? 
			AND t.jenis = 'Pengeluaran' 
			AND EXTRACT(MONTH FROM t.tanggal) = ? 
			AND EXTRACT(YEAR FROM t.tanggal) = ?
	`

	err = r.db.Raw(queryTotal, userID, bulan, tahun).Row().Scan(&totalPengeluaranSemua, &totalTransaksiSemua)
	if err != nil {
		return nil, err
	}

	var topKantong []domain.DataTopKantong
	for i, result := range results {
		persentaseDariTotal := float64(0)
		if totalPengeluaranSemua > 0 {
			persentaseDariTotal = (result.TotalPengeluaran / totalPengeluaranSemua) * 100
		}

		rataRataPengeluaran := float64(0)
		if result.JumlahTransaksi > 0 {
			rataRataPengeluaran = result.TotalPengeluaran / float64(result.JumlahTransaksi)
		}

		topKantong = append(topKantong, domain.DataTopKantong{
			Ranking:             i + 1,
			KantongID:           result.KantongID,
			KantongNama:         result.KantongNama,
			Kategori:            result.Kategori,
			TotalPengeluaran:    result.TotalPengeluaran,
			JumlahTransaksi:     result.JumlahTransaksi,
			PersentaseDariTotal: persentaseDariTotal,
			RataRataPengeluaran: rataRataPengeluaran,
		})
	}

	namaBulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	return &domain.TopKantongPengeluaran{
		Periode: domain.PeriodeBulan{
			Bulan:     bulan,
			NamaBulan: namaBulan[bulan],
			Tahun:     tahun,
		},
		TopKantong:            topKantong,
		TotalPengeluaranSemua: totalPengeluaranSemua,
		TotalTransaksiSemua:   totalTransaksiSemua,
	}, nil
}

func (r *laporanRepository) GetStatistikKantongPeriode(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.StatistikKantongPeriode, error) {
	type queryResult struct {
		KantongID        string  `json:"kantong_id"`
		KantongNama      string  `json:"kantong_nama"`
		TotalPengeluaran float64 `json:"total_pengeluaran"`
	}

	var results []queryResult
	query := `
		SELECT 
			k.id as kantong_id,
			k.nama as kantong_nama,
			COALESCE(SUM(t.jumlah), 0) as total_pengeluaran
		FROM kantongs k
		LEFT JOIN transaksis t ON k.id = t.kantong_id 
			AND t.jenis = 'Pengeluaran' 
			AND t.tanggal BETWEEN ? AND ?
		WHERE k.user_id = ?
		GROUP BY k.id, k.nama
		ORDER BY total_pengeluaran DESC
	`

	err := r.db.Raw(query, tanggalMulai, tanggalSelesai, userID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var totalPengeluaran float64
	var dataKantong []domain.DataKantongPeriode

	for _, result := range results {
		totalPengeluaran += result.TotalPengeluaran
		dataKantong = append(dataKantong, domain.DataKantongPeriode{
			KantongID:        result.KantongID,
			KantongNama:      result.KantongNama,
			TotalPengeluaran: result.TotalPengeluaran,
		})
	}

	return &domain.StatistikKantongPeriode{
		Periode: domain.PeriodeTanggal{
			TanggalMulai:   tanggalMulai.Format("2006-01-02"),
			TanggalSelesai: tanggalSelesai.Format("2006-01-02"),
		},
		DataKantong:      dataKantong,
		TotalPengeluaran: totalPengeluaran,
	}, nil
}

func (r *laporanRepository) GetPengeluaranKantongDetail(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.PengeluaranKantongDetail, error) {
	type queryResult struct {
		KantongID        string  `json:"kantong_id"`
		KantongNama      string  `json:"kantong_nama"`
		TotalPengeluaran float64 `json:"total_pengeluaran"`
		JumlahTransaksi  int     `json:"jumlah_transaksi"`
		SaldoKantong     float64 `json:"saldo_kantong"`
	}

	var results []queryResult
	query := `
		SELECT 
			k.id as kantong_id,
			k.nama as kantong_nama,
			COALESCE(SUM(CASE WHEN t.jenis = 'Pengeluaran' THEN t.jumlah ELSE 0 END), 0) as total_pengeluaran,
			COUNT(CASE WHEN t.jenis = 'Pengeluaran' THEN 1 END) as jumlah_transaksi,
			k.saldo as saldo_kantong
		FROM kantongs k
		LEFT JOIN transaksis t ON k.id = t.kantong_id 
			AND t.tanggal BETWEEN ? AND ?
		WHERE k.user_id = ?
		GROUP BY k.id, k.nama, k.saldo
		ORDER BY total_pengeluaran DESC
	`

	err := r.db.Raw(query, tanggalMulai, tanggalSelesai, userID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var totalPengeluaran, totalSaldoSemuaKantong float64
	var dataKantong []domain.DataKantongDetail

	for _, result := range results {
		totalPengeluaran += result.TotalPengeluaran
		totalSaldoSemuaKantong += result.SaldoKantong

		persentaseDariSaldo := float64(0)
		if result.SaldoKantong > 0 {
			persentaseDariSaldo = (result.TotalPengeluaran / result.SaldoKantong) * 100
		}

		rataRataPengeluaran := float64(0)
		if result.JumlahTransaksi > 0 {
			rataRataPengeluaran = result.TotalPengeluaran / float64(result.JumlahTransaksi)
		}

		dataKantong = append(dataKantong, domain.DataKantongDetail{
			KantongID:           result.KantongID,
			KantongNama:         result.KantongNama,
			TotalPengeluaran:    result.TotalPengeluaran,
			PersentaseDariSaldo: persentaseDariSaldo,
			JumlahTransaksi:     result.JumlahTransaksi,
			RataRataPengeluaran: rataRataPengeluaran,
			SaldoKantong:        result.SaldoKantong,
		})
	}

	return &domain.PengeluaranKantongDetail{
		Periode: domain.PeriodeTanggal{
			TanggalMulai:   tanggalMulai.Format("2006-01-02"),
			TanggalSelesai: tanggalSelesai.Format("2006-01-02"),
		},
		DataKantong:            dataKantong,
		TotalPengeluaran:       totalPengeluaran,
		TotalSaldoSemuaKantong: totalSaldoSemuaKantong,
	}, nil
}
