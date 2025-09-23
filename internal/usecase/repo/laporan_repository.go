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

func (r *laporanRepository) GetTrenBulanan(userID uint, tahun int) (*domain.TrenBulanan, error) {
	var dataTren []domain.DataBulanan
	var totalPemasukanTahun, totalPengeluaranTahun float64

	namaBulan := []string{
		"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	for bulan := 1; bulan <= 12; bulan++ {
		var totalPemasukan, totalPengeluaran float64

		err := r.db.Table("transaksis").
			Where("user_id = ? AND EXTRACT(year FROM tanggal) = ? AND EXTRACT(month FROM tanggal) = ? AND jenis = ?",
				userID, tahun, bulan, "Pemasukan").
			Select("COALESCE(SUM(jumlah), 0)").
			Row().
			Scan(&totalPemasukan)
		if err != nil {
			return nil, err
		}

		err = r.db.Table("transaksis").
			Where("user_id = ? AND EXTRACT(year FROM tanggal) = ? AND EXTRACT(month FROM tanggal) = ? AND jenis = ?",
				userID, tahun, bulan, "Pengeluaran").
			Select("COALESCE(SUM(jumlah), 0)").
			Row().
			Scan(&totalPengeluaran)
		if err != nil {
			return nil, err
		}

		dataTren = append(dataTren, domain.DataBulanan{
			Bulan:            bulan,
			NamaBulan:        namaBulan[bulan-1],
			TotalPemasukan:   totalPemasukan,
			TotalPengeluaran: totalPengeluaran,
		})

		totalPemasukanTahun += totalPemasukan
		totalPengeluaranTahun += totalPengeluaran
	}

	return &domain.TrenBulanan{
		Tahun:                 tahun,
		DataTren:              dataTren,
		TotalPemasukanTahun:   totalPemasukanTahun,
		TotalPengeluaranTahun: totalPengeluaranTahun,
	}, nil
}

func (r *laporanRepository) GetPerbandinganKantong(userID uint, bulanIni, tahunIni, bulanLalu, tahunLalu int) (*domain.PerbandinganKantong, error) {
	type KantongResult struct {
		KantongID       string  `gorm:"column:kantong_id"`
		KantongNama     string  `gorm:"column:kantong_nama"`
		JumlahBulanIni  float64 `gorm:"column:jumlah_bulan_ini"`
		JumlahBulanLalu float64 `gorm:"column:jumlah_bulan_lalu"`
	}

	var results []KantongResult

	query := `
		WITH kantong_bulan_ini AS (
			SELECT 
				k.id as kantong_id,
				k.nama as kantong_nama,
				COALESCE(SUM(t.jumlah), 0) as jumlah_bulan_ini
			FROM kantongs k
			LEFT JOIN transaksis t ON k.id = t.kantong_id 
				AND t.jenis = 'Pengeluaran'
				AND EXTRACT(year FROM t.tanggal) = ?
				AND EXTRACT(month FROM t.tanggal) = ?
			WHERE k.user_id = ?
			GROUP BY k.id, k.nama
		),
		kantong_bulan_lalu AS (
			SELECT 
				k.id as kantong_id,
				COALESCE(SUM(t.jumlah), 0) as jumlah_bulan_lalu
			FROM kantongs k
			LEFT JOIN transaksis t ON k.id = t.kantong_id 
				AND t.jenis = 'Pengeluaran'
				AND EXTRACT(year FROM t.tanggal) = ?
				AND EXTRACT(month FROM t.tanggal) = ?
			WHERE k.user_id = ?
			GROUP BY k.id
		)
		SELECT 
			kbi.kantong_id,
			kbi.kantong_nama,
			kbi.jumlah_bulan_ini,
			COALESCE(kbl.jumlah_bulan_lalu, 0) as jumlah_bulan_lalu
		FROM kantong_bulan_ini kbi
		LEFT JOIN kantong_bulan_lalu kbl ON kbi.kantong_id = kbl.kantong_id
		ORDER BY kbi.kantong_nama`

	err := r.db.Raw(query, tahunIni, bulanIni, userID, tahunLalu, bulanLalu, userID).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var dataKantong []domain.DataPerbandinganKantong
	var totalBulanIni, totalBulanLalu float64

	for _, result := range results {
		dataKantong = append(dataKantong, domain.DataPerbandinganKantong{
			KantongID:       result.KantongID,
			KantongNama:     result.KantongNama,
			JumlahBulanIni:  result.JumlahBulanIni,
			JumlahBulanLalu: result.JumlahBulanLalu,
		})
		totalBulanIni += result.JumlahBulanIni
		totalBulanLalu += result.JumlahBulanLalu
	}

	namaBulan := []string{
		"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember",
	}

	return &domain.PerbandinganKantong{
		BulanIni: domain.PeriodeBulan{
			Bulan:     bulanIni,
			NamaBulan: namaBulan[bulanIni],
			Tahun:     tahunIni,
		},
		BulanSebelumnya: domain.PeriodeBulan{
			Bulan:     bulanLalu,
			NamaBulan: namaBulan[bulanLalu],
			Tahun:     tahunLalu,
		},
		DataKantong:    dataKantong,
		TotalBulanIni:  totalBulanIni,
		TotalBulanLalu: totalBulanLalu,
	}, nil
}

func (r *laporanRepository) GetDetailPerbandinganKantong(userID uint, bulanIni, tahunIni, bulanLalu, tahunLalu int) (*domain.DetailPerbandinganKantong, error) {
	perbandinganKantong, err := r.GetPerbandinganKantong(userID, bulanIni, tahunIni, bulanLalu, tahunLalu)
	if err != nil {
		return nil, err
	}

	var dataKantong []domain.DataDetailPerbandinganKantong

	for _, kantong := range perbandinganKantong.DataKantong {
		rataRata := (kantong.JumlahBulanIni + kantong.JumlahBulanLalu) / 2

		var persentase float64
		var trend string

		if kantong.JumlahBulanLalu > 0 {
			persentase = ((kantong.JumlahBulanIni - kantong.JumlahBulanLalu) / kantong.JumlahBulanLalu) * 100
		} else if kantong.JumlahBulanIni > 0 {
			persentase = 100
		}

		if persentase > 0 {
			trend = "naik"
		} else if persentase < 0 {
			trend = "turun"
		} else {
			trend = "stabil"
		}

		dataKantong = append(dataKantong, domain.DataDetailPerbandinganKantong{
			KantongID:           kantong.KantongID,
			KantongNama:         kantong.KantongNama,
			JumlahBulanIni:      kantong.JumlahBulanIni,
			JumlahBulanLalu:     kantong.JumlahBulanLalu,
			RataRataPengeluaran: rataRata,
			Persentase:          persentase,
			Trend:               trend,
		})
	}

	rataRataTotal := (perbandinganKantong.TotalBulanIni + perbandinganKantong.TotalBulanLalu) / 2

	var persentaseTotal float64
	var trendTotal string

	if perbandinganKantong.TotalBulanLalu > 0 {
		persentaseTotal = ((perbandinganKantong.TotalBulanIni - perbandinganKantong.TotalBulanLalu) / perbandinganKantong.TotalBulanLalu) * 100
	} else if perbandinganKantong.TotalBulanIni > 0 {
		persentaseTotal = 100
	}

	if persentaseTotal > 0 {
		trendTotal = "naik"
	} else if persentaseTotal < 0 {
		trendTotal = "turun"
	} else {
		trendTotal = "stabil"
	}

	return &domain.DetailPerbandinganKantong{
		BulanIni:        perbandinganKantong.BulanIni,
		BulanSebelumnya: perbandinganKantong.BulanSebelumnya,
		DataKantong:     dataKantong,
		TotalBulanIni:   perbandinganKantong.TotalBulanIni,
		TotalBulanLalu:  perbandinganKantong.TotalBulanLalu,
		RataRataTotal:   rataRataTotal,
		PersentaseTotal: persentaseTotal,
		TrendTotal:      trendTotal,
	}, nil
}
