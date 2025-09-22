package repo

import (
	"errors"
	"fiber-boiler-plate/internal/domain"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type transaksiRepository struct {
	db *gorm.DB
}

func NewTransaksiRepository(db *gorm.DB) TransaksiRepository {
	return &transaksiRepository{db: db}
}

func (r *transaksiRepository) GetByUserID(userID uint, req *domain.TransaksiListRequest) ([]*domain.TransaksiResponse, int, error) {
	var transaksiList []struct {
		domain.Transaksi
		KantongNama string `json:"kantong_nama"`
	}

	query := r.db.Table("transaksis t").
		Select("t.*, k.nama as kantong_nama").
		Joins("LEFT JOIN kantongs k ON t.kantong_id = k.id").
		Where("t.user_id = ?", userID)

	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		query = query.Where("k.nama ILIKE ? OR t.catatan ILIKE ?", searchTerm, searchTerm)
	}

	if req.Jenis != nil && *req.Jenis != "" {
		query = query.Where("t.jenis = ?", *req.Jenis)
	}

	if req.KantongNama != nil && *req.KantongNama != "" {
		query = query.Where("k.nama ILIKE ?", "%"+*req.KantongNama+"%")
	}

	if req.TanggalMulai != nil && *req.TanggalMulai != "" {
		query = query.Where("t.tanggal >= ?", *req.TanggalMulai)
	}

	if req.TanggalSelesai != nil && *req.TanggalSelesai != "" {
		query = query.Where("t.tanggal <= ?", *req.TanggalSelesai)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderBy := "t.tanggal"
	if req.SortBy == "jumlah" {
		orderBy = "t.jumlah"
	}

	direction := "DESC"
	if req.SortDirection == "asc" {
		direction = "ASC"
	}

	offset := (req.Page - 1) * req.PerPage
	if err := query.Order(fmt.Sprintf("%s %s", orderBy, direction)).
		Offset(offset).
		Limit(req.PerPage).
		Find(&transaksiList).Error; err != nil {
		return nil, 0, err
	}

	var result []*domain.TransaksiResponse
	for _, t := range transaksiList {
		result = append(result, &domain.TransaksiResponse{
			ID:          t.ID,
			Tanggal:     t.Tanggal.Format("2006-01-02"),
			Jenis:       t.Jenis,
			Jumlah:      t.Jumlah,
			KantongID:   t.KantongID,
			KantongNama: t.KantongNama,
			Catatan:     t.Catatan,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
		})
	}

	return result, int(total), nil
}

func (r *transaksiRepository) GetByID(id string, userID uint) (*domain.TransaksiResponse, error) {
	var transaksi struct {
		domain.Transaksi
		KantongNama string `json:"kantong_nama"`
	}

	err := r.db.Table("transaksis t").
		Select("t.*, k.nama as kantong_nama").
		Joins("LEFT JOIN kantongs k ON t.kantong_id = k.id").
		Where("t.id = ? AND t.user_id = ?", id, userID).
		First(&transaksi).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaksi tidak ditemukan")
		}
		return nil, err
	}

	return &domain.TransaksiResponse{
		ID:          transaksi.ID,
		Tanggal:     transaksi.Tanggal.Format("2006-01-02"),
		Jenis:       transaksi.Jenis,
		Jumlah:      transaksi.Jumlah,
		KantongID:   transaksi.KantongID,
		KantongNama: transaksi.KantongNama,
		Catatan:     transaksi.Catatan,
		CreatedAt:   transaksi.CreatedAt,
		UpdatedAt:   transaksi.UpdatedAt,
	}, nil
}

func (r *transaksiRepository) Create(transaksi *domain.Transaksi) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var kantong domain.Kantong
		if err := tx.Where("id = ? AND user_id = ?", transaksi.KantongID, transaksi.UserID).First(&kantong).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("kantong tidak ditemukan")
			}
			return err
		}

		if err := tx.Create(transaksi).Error; err != nil {
			return err
		}

		if transaksi.Jenis == "Pemasukan" {
			kantong.Saldo += transaksi.Jumlah
		} else if transaksi.Jenis == "Pengeluaran" {
			if kantong.Saldo < transaksi.Jumlah {
				return errors.New("saldo tidak mencukupi")
			}
			kantong.Saldo -= transaksi.Jumlah
		}

		return tx.Save(&kantong).Error
	})
}

func (r *transaksiRepository) Update(transaksi *domain.Transaksi) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingTransaksi domain.Transaksi
		if err := tx.Where("id = ? AND user_id = ?", transaksi.ID, transaksi.UserID).First(&existingTransaksi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("transaksi tidak ditemukan")
			}
			return err
		}

		var oldKantong, newKantong domain.Kantong
		if err := tx.Where("id = ?", existingTransaksi.KantongID).First(&oldKantong).Error; err != nil {
			return err
		}

		if err := tx.Where("id = ? AND user_id = ?", transaksi.KantongID, transaksi.UserID).First(&newKantong).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("kantong tujuan tidak ditemukan")
			}
			return err
		}

		if existingTransaksi.Jenis == "Pemasukan" {
			oldKantong.Saldo -= existingTransaksi.Jumlah
		} else {
			oldKantong.Saldo += existingTransaksi.Jumlah
		}

		if transaksi.Jenis == "Pemasukan" {
			newKantong.Saldo += transaksi.Jumlah
		} else {
			if newKantong.Saldo < transaksi.Jumlah {
				return errors.New("saldo tidak mencukupi")
			}
			newKantong.Saldo -= transaksi.Jumlah
		}

		if err := tx.Save(&oldKantong).Error; err != nil {
			return err
		}

		if oldKantong.ID != newKantong.ID {
			if err := tx.Save(&newKantong).Error; err != nil {
				return err
			}
		}

		transaksi.UpdatedAt = time.Now()
		return tx.Model(&existingTransaksi).Updates(transaksi).Error
	})
}

func (r *transaksiRepository) Delete(id string, userID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var transaksi domain.Transaksi
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&transaksi).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("transaksi tidak ditemukan")
			}
			return err
		}

		var kantong domain.Kantong
		if err := tx.Where("id = ?", transaksi.KantongID).First(&kantong).Error; err != nil {
			return err
		}

		if transaksi.Jenis == "Pemasukan" {
			kantong.Saldo -= transaksi.Jumlah
		} else {
			kantong.Saldo += transaksi.Jumlah
		}

		if err := tx.Save(&kantong).Error; err != nil {
			return err
		}

		return tx.Delete(&transaksi).Error
	})
}
