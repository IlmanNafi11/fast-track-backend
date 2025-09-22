package helper

import (
	"fiber-boiler-plate/internal/domain"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(data interface{}) []domain.ValidationError {
	var validationErrors []domain.ValidationError

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationError := domain.ValidationError{
				Field:   err.Field(),
				Message: getValidationMessage(err),
			}
			validationErrors = append(validationErrors, validationError)
		}
	}

	return validationErrors
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + " wajib diisi"
	case "email":
		return "Format email tidak valid"
	case "min":
		return err.Field() + " minimal " + err.Param() + " karakter"
	case "max":
		return err.Field() + " maksimal " + err.Param() + " karakter"
	default:
		return err.Field() + " tidak valid"
	}
}
