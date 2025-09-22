package helper_test

import (
	"fiber-boiler-plate/internal/domain"
	"fiber-boiler-plate/internal/helper"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateStruct_ValidData(t *testing.T) {
	req := domain.RegisterRequest{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
	}

	errors := helper.ValidateStruct(req)

	assert.Empty(t, errors)
}

func TestValidateStruct_InvalidData(t *testing.T) {
	req := domain.RegisterRequest{
		Name:     "",
		Email:    "invalid-email",
		Password: "123",
	}

	errors := helper.ValidateStruct(req)

	assert.NotEmpty(t, errors)
	assert.Greater(t, len(errors), 0)

	for _, validationError := range errors {
		assert.NotEmpty(t, validationError.Field)
		assert.NotEmpty(t, validationError.Message)
	}
}

func TestValidateStruct_EmptyEmail(t *testing.T) {
	req := domain.AuthRequest{
		Email:    "",
		Password: "password123",
	}

	errors := helper.ValidateStruct(req)

	assert.NotEmpty(t, errors)
	assert.Greater(t, len(errors), 0)
}

func TestValidateStruct_ShortPassword(t *testing.T) {
	req := domain.AuthRequest{
		Email:    "test@example.com",
		Password: "123",
	}

	errors := helper.ValidateStruct(req)

	assert.NotEmpty(t, errors)
	assert.Greater(t, len(errors), 0)
}

func TestValidateStruct_ResetPasswordRequest(t *testing.T) {
	validReq := domain.ResetPasswordRequest{
		Email: "test@example.com",
	}

	validErrors := helper.ValidateStruct(validReq)
	assert.Empty(t, validErrors)

	invalidReq := domain.ResetPasswordRequest{
		Email: "invalid-email",
	}

	invalidErrors := helper.ValidateStruct(invalidReq)
	assert.NotEmpty(t, invalidErrors)
}

func TestValidateStruct_RefreshTokenRequest(t *testing.T) {
	validReq := domain.RefreshTokenRequest{
		RefreshToken: "valid_token_123",
	}

	validErrors := helper.ValidateStruct(validReq)
	assert.Empty(t, validErrors)

	invalidReq := domain.RefreshTokenRequest{
		RefreshToken: "",
	}

	invalidErrors := helper.ValidateStruct(invalidReq)
	assert.NotEmpty(t, invalidErrors)
}

func TestValidateStruct_NewPasswordRequest(t *testing.T) {
	validReq := domain.NewPasswordRequest{
		Token:       "valid_token",
		NewPassword: "newpassword123",
	}

	validErrors := helper.ValidateStruct(validReq)
	assert.Empty(t, validErrors)

	invalidReq := domain.NewPasswordRequest{
		Token:       "",
		NewPassword: "123",
	}

	invalidErrors := helper.ValidateStruct(invalidReq)
	assert.NotEmpty(t, invalidErrors)
	assert.Greater(t, len(invalidErrors), 1)
}
