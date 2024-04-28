package get

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
)

func TestFindViaID(t *testing.T) {
	t.Run("Should return true when UUID is not empty", func(t *testing.T) {
		// Arrange
		dto := GetOrderDto{
			OrderId: "valid-uuid",
		}

		// Act
		result := dto.FindViaID()

		// Assert
		assert.True(t, result)
	})

	t.Run("Should return false when UUID is empty", func(t *testing.T) {
		// Arrange
		dto := GetOrderDto{}

		// Act
		result := dto.FindViaID()

		// Assert
		assert.False(t, result)
	})
}

func TestValidate(t *testing.T) {
	t.Run("Should return nil when dto is valid with id", func(t *testing.T) {
		// Arrange
		dto := GetOrderDto{
			OrderId: uuid.NewString(),
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return nil when dto is valid with trackId", func(t *testing.T) {
		// Arrange
		dto := GetOrderDto{
			TrackId: "ABC456",
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when dto is invalid with id", func(t *testing.T) {
		// Arrange
		dto := GetOrderDto{
			OrderId: "invalid-uuid",
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrRequestNotValid)
	})

	t.Run("Should return error when dto is invalid with trackId", func(t *testing.T) {
		// Arrange
		dto := GetOrderDto{
			TrackId: "invalid-track-id",
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrRequestNotValid)
	})

	t.Run("Should return error when dto is invalid with id and trackId", func(t *testing.T) {
		// Arrange
		dto := GetOrderDto{
			OrderId: "invalid-uuid",
			TrackId: "invalid-track-id",
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrRequestNotValid)
	})

	t.Run("Should return error when dto is empty", func(t *testing.T) {
		// Arrange
		dto := GetOrderDto{}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrRequestNotValid)
	})
}
