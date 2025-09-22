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
	FlushAll() error
	Ping() error
}
