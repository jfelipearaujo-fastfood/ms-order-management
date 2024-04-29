package send_to_pay

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	t.Run("Should return nil when request is valid", func(t *testing.T) {
		// Arrange
		dto := SendToPayDto{
			OrderID: uuid.NewString(),
		}

		// Act
		err := dto.Validate()

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return error when request is not valid", func(t *testing.T) {
		// Arrange
		dto := SendToPayDto{}

		// Act
		err := dto.Validate()

		// Assert
		assert.Error(t, err)
	})
}
