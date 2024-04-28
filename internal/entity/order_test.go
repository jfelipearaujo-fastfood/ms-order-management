package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewOrder(t *testing.T) {
	t.Run("Should create a new order", func(t *testing.T) {
		// Arrange
		now := time.Now()

		// Act
		res := NewOrder("customer_id", now)

		// Assert
		assert.NotEmpty(t, res.UUID)
		assert.Equal(t, "customer_id", res.CustomerID)
		assert.NotEmpty(t, res.TrackID)
		assert.Equal(t, Created, res.State)
		assert.Equal(t, now, res.StateUpdatedAt)
		assert.Equal(t, 0, res.TotalItems)
		assert.Equal(t, 0.0, res.TotalPrice)
		assert.Empty(t, res.Items)
		assert.Equal(t, now, res.CreatedAt)
		assert.Equal(t, now, res.UpdatedAt)
	})

	t.Run("Should add an item to the order", func(t *testing.T) {
		// Arrange
		now := time.Now()

		expectedItem := Item{
			UUID:      "item_id",
			UnitPrice: 1.23,
			Quantity:  1,
		}

		order := NewOrder("customer_id", now)

		// Act
		err := order.AddItem(NewItem("item_id", 1.23, 1), now)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, order.Items, 1)
		assert.Contains(t, order.Items, expectedItem)
	})

	t.Run("Should calculate the total items and price", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)

		var err error

		err = order.AddItem(NewItem("item_id_1", 1.23, 1), now)
		assert.NoError(t, err)

		err = order.AddItem(NewItem("item_id_2", 2.34, 2), now)
		assert.NoError(t, err)

		// Act
		order.calculateTotals()

		// Assert
		assert.Equal(t, 3, order.TotalItems)
		assert.Equal(t, 5.91, order.TotalPrice)
	})

	t.Run("Should update the state of the order", func(t *testing.T) {
		// Arrange
		past := time.Now().Add(-time.Hour)
		now := time.Now()

		order := NewOrder("customer_id", past)

		// Act
		err := order.UpdateState(Received, now)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, Received, order.State)
		assert.Equal(t, now, order.StateUpdatedAt)
		assert.Equal(t, now, order.UpdatedAt)
	})

	t.Run("Should return an error when trying to update the state to an invalid state", func(t *testing.T) {
		// Arrange
		past := time.Now().Add(-time.Hour)
		now := time.Now()

		order := NewOrder("customer_id", past)

		// Act
		err := order.UpdateState(Completed, now)

		// Assert
		assert.Error(t, err)
	})

	t.Run("Should not update the state if is the same state", func(t *testing.T) {
		// Arrange
		past := time.Now().Add(-time.Hour)
		now := time.Now()

		order := NewOrder("customer_id", past)

		// Act
		err := order.UpdateState(Created, now)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, Created, order.State)
		assert.Equal(t, past, order.StateUpdatedAt)
		assert.Equal(t, past, order.UpdatedAt)
	})

	t.Run("Should refresh the state title", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)

		// Act
		order.RefreshStateTitle()

		// Assert
		assert.Equal(t, "Created", order.StateTitle)
	})

	t.Run("Should allow add an item to the order", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)

		// Act
		res := order.CanAddItems()

		// Assert
		assert.True(t, res)
	})

	t.Run("Should not allow add an item to the order", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		order.State = Received

		// Act
		res := order.CanAddItems()

		// Assert
		assert.False(t, res)
	})

	t.Run("Should return true if the order is already completed", func(t *testing.T) {
		// Arrange
		states := []State{Delivered, Cancelled}

		for _, state := range states {
			now := time.Now()

			order := NewOrder("customer_id", now)
			order.State = state

			// Act
			res := order.IsCompleted()

			// Assert
			assert.True(t, res)
		}
	})

	t.Run("Should return false if the order is not completed", func(t *testing.T) {
		// Arrange
		states := []State{
			None,
			Created,
			Received,
			Processing,
			Completed,
		}

		for _, state := range states {
			now := time.Now()

			order := NewOrder("customer_id", now)
			order.State = state

			// Act
			res := order.IsCompleted()

			// Assert
			assert.False(t, res)
		}
	})
}
