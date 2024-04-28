package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrackID(t *testing.T) {
	t.Run("Should return a new track ID", func(t *testing.T) {
		// Arrange
		expectedLength := 7

		// Act
		res := NewTrackID()

		// Assert
		assert.Len(t, res, expectedLength)
		assert.Regexp(t, "^[A-Z]{3}-[0-9]{3}$", res)
	})

	t.Run("Should return a different track ID each time", func(t *testing.T) {
		// Arrange
		// Act
		res1 := NewTrackID()
		res2 := NewTrackID()

		// Assert
		assert.NotEqual(t, res1, res2)
	})

	t.Run("Should return a track ID from a string", func(t *testing.T) {
		// Arrange
		expected := TrackID("ABC-123")

		// Act
		res := NewTrackIDFrom("ABC-123")

		// Assert
		assert.Equal(t, expected, res)
	})
}
