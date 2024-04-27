package create

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/errors"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil when dto is valid", func(t *testing.T) {
		// Arrange
		dto := CreateOrderDto{
			CustomerID: uuid.NewString(),
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when dto is invalid", func(t *testing.T) {
		// Arrange
		dto := CreateOrderDto{
			CustomerID: "abc",
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrRequestNotValid)
	})

	t.Run("Should return error when dto is empty", func(t *testing.T) {
		// Arrange
		dto := CreateOrderDto{}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, errors.ErrRequestNotValid)
	})
}
