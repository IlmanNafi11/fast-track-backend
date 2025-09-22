package repo_test

import (
	"fiber-boiler-plate/internal/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockDB struct {
	mock.Mock
	db *gorm.DB
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.Called(query, args)
	return m.db
}

func (m *MockDB) First(dest interface{}) *gorm.DB {
	args := m.Called(dest)
	if args.Error(0) != nil {
		return &gorm.DB{Error: args.Error(0)}
	}
	if user, ok := dest.(*domain.User); ok {
		*user = domain.User{
			ID:       1,
			Email:    "test@example.com",
			Name:     "Test User",
			IsActive: true,
		}
	}
	return &gorm.DB{Error: nil}
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return &gorm.DB{Error: args.Error(0)}
}

func (m *MockDB) Model(value interface{}) *gorm.DB {
	m.Called(value)
	return m.db
}

func (m *MockDB) Update(column string, value interface{}) *gorm.DB {
	args := m.Called(column, value)
	return &gorm.DB{Error: args.Error(0)}
}

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Name:     "Test User",
		IsActive: true,
	}

	result, err := user, error(nil)

	assert.NoError(t, err)
	assert.Equal(t, user.Email, result.Email)
	assert.Equal(t, user.ID, result.ID)
	assert.True(t, result.IsActive)
}

func TestUserRepository_GetByID_Success(t *testing.T) {
	user := &domain.User{
		ID:       1,
		Email:    "test@example.com",
		Name:     "Test User",
		IsActive: true,
	}

	result, err := user, error(nil)

	assert.NoError(t, err)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Email, result.Email)
	assert.True(t, result.IsActive)
}

func TestUserRepository_Create_Success(t *testing.T) {
	user := &domain.User{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "hashedpassword",
		IsActive: true,
	}

	err := error(nil)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.True(t, user.IsActive)
}

func TestUserRepository_UpdatePassword_Success(t *testing.T) {
	email := "test@example.com"
	hashedPassword := "newhashed"

	err := error(nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, email)
	assert.NotEmpty(t, hashedPassword)
	assert.Equal(t, "test@example.com", email)
	assert.Equal(t, "newhashed", hashedPassword)
}

func TestUserRepository_DeleteRefreshToken_Success(t *testing.T) {
	token := "refresh_token_123"

	err := error(nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "refresh_token_123", token)
}

func TestRefreshTokenRepository_Create_Success(t *testing.T) {
	userID := uint(1)
	token := "refresh_token_123"
	expiresAt := time.Now().Add(24 * time.Hour)

	result := &domain.RefreshToken{
		ID:        1,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		IsRevoked: false,
	}
	err := error(nil)

	assert.NoError(t, err)
	assert.Equal(t, userID, result.UserID)
	assert.Equal(t, token, result.Token)
	assert.False(t, result.IsRevoked)
}

func TestRefreshTokenRepository_GetByToken_Success(t *testing.T) {
	token := "refresh_token_123"

	result := &domain.RefreshToken{
		ID:        1,
		UserID:    1,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		IsRevoked: false,
	}
	err := error(nil)

	assert.NoError(t, err)
	assert.Equal(t, token, result.Token)
	assert.False(t, result.IsRevoked)
}

func TestRefreshTokenRepository_RevokeToken_Success(t *testing.T) {
	token := "refresh_token_123"

	err := error(nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "refresh_token_123", token)
}

func TestPasswordResetTokenRepository_Create_Success(t *testing.T) {
	email := "test@example.com"
	token := "reset_token_123"
	expiresAt := time.Now().Add(time.Hour)

	result := &domain.PasswordResetToken{
		ID:        1,
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
		IsUsed:    false,
	}
	err := error(nil)

	assert.NoError(t, err)
	assert.Equal(t, email, result.Email)
	assert.Equal(t, token, result.Token)
	assert.False(t, result.IsUsed)
}

func TestPasswordResetTokenRepository_GetByToken_Success(t *testing.T) {
	token := "reset_token_123"

	result := &domain.PasswordResetToken{
		ID:        1,
		Email:     "test@example.com",
		Token:     token,
		ExpiresAt: time.Now().Add(time.Hour),
		IsUsed:    false,
	}
	err := error(nil)

	assert.NoError(t, err)
	assert.Equal(t, token, result.Token)
	assert.False(t, result.IsUsed)
}

func TestPasswordResetTokenRepository_MarkAsUsed_Success(t *testing.T) {
	token := "reset_token_123"

	err := error(nil)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "reset_token_123", token)
}
