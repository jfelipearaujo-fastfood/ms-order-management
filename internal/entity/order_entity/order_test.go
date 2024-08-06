package order_entity

import (
	"testing"
	"time"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	t.Run("Should create a new order", func(t *testing.T) {
		// Arrange
		now := time.Now()

		// Act
		res := NewOrder("customer_id", now)

		// Assert
		assert.NotEmpty(t, res.Id)
		assert.Equal(t, "customer_id", res.CustomerId)
		assert.NotEmpty(t, res.TrackId)
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
			Id:        "item_id",
			Name:      "name",
			UnitPrice: 1.23,
			Quantity:  1,
		}

		order := NewOrder("customer_id", now)

		// Act
		err := order.AddItem(NewItem("item_id", "name", 1.23, 1), now)

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

		err = order.AddItem(NewItem("item_id_1", "name", 1.23, 1), now)
		assert.NoError(t, err)

		err = order.AddItem(NewItem("item_id_2", "name", 2.34, 2), now)
		assert.NoError(t, err)

		// Act
		order.CalculateTotals()

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
		states := []OrderState{Delivered, Cancelled}

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
		states := []OrderState{
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

	t.Run("Should return an error when trying to add an item that already exists", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		order.Items = append(order.Items, NewItem("item_id", "name", 1.23, 1))

		// Act
		err := order.AddItem(NewItem("item_id", "name", 1.23, 1), now)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrOrderItemAlreadyExists)
	})

	t.Run("Should return true if the order has items", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		order.Items = append(order.Items, NewItem("item_id", "name", 1.23, 1))

		// Act
		res := order.HasItems()

		// Assert
		assert.True(t, res)
	})

	t.Run("Should return true if the order has on going payments", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		payment := payment_entity.NewPayment(order.Id, "payment_id", 1, 1.23, now)

		order.Payments = append(order.Payments, payment)

		// Act
		res := order.HasOnGoingPayments()

		// Assert
		assert.True(t, res)
	})

	t.Run("Should return false if the order has no on going payments", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		payment := payment_entity.NewPayment(order.Id, "payment_id", 1, 1.23, now)
		payment.State = payment_entity.Rejected

		order.Payments = append(order.Payments, payment)

		// Act
		res := order.HasOnGoingPayments()

		// Assert
		assert.False(t, res)
	})

	t.Run("Should return the payment by id", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		payment := payment_entity.NewPayment(order.Id, "payment_id", 1, 1.23, now)

		order.Payments = append(order.Payments, payment)

		// Act
		res := order.GetPaymentByID("payment_id")

		// Assert
		assert.NotNil(t, res)
		assert.Equal(t, payment.PaymentId, res.PaymentId)
	})

	t.Run("Should return nil when the payment id does not exist", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)

		// Act
		res := order.GetPaymentByID("invalid_payment_id")

		// Assert
		assert.Nil(t, res)
	})

	t.Run("Should return the on going payment", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		payment := payment_entity.NewPayment(order.Id, "payment_id", 1, 1.23, now)
		payment.State = payment_entity.WaitingForApproval

		order.Payments = append(order.Payments, payment)

		// Act
		res := order.GetOnGoingPayment()

		// Assert
		assert.NotNil(t, res)
		assert.Equal(t, payment.PaymentId, res.PaymentId)
	})

	t.Run("Should return nil when there is no on going payment", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)

		// Act
		res := order.GetOnGoingPayment()

		// Assert
		assert.Nil(t, res)
	})
}

func TestOrder_ShouldCancel(t *testing.T) {
	t.Run("should cancel when all payments are rejected", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)

		order.Payments = []payment_entity.Payment{
			{PaymentId: "1", State: payment_entity.Rejected},
			{PaymentId: "2", State: payment_entity.Rejected},
			{PaymentId: "3", State: payment_entity.Rejected},
		}

		// Act
		res := order.ShouldCancel()

		// Assert
		assert.True(t, res)
	})

	t.Run("should not cancel when there is at least one payment that is not rejected", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)
		payment := payment_entity.NewPayment(order.Id, "payment_id", 1, 1.23, now)
		payment.State = payment_entity.Approved

		order.Payments = append(order.Payments, payment)

		// Act
		res := order.ShouldCancel()

		// Assert
		assert.False(t, res)
	})

	t.Run("should not cancel when there are no payments", func(t *testing.T) {
		// Arrange
		now := time.Now()

		order := NewOrder("customer_id", now)

		// Act
		res := order.ShouldCancel()

		// Assert
		assert.False(t, res)
	})
}
