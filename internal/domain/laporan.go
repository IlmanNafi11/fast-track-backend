package domain

import "time"

type RingkasanLaporan struct {
	TotalPemasukan            float64        `json:"total_pemasukan"`
	TotalPengeluaran          float64        `json:"total_pengeluaran"`
	TotalSaldo                float64        `json:"total_saldo"`
	RataRataPengeluaranHarian float64        `json:"rata_rata_pengeluaran_harian"`
	Periode                   PeriodeTanggal `json:"periode"`
}

type PeriodeTanggal struct {
	TanggalMulai   string `json:"tanggal_mulai"`
	TanggalSelesai string `json:"tanggal_selesai"`
}

type StatistikTahunan struct {
	Tahun                 int           `json:"tahun"`
	DataBulanan           []DataBulanan `json:"data_bulanan"`
	TotalPemasukanTahun   float64       `json:"total_pemasukan_tahun"`
	TotalPengeluaranTahun float64       `json:"total_pengeluaran_tahun"`
}

type DataBulanan struct {
	Bulan            int     `json:"bulan"`
	NamaBulan        string  `json:"nama_bulan"`
	TotalPemasukan   float64 `json:"total_pemasukan"`
	TotalPengeluaran float64 `json:"total_pengeluaran"`
}

type StatistikKantongBulanan struct {
	Periode          PeriodeBulan         `json:"periode"`
	DataKantong      []DataKantongBulanan `json:"data_kantong"`
	TotalPengeluaran float64              `json:"total_pengeluaran"`
	TotalTransaksi   int                  `json:"total_transaksi"`
}

type PeriodeBulan struct {
	Bulan     int    `json:"bulan"`
	NamaBulan string `json:"nama_bulan"`
	Tahun     int    `json:"tahun"`
}

type DataKantongBulanan struct {
	KantongID        string  `json:"kantong_id"`
	KantongNama      string  `json:"kantong_nama"`
	Kategori         string  `json:"kategori"`
	TotalPengeluaran float64 `json:"total_pengeluaran"`
	JumlahTransaksi  int     `json:"jumlah_transaksi"`
	Persentase       float64 `json:"persentase"`
}

type TopKantongPengeluaran struct {
	Periode               PeriodeBulan     `json:"periode"`
	TopKantong            []DataTopKantong `json:"top_kantong"`
	TotalPengeluaranSemua float64          `json:"total_pengeluaran_semua"`
	TotalTransaksiSemua   int              `json:"total_transaksi_semua"`
}

type DataTopKantong struct {
	Ranking             int     `json:"ranking"`
	KantongID           string  `json:"kantong_id"`
	KantongNama         string  `json:"kantong_nama"`
	Kategori            string  `json:"kategori"`
	TotalPengeluaran    float64 `json:"total_pengeluaran"`
	JumlahTransaksi     int     `json:"jumlah_transaksi"`
	PersentaseDariTotal float64 `json:"persentase_dari_total"`
	RataRataPengeluaran float64 `json:"rata_rata_pengeluaran"`
}

type RingkasanLaporanRequest struct {
	TanggalMulai   *string `json:"tanggal_mulai" query:"tanggal_mulai"`
	TanggalSelesai *string `json:"tanggal_selesai" query:"tanggal_selesai"`
}

type StatistikTahunanRequest struct {
	Tahun *int `json:"tahun" query:"tahun" validate:"omitempty,min=2020,max=2030"`
}

type StatistikKantongBulananRequest struct {
	Bulan *int `json:"bulan" query:"bulan" validate:"omitempty,min=1,max=12"`
	Tahun *int `json:"tahun" query:"tahun" validate:"omitempty,min=2020,max=2030"`
}

type TopKantongPengeluaranRequest struct {
	Bulan *int `json:"bulan" query:"bulan" validate:"omitempty,min=1,max=12"`
	Tahun *int `json:"tahun" query:"tahun" validate:"omitempty,min=2020,max=2030"`
	Limit *int `json:"limit" query:"limit" validate:"omitempty,min=1,max=10"`
}

type RingkasanLaporanResponse struct {
	Success   bool             `json:"success"`
	Message   string           `json:"message"`
	Code      int              `json:"code"`
	Data      RingkasanLaporan `json:"data"`
	Timestamp time.Time        `json:"timestamp"`
}

type StatistikTahunanResponse struct {
	Success   bool             `json:"success"`
	Message   string           `json:"message"`
	Code      int              `json:"code"`
	Data      StatistikTahunan `json:"data"`
	Timestamp time.Time        `json:"timestamp"`
}

type StatistikKantongBulananResponse struct {
	Success   bool                    `json:"success"`
	Message   string                  `json:"message"`
	Code      int                     `json:"code"`
	Data      StatistikKantongBulanan `json:"data"`
	Timestamp time.Time               `json:"timestamp"`
}

type TopKantongPengeluaranResponse struct {
	Success   bool                  `json:"success"`
	Message   string                `json:"message"`
	Code      int                   `json:"code"`
	Data      TopKantongPengeluaran `json:"data"`
	Timestamp time.Time             `json:"timestamp"`
}

type StatistikKantongPeriode struct {
	Periode          PeriodeTanggal       `json:"periode"`
	DataKantong      []DataKantongPeriode `json:"data_kantong"`
	TotalPengeluaran float64              `json:"total_pengeluaran"`
}

type DataKantongPeriode struct {
	KantongID        string  `json:"kantong_id"`
	KantongNama      string  `json:"kantong_nama"`
	TotalPengeluaran float64 `json:"total_pengeluaran"`
}

type PengeluaranKantongDetail struct {
	Periode                PeriodeTanggal      `json:"periode"`
	DataKantong            []DataKantongDetail `json:"data_kantong"`
	TotalPengeluaran       float64             `json:"total_pengeluaran"`
	TotalSaldoSemuaKantong float64             `json:"total_saldo_semua_kantong"`
}

type DataKantongDetail struct {
	KantongID           string  `json:"kantong_id"`
	KantongNama         string  `json:"kantong_nama"`
	TotalPengeluaran    float64 `json:"total_pengeluaran"`
	PersentaseDariSaldo float64 `json:"persentase_dari_saldo"`
	JumlahTransaksi     int     `json:"jumlah_transaksi"`
	RataRataPengeluaran float64 `json:"rata_rata_pengeluaran"`
	SaldoKantong        float64 `json:"saldo_kantong"`
}

type StatistikKantongPeriodeRequest struct {
	TanggalMulai   *string `json:"tanggal_mulai" query:"tanggal_mulai"`
	TanggalSelesai *string `json:"tanggal_selesai" query:"tanggal_selesai"`
}

type PengeluaranKantongDetailRequest struct {
	TanggalMulai   *string `json:"tanggal_mulai" query:"tanggal_mulai"`
	TanggalSelesai *string `json:"tanggal_selesai" query:"tanggal_selesai"`
}

type StatistikKantongPeriodeResponse struct {
	Success   bool                    `json:"success"`
	Message   string                  `json:"message"`
	Code      int                     `json:"code"`
	Data      StatistikKantongPeriode `json:"data"`
	Timestamp time.Time               `json:"timestamp"`
}

type PengeluaranKantongDetailResponse struct {
	Success   bool                     `json:"success"`
	Message   string                   `json:"message"`
	Code      int                      `json:"code"`
	Data      PengeluaranKantongDetail `json:"data"`
	Timestamp time.Time                `json:"timestamp"`
}

type TrenBulanan struct {
	Tahun                 int           `json:"tahun"`
	DataTren              []DataBulanan `json:"data_tren"`
	TotalPemasukanTahun   float64       `json:"total_pemasukan_tahun"`
	TotalPengeluaranTahun float64       `json:"total_pengeluaran_tahun"`
}

type TrenBulananRequest struct {
	Tahun *int `json:"tahun" query:"tahun" validate:"omitempty,min=2020,max=2030"`
}

type TrenBulananResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Code      int         `json:"code"`
	Data      TrenBulanan `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
}

type PerbandinganKantong struct {
	BulanIni        PeriodeBulan              `json:"bulan_ini"`
	BulanSebelumnya PeriodeBulan              `json:"bulan_sebelumnya"`
	DataKantong     []DataPerbandinganKantong `json:"data_kantong"`
	TotalBulanIni   float64                   `json:"total_bulan_ini"`
	TotalBulanLalu  float64                   `json:"total_bulan_lalu"`
}

type DataPerbandinganKantong struct {
	KantongID       string  `json:"kantong_id"`
	KantongNama     string  `json:"kantong_nama"`
	JumlahBulanIni  float64 `json:"jumlah_bulan_ini"`
	JumlahBulanLalu float64 `json:"jumlah_bulan_lalu"`
}

type PerbandinganKantongResponse struct {
	Success   bool                `json:"success"`
	Message   string              `json:"message"`
	Code      int                 `json:"code"`
	Data      PerbandinganKantong `json:"data"`
	Timestamp time.Time           `json:"timestamp"`
}

type DetailPerbandinganKantong struct {
	BulanIni        PeriodeBulan                    `json:"bulan_ini"`
	BulanSebelumnya PeriodeBulan                    `json:"bulan_sebelumnya"`
	DataKantong     []DataDetailPerbandinganKantong `json:"data_kantong"`
	TotalBulanIni   float64                         `json:"total_bulan_ini"`
	TotalBulanLalu  float64                         `json:"total_bulan_lalu"`
	RataRataTotal   float64                         `json:"rata_rata_total"`
	PersentaseTotal float64                         `json:"persentase_total"`
	TrendTotal      string                          `json:"trend_total"`
}

type DataDetailPerbandinganKantong struct {
	KantongID           string  `json:"kantong_id"`
	KantongNama         string  `json:"kantong_nama"`
	JumlahBulanIni      float64 `json:"jumlah_bulan_ini"`
	JumlahBulanLalu     float64 `json:"jumlah_bulan_lalu"`
	RataRataPengeluaran float64 `json:"rata_rata_pengeluaran"`
	Persentase          float64 `json:"persentase"`
	Trend               string  `json:"trend"`
}

type DetailPerbandinganKantongResponse struct {
	Success   bool                      `json:"success"`
	Message   string                    `json:"message"`
	Code      int                       `json:"code"`
	Data      DetailPerbandinganKantong `json:"data"`
	Timestamp time.Time                 `json:"timestamp"`
}
