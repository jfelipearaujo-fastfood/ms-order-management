package get_all

import (
	"testing"

	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return error when state is invalid", func(t *testing.T) {
		// Arrange
		dto := GetOrdersDto{
			State: 10,
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrRequestNotValid)
	})

	t.Run("Should return nil when state is valid", func(t *testing.T) {
		// Arrange
		dto := GetOrdersDto{
			State: 3,
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.NoError(t, err)
	})
}
