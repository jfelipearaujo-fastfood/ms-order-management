package payment

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("Should create a payment", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectExec("INSERT INTO (.+)?order_payments(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Create(ctx, &payment_entity.Payment{})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error when the insert fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectExec("INSERT INTO (.+)?order_payments(.+)?").
			WillReturnError(assert.AnError)

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Create(ctx, &payment_entity.Payment{})

		// Assert
		assert.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Should update a payment", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectExec("UPDATE (.+)?order_payments(.+)?").
			WillReturnResult(sqlmock.NewResult(1, 1))

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Update(ctx, &payment_entity.Payment{})

		// Assert
		assert.NoError(t, err)
	})

	t.Run("Should return an error when the update fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectExec("UPDATE (.+)?order_payments(.+)?").
			WillReturnError(assert.AnError)

		repo := NewPaymentRepository(db)

		// Act
		err = repo.Update(ctx, &payment_entity.Payment{})

		// Assert
		assert.Error(t, err)
	})
}
