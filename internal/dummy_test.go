package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDummy(t *testing.T) {
	t.Run("Dummy", func(t *testing.T) {
		err := Dummy()

		assert.Nil(t, err)
	})
}
