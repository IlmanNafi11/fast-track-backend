package repo

import (
	"fiber-boiler-plate/internal/domain"
	"time"

	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) Create(email, token string, expiresAt time.Time) (*domain.PasswordResetToken, error) {
	resetToken := &domain.PasswordResetToken{
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	err := r.db.Create(resetToken).Error
	if err != nil {
		return nil, err
	}
	return resetToken, nil
}

func (r *passwordResetTokenRepository) GetByToken(token string) (*domain.PasswordResetToken, error) {
	var resetToken domain.PasswordResetToken
	err := r.db.Where("token = ? AND is_used = ? AND expires_at > ?", token, false, time.Now()).First(&resetToken).Error
	if err != nil {
		return nil, err
	}
	return &resetToken, nil
}

func (r *passwordResetTokenRepository) MarkAsUsed(token string) error {
	return r.db.Model(&domain.PasswordResetToken{}).Where("token = ?", token).Update("is_used", true).Error
}

func (r *passwordResetTokenRepository) CleanupExpired() error {
	return r.db.Where("expires_at < ? OR is_used = ?", time.Now(), true).Delete(&domain.PasswordResetToken{}).Error
}
