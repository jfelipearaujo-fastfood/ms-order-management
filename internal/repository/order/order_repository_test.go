package order_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	t.Run("Should create an order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := entity.NewOrder("customer_id", now)
		item := entity.NewItem("product_id", 1, 10.0)
		order.AddItem(item, now)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.UUID, order.CustomerID, order.TrackID, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.UUID, item.UUID, item.Quantity, item.UnitPrice).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error if something got wrong while try to insert the items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := entity.NewOrder("customer_id", now)
		item := entity.NewItem("product_id", 1, 10.0)
		order.AddItem(item, now)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.UUID, order.CustomerID, order.TrackID, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.UUID, item.UUID, item.Quantity, item.UnitPrice).
			WillReturnError(errors.New("something got wrong"))

		mock.ExpectRollback()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error if something got wrong while try to insert the order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := entity.NewOrder("customer_id", now)
		item := entity.NewItem("product_id", 1, 10.0)
		order.AddItem(item, now)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.UUID, order.CustomerID, order.TrackID, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
			WillReturnError(errors.New("something got wrong"))

		mock.ExpectRollback()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error while try to rollback the order insertion", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := entity.NewOrder("customer_id", now)
		item := entity.NewItem("product_id", 1, 10.0)
		order.AddItem(item, now)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.UUID, order.CustomerID, order.TrackID, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
			WillReturnError(errors.New("something got wrong"))

		mock.ExpectRollback().
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error while try to rollback the order item insertion", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := entity.NewOrder("customer_id", now)
		item := entity.NewItem("product_id", 1, 10.0)
		order.AddItem(item, now)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.UUID, order.CustomerID, order.TrackID, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.UUID, item.UUID, item.Quantity, item.UnitPrice).
			WillReturnError(errors.New("something got wrong"))

		mock.ExpectRollback().
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error if something got wrong while try to create a transaction", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := entity.NewOrder("customer_id", now)
		item := entity.NewItem("product_id", 1, 10.0)
		order.AddItem(item, now)

		mock.ExpectBegin().
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		err = repo.Create(ctx, &order)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
