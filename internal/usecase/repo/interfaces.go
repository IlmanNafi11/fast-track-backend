package repo

import (
	"fiber-boiler-plate/internal/domain"
	"time"
)

type UserRepository interface {
	GetByEmail(email string) (*domain.User, error)
	GetByID(id uint) (*domain.User, error)
	Create(user *domain.User) error
	UpdatePassword(email, hashedPassword string) error
}

type RefreshTokenRepository interface {
	Create(userID uint, token string, expiresAt time.Time) (*domain.RefreshToken, error)
	GetByToken(token string) (*domain.RefreshToken, error)
	RevokeToken(token string) error
	RevokeAllUserTokens(userID uint) error
	CleanupExpired() error
}

type PasswordResetTokenRepository interface {
	Create(email, token string, expiresAt time.Time) (*domain.PasswordResetToken, error)
	GetByToken(token string) (*domain.PasswordResetToken, error)
	MarkAsUsed(token string) error
	CleanupExpired() error
}

type KantongRepository interface {
	GetByUserID(userID uint, req *domain.KantongListRequest) ([]*domain.Kantong, int, error)
	GetByID(id string, userID uint) (*domain.Kantong, error)
	GetByIDKartu(idKartu string, userID uint) (*domain.Kantong, error)
	Create(kantong *domain.Kantong) error
	Update(kantong *domain.Kantong) error
	Delete(id string, userID uint) error
	IsNameExistForUser(nama string, userID uint, excludeID ...string) (bool, error)
	GenerateUniqueIDKartu() (string, error)
	Transfer(kantongAsalID, kantongTujuanID string, jumlah float64, userID uint) (*domain.Kantong, *domain.Kantong, error)
}

type TransaksiRepository interface {
	GetByUserID(userID uint, req *domain.TransaksiListRequest) ([]*domain.TransaksiResponse, int, error)
	GetByID(id string, userID uint) (*domain.TransaksiResponse, error)
	Create(transaksi *domain.Transaksi) error
	Update(transaksi *domain.Transaksi) error
	Delete(id string, userID uint) error
}

type AnggaranRepository interface {
	GetByUserID(userID uint, req *domain.AnggaranListRequest) ([]*domain.AnggaranItem, int, error)
	GetByKantongID(kantongID string, userID uint, bulan, tahun int) (*domain.AnggaranItem, error)
	CreateOrUpdate(anggaran *domain.AnggaranItem) error
	CreatePenyesuaian(userID uint, req *domain.PenyesuaianAnggaranRequest) (*domain.AnggaranItem, error)
	GetStatistikBulan(kantongID string, userID uint, bulan, tahun int) ([]domain.StatistikHarian, error)
	RecalculateAnggaran(kantongID string, userID uint, bulan, tahun int) (*domain.AnggaranItem, error)
	CreateAnggaranForKantong(kantong *domain.Kantong) error
	UpdateAnggaranAfterTransaksi(kantongID string, userID uint) error
}

type LaporanRepository interface {
	GetRingkasanLaporan(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.RingkasanLaporan, error)
	GetStatistikTahunan(userID uint, tahun int) (*domain.StatistikTahunan, error)
	GetStatistikKantongBulanan(userID uint, bulan, tahun int) (*domain.StatistikKantongBulanan, error)
	GetTopKantongPengeluaran(userID uint, bulan, tahun, limit int) (*domain.TopKantongPengeluaran, error)
	GetStatistikKantongPeriode(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.StatistikKantongPeriode, error)
	GetPengeluaranKantongDetail(userID uint, tanggalMulai, tanggalSelesai time.Time) (*domain.PengeluaranKantongDetail, error)
	GetTrenBulanan(userID uint, tahun int) (*domain.TrenBulanan, error)
	GetPerbandinganKantong(userID uint, bulanIni, tahunIni, bulanLalu, tahunLalu int) (*domain.PerbandinganKantong, error)
	GetDetailPerbandinganKantong(userID uint, bulanIni, tahunIni, bulanLalu, tahunLalu int) (*domain.DetailPerbandinganKantong, error)
}

type SubscriptionPlanRepository interface {
	GetAll(req *domain.SubscriptionPlanListRequest) ([]*domain.SubscriptionPlan, int, error)
	GetByID(id string) (*domain.SubscriptionPlan, error)
	GetByKode(kode string) (*domain.SubscriptionPlan, error)
	Create(plan *domain.SubscriptionPlan) error
	Update(plan *domain.SubscriptionPlan) error
	Delete(id string) error
	IsNameExists(nama string, excludeID ...string) (bool, error)
	IsKodeExists(kode string, excludeID ...string) (bool, error)
	CountActiveUsers(planID string) (int64, error)
}

type PermissionRepository interface {
	GetAll(req *domain.PermissionListRequest) ([]*domain.Permission, int, error)
	GetByID(id string) (*domain.Permission, error)
	GetByNama(nama string) (*domain.Permission, error)
	Create(permission *domain.Permission) error
	Update(permission *domain.Permission) error
	Delete(id string) error
	IsNameExists(nama string, excludeID ...string) (bool, error)
	IsUsedByRoles(id string) (bool, error)
}

type RoleRepository interface {
	GetAll(req *domain.RoleListRequest) ([]*domain.Role, int, error)
	GetByID(id string) (*domain.Role, error)
	GetByNama(nama string) (*domain.Role, error)
	Create(role *domain.Role) error
	CreateWithPermissions(role *domain.Role, permissionIDs []string) error
	Update(role *domain.Role) error
	Delete(id string) error
	IsNameExists(nama string, excludeID ...string) (bool, error)
	IsUsedByUsers(id string) (bool, error)
	GetRolePermissions(roleID string, req *domain.RolePermissionListRequest) ([]*domain.Permission, int, error)
	UpdateRolePermissions(roleID string, permissionIDs []string) error
	ValidatePermissions(permissionIDs []string) ([]string, error)
}

type RedisRepository interface {
	Set(key string, value interface{}, ttl time.Duration) error
	Get(key string) (string, error)
	GetJSON(key string, dest interface{}) error
	SetJSON(key string, value interface{}, ttl time.Duration) error
	Delete(key string) error
	Exists(key string) (bool, error)
	Increment(key string) (int64, error)
	Decrement(key string) (int64, error)
	SetExpire(key string, ttl time.Duration) error
	GetTTL(key string) (time.Duration, error)
	GetKeys(pattern string) ([]string, error)
	FlushAll() error
	Ping() error
}
