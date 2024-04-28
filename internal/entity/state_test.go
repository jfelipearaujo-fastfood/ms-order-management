package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCanTransitionTo(t *testing.T) {
	t.Run("Should return true when transition is allowed", func(t *testing.T) {
		// Arrange
		cases := []struct {
			from State
			to   State
		}{
			{None, Created},
			{Created, Received},
			{Created, Cancelled},
			{Received, Processing},
			{Received, Cancelled},
			{Processing, Completed},
			{Processing, Cancelled},
			{Completed, Delivered},
		}

		// Act
		for _, c := range cases {
			res := c.from.CanTransitionTo(c.to)

			// Assert
			assert.True(t, res)
		}
	})

	t.Run("Should return false when transition is not allowed", func(t *testing.T) {
		// Arrange
		cases := []struct {
			from State
			to   State
		}{
			{None, Received},
			{Created, Processing},
			{Received, Completed},
			{Processing, Received},
			{Completed, Created},
		}

		// Act
		for _, c := range cases {
			res := c.from.CanTransitionTo(c.to)

			// Assert
			assert.False(t, res)
		}
	})
}

func TestIsValidState(t *testing.T) {
	t.Run("Should return true when state is valid", func(t *testing.T) {
		// Arrange
		state := Created

		// Act
		res := IsValidState(state)

		// Assert
		assert.True(t, res)
	})

	t.Run("Should return false when state is invalid", func(t *testing.T) {
		// Arrange
		state := State(0)

		// Act
		res := IsValidState(state)

		// Assert
		assert.False(t, res)
	})
}
