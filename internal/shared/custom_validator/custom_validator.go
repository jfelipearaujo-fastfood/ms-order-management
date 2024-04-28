package custom_validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

const (
	ValidatorUUIDWhenNotEmpty    = "uuid-when-not-empty"
	ValidatorTrackIDWhenNotEmpty = "track-id-when-not-empty"
)

func RegisterCustomValidations(validator *validator.Validate) error {
	var err error

	err = validator.RegisterValidation(ValidatorUUIDWhenNotEmpty, ValidateUUIDWhenNotEmpty, true)
	if err != nil {
		return err
	}

	err = validator.RegisterValidation(ValidatorTrackIDWhenNotEmpty, ValidateTrackIDWhenNotEmpty, true)
	if err != nil {
		return err
	}

	return nil
}

func ValidateUUIDWhenNotEmpty(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}

	if _, err := uuid.Parse(fl.Field().String()); err != nil {
		return false
	}

	return true
}

func ValidateTrackIDWhenNotEmpty(fl validator.FieldLevel) bool {
	if fl.Field().String() == "" {
		return true
	}

	pattern := "^[A-Z]{3}-[0-9]{3}$"

	regex := regexp.MustCompile(pattern)

	if ok := regex.MatchString(fl.Field().String()); !ok {
		return false
	}

	return true
}
