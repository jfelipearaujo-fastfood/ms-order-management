package custom_validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateUUIDWhenNotEmpty(t *testing.T) {
	t.Run("Should return nil when UUID is valid", func(t *testing.T) {
		// Arrange
		type test struct {
			Id string `validate:"uuid-when-not-empty"`
		}

		uuid := uuid.NewString()

		validator := validator.New()
		err := RegisterCustomValidations(validator)
		assert.NoError(t, err)

		data := test{
			Id: uuid,
		}

		// Act
		err = validator.Struct(data)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return nil when UUID is empty", func(t *testing.T) {
		// Arrange
		type test struct {
			Id string `validate:"uuid-when-not-empty"`
		}

		validator := validator.New()
		err := RegisterCustomValidations(validator)
		assert.NoError(t, err)

		data := test{}

		// Act
		err = validator.Struct(data)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when UUID is invalid", func(t *testing.T) {
		// Arrange
		type test struct {
			Id string `validate:"uuid-when-not-empty"`
		}

		validator := validator.New()
		err := RegisterCustomValidations(validator)
		assert.NoError(t, err)

		data := test{
			Id: "invalid-uuid",
		}

		// Act
		err = validator.Struct(data)

		// Assert
		assert.Error(t, err)
	})
}

func TestValidateTrackIDWhenNotEmpty(t *testing.T) {
	t.Run("Should return nil when TrackID is valid", func(t *testing.T) {
		// Arrange
		type test struct {
			TrackID string `validate:"track-id-when-not-empty"`
		}

		validator := validator.New()
		err := RegisterCustomValidations(validator)
		assert.NoError(t, err)

		data := test{
			TrackID: "ABC123",
		}

		// Act
		err = validator.Struct(data)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return nil when TrackID is empty", func(t *testing.T) {
		// Arrange
		type test struct {
			TrackID string `validate:"track-id-when-not-empty"`
		}

		validator := validator.New()
		err := RegisterCustomValidations(validator)
		assert.NoError(t, err)

		data := test{}

		// Act
		err = validator.Struct(data)

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when TrackID is invalid", func(t *testing.T) {
		// Arrange
		type test struct {
			TrackID string `validate:"track-id-when-not-empty"`
		}

		validator := validator.New()
		err := RegisterCustomValidations(validator)
		assert.NoError(t, err)

		data := test{
			TrackID: "invalid-track-id",
		}

		// Act
		err = validator.Struct(data)

		// Assert
		assert.Error(t, err)
	})
}
