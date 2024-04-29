package payment_entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	t.Run("Should return the string representation of the state", func(t *testing.T) {
		// Arrange
		cases := []struct {
			state    PaymentState
			expected string
		}{
			{None, "None"},
			{WaitingForApproval, "WaitingForApproval"},
			{Approved, "Approved"},
			{Rejected, "Rejected"},
			{PaymentState(100), "Unknown"},
		}

		for _, c := range cases {
			// Act
			res := c.state.String()

			// Assert
			assert.Equal(t, c.expected, res)
		}
	})
}
