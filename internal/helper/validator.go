package helper

import (
	"fiber-boiler-plate/internal/domain"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("permissionname", validatePermissionName)
	validate.RegisterValidation("rolename", validateRoleName)
}

func validatePermissionName(fl validator.FieldLevel) bool {
	permissionName := fl.Field().String()
	match, _ := regexp.MatchString("^[a-z0-9._-]+$", permissionName)
	return match
}

func validateRoleName(fl validator.FieldLevel) bool {
	roleName := fl.Field().String()
	match, _ := regexp.MatchString("^[a-z0-9_-]+$", roleName)
	return match
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
	case "oneof":
		return err.Field() + " harus berupa salah satu dari: " + err.Param()
	case "uuid":
		return err.Field() + " harus berupa UUID yang valid"
	case "permissionname":
		return err.Field() + " hanya boleh mengandung huruf kecil, angka, titik, underscore, dan dash"
	case "rolename":
		return err.Field() + " hanya boleh mengandung huruf kecil, angka, underscore, dan dash"
	case "dive":
		return err.Field() + " mengandung data yang tidak valid"
	default:
		return err.Field() + " tidak valid"
	}
}
