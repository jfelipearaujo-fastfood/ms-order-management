package order_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/common"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
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

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.Id, order.CustomerId, order.TrackId, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.Id, item.Id, item.Quantity, item.UnitPrice).
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

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.Id, order.CustomerId, order.TrackId, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.Id, item.Id, item.Quantity, item.UnitPrice).
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

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.Id, order.CustomerId, order.TrackId, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
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

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.Id, order.CustomerId, order.TrackId, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
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

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO orders").
			WithArgs(order.Id, order.CustomerId, order.TrackId, order.State, order.StateUpdatedAt, order.CreatedAt, order.UpdatedAt).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.Id, item.Id, item.Quantity, item.UnitPrice).
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

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

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

func TestGetByID(t *testing.T) {
	t.Run("Should get an order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()
		paymentId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, "ABC123", order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		paymentRows := sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
			AddRow(orderId, paymentId, 1, 10.5, payment_entity.WaitingForApproval, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_payments(.+)?").
			WillReturnRows(paymentRows)

		productId := uuid.NewString()

		orderItemRows := sqlmock.NewRows([]string{"product_id", "quantity", "price"}).
			AddRow(productId, 1, 10.0)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_items(.+)?").
			WillReturnRows(orderItemRows)

		expected := order_entity.Order{
			Id:             orderId,
			CustomerId:     customerId,
			TrackId:        "ABC123",
			State:          order_entity.Created,
			StateUpdatedAt: now,
			Items: []order_entity.Item{
				{
					Id:        productId,
					Quantity:  1,
					UnitPrice: 10.0,
				},
			},
			Payments: []payment_entity.Payment{
				{
					OrderId:    orderId,
					PaymentId:  paymentId,
					TotalItems: 1,
					Amount:     10.5,
					State:      payment_entity.WaitingForApproval,
					StateTitle: "WaitingForApproval",
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByID(ctx, orderId)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEmpty(t, res)
		assert.Equal(t, expected, res)
	})

	t.Run("Should not return error when no order items were found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()
		paymentId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, "ABC123", order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		paymentRows := sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
			AddRow(orderId, paymentId, 1, 10.5, payment_entity.WaitingForApproval, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_payments(.+)?").
			WillReturnRows(paymentRows)

		orderItemRows := sqlmock.NewRows([]string{"product_id", "quantity", "price"})

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_items(.+)?").
			WillReturnRows(orderItemRows)

		expected := order_entity.Order{
			Id:             orderId,
			CustomerId:     customerId,
			TrackId:        "ABC123",
			State:          order_entity.Created,
			StateUpdatedAt: now,
			Items:          []order_entity.Item{},
			Payments: []payment_entity.Payment{
				{
					OrderId:    orderId,
					PaymentId:  paymentId,
					TotalItems: 1,
					Amount:     10.5,
					State:      payment_entity.WaitingForApproval,
					StateTitle: "WaitingForApproval",
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByID(ctx, orderId)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEmpty(t, res)
		assert.Equal(t, expected, res)
	})

	t.Run("Should return error when order is not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		orderId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{})

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByID(ctx, orderId)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrOrderNotFound)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})

	t.Run("Should return error when try to query the order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		orderId := uuid.NewString()

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByID(ctx, orderId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})

	t.Run("Should return a scan error while try to parse the order rows", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, "ABC123", "Created", now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByID(ctx, orderId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})

	t.Run("Should return error when try to query the order items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()
		paymentId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, "ABC123", order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		paymentRows := sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
			AddRow(orderId, paymentId, 1, 10.5, payment_entity.WaitingForApproval, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_payments(.+)?").
			WillReturnRows(paymentRows)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_items(.+)?").
			WillReturnError(assert.AnError)

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByID(ctx, orderId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})

	t.Run("Should return a scan error while try to parse the order items rows", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()
		paymentId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, "ABC123", order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		paymentRows := sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
			AddRow(orderId, paymentId, 1, 10.5, payment_entity.WaitingForApproval, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_payments(.+)?").
			WillReturnRows(paymentRows)

		productId := uuid.NewString()

		orderItemRows := sqlmock.NewRows([]string{"product_id", "quantity", "price"}).
			AddRow(productId, "abc", 10.0)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_items(.+)?").
			WillReturnRows(orderItemRows)

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByID(ctx, orderId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})
}

func TestGetByTrackID(t *testing.T) {
	t.Run("Should get an order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		trackId := "ABC-123"
		customerId := uuid.NewString()
		paymentId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, trackId, order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		paymentRows := sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
			AddRow(orderId, paymentId, 1, 10.5, payment_entity.WaitingForApproval, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_payments(.+)?").
			WillReturnRows(paymentRows)

		productId := uuid.NewString()

		orderItemRows := sqlmock.NewRows([]string{"product_id", "quantity", "price"}).
			AddRow(productId, 1, 10.0)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_items(.+)?").
			WillReturnRows(orderItemRows)

		expected := order_entity.Order{
			Id:             orderId,
			CustomerId:     customerId,
			TrackId:        order_entity.NewTrackIdFrom(trackId),
			State:          order_entity.Created,
			StateUpdatedAt: now,
			Items: []order_entity.Item{
				{
					Id:        productId,
					Quantity:  1,
					UnitPrice: 10.0,
				},
			},
			Payments: []payment_entity.Payment{
				{
					OrderId:    orderId,
					PaymentId:  paymentId,
					TotalItems: 1,
					Amount:     10.5,
					State:      payment_entity.WaitingForApproval,
					StateTitle: "WaitingForApproval",
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByTrackID(ctx, trackId)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEmpty(t, res)
		assert.Equal(t, expected, res)
	})

	t.Run("Should not return error when no order items were found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		trackId := "ABC-123"
		customerId := uuid.NewString()
		paymentId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, trackId, order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		paymentRows := sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
			AddRow(orderId, paymentId, 1, 10.5, payment_entity.WaitingForApproval, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_payments(.+)?").
			WillReturnRows(paymentRows)

		orderItemRows := sqlmock.NewRows([]string{"product_id", "quantity", "price"})

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_items(.+)?").
			WillReturnRows(orderItemRows)

		expected := order_entity.Order{
			Id:             orderId,
			CustomerId:     customerId,
			TrackId:        order_entity.NewTrackIdFrom(trackId),
			State:          order_entity.Created,
			StateUpdatedAt: now,
			Items:          []order_entity.Item{},
			Payments: []payment_entity.Payment{
				{
					OrderId:    orderId,
					PaymentId:  paymentId,
					TotalItems: 1,
					Amount:     10.5,
					State:      payment_entity.WaitingForApproval,
					StateTitle: "WaitingForApproval",
					CreatedAt:  now,
					UpdatedAt:  now,
				},
			},
			CreatedAt: now,
			UpdatedAt: now,
		}

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByTrackID(ctx, trackId)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEmpty(t, res)
		assert.Equal(t, expected, res)
	})

	t.Run("Should return error when order is not found", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		trackId := "ABC-123"

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"})

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByTrackID(ctx, trackId)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrOrderNotFound)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})

	t.Run("Should return error when try to query the order", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		trackId := "ABC123"

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByTrackID(ctx, trackId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})

	t.Run("Should return a scan error while try to parse the order rows", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()
		trackId := "ABC123"

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, trackId, "Created", now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByTrackID(ctx, trackId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})

	t.Run("Should return error when try to query the order items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		trackId := "ABC-123"
		customerId := uuid.NewString()
		paymentId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, trackId, order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		paymentRows := sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
			AddRow(orderId, paymentId, 1, 10.5, payment_entity.WaitingForApproval, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_payments(.+)?").
			WillReturnRows(paymentRows)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_items(.+)?").
			WillReturnError(assert.AnError)

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByTrackID(ctx, trackId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})

	t.Run("Should return a scan error while try to parse the order items rows", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		trackId := "ABC-123"
		customerId := uuid.NewString()
		paymentId := uuid.NewString()

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, trackId, order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)?").
			WillReturnRows(orderRows)

		paymentRows := sqlmock.NewRows([]string{"order_id", "payment_id", "total_items", "amount", "state", "created_at", "updated_at"}).
			AddRow(orderId, paymentId, 1, 10.5, payment_entity.WaitingForApproval, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_payments(.+)?").
			WillReturnRows(paymentRows)

		productId := uuid.NewString()

		orderItemRows := sqlmock.NewRows([]string{"product_id", "quantity", "price"}).
			AddRow(productId, "abc", 10.0)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?order_items(.+)?").
			WillReturnRows(orderItemRows)

		repo := NewOrderRepository(db)

		// Act
		res, err := repo.GetByTrackID(ctx, trackId)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Should get the orders without filtering", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()

		mock.ExpectQuery("SELECT COUNT(.+) FROM (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, "ABC123", order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)? ORDER BY (.+)").
			WillReturnRows(orderRows)

		repo := NewOrderRepository(db)

		pagination := common.Pagination{
			Page: 1,
			Size: 10,
		}

		filter := repository.GetAllOrdersFilter{}

		// Act
		count, res, err := repo.GetAll(ctx, pagination, filter)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, count)
	})

	t.Run("Should return error when something got wrong with order count", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT COUNT(.+) FROM (.+)?orders(.+)?").
			WillReturnError(assert.AnError)

		repo := NewOrderRepository(db)

		pagination := common.Pagination{
			Page: 1,
			Size: 10,
		}

		filter := repository.GetAllOrdersFilter{}

		// Act
		count, res, err := repo.GetAll(ctx, pagination, filter)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
		assert.Equal(t, 0, count)
	})

	t.Run("Should return scan error from order count", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT COUNT(.+) FROM (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow("abc"))

		repo := NewOrderRepository(db)

		pagination := common.Pagination{
			Page: 1,
			Size: 10,
		}

		filter := repository.GetAllOrdersFilter{}

		// Act
		count, res, err := repo.GetAll(ctx, pagination, filter)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
		assert.Equal(t, 0, count)
	})

	t.Run("Should get the orders with filtering", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()

		mock.ExpectQuery("SELECT COUNT(.+) FROM (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, "ABC123", order_entity.Created, now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)? ORDER BY (.+)").
			WillReturnRows(orderRows)

		repo := NewOrderRepository(db)

		pagination := common.Pagination{
			Page: 1,
			Size: 10,
		}

		filter := repository.GetAllOrdersFilter{
			StateFrom: order_entity.Created,
		}

		// Act
		count, res, err := repo.GetAll(ctx, pagination, filter)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.NotEmpty(t, res)
		assert.Equal(t, 1, count)
	})

	t.Run("Should return error when something got wrong with order query", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		mock.ExpectQuery("SELECT COUNT(.+) FROM (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)? ORDER BY (.+)").
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		pagination := common.Pagination{
			Page: 1,
			Size: -10,
		}

		filter := repository.GetAllOrdersFilter{}

		// Act
		count, res, err := repo.GetAll(ctx, pagination, filter)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
		assert.Equal(t, 0, count)
	})

	t.Run("Should return scan error", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		orderId := uuid.NewString()
		customerId := uuid.NewString()

		mock.ExpectQuery("SELECT COUNT(.+) FROM (.+)?orders(.+)?").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		orderRows := sqlmock.NewRows([]string{"id", "customer_id", "track_id", "state", "state_updated_at", "created_at", "updated_at"}).
			AddRow(orderId, customerId, "ABC123", "Created", now, now, now)

		mock.ExpectQuery("SELECT (.+) FROM (.+)?orders(.+)? ORDER BY (.+)").
			WillReturnRows(orderRows)

		repo := NewOrderRepository(db)

		pagination := common.Pagination{
			Page: 1,
			Size: 10,
		}

		filter := repository.GetAllOrdersFilter{}

		// Act
		count, res, err := repo.GetAll(ctx, pagination, filter)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.Empty(t, res)
		assert.Equal(t, 0, count)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Should update an order without update the items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, false)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error when no order were updated", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectRollback()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
		assert.ErrorIs(t, err, custom_error.ErrOrderNotFound)
	})

	t.Run("Should return error if rows affected return error when no order were updated", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("something got wrong")))
		mock.ExpectRollback()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error if rollback fails when rows affected return error and no order were updated", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("something got wrong")))

		mock.ExpectRollback().
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error if rollback fails when no order were updated", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		mock.ExpectRollback().
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error when update order fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnError(errors.New("something got wrong"))
		mock.ExpectRollback()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error if rollback fails when update order fails", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnError(errors.New("something got wrong"))
		mock.ExpectRollback().
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, false)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should update an order and update the items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("DELETE FROM order_items").
			WithArgs(order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.Id, item.Id, item.Quantity, item.UnitPrice).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, true)

		// Assert
		assert.NoError(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error if rollback fails when insert the new items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("DELETE FROM order_items").
			WithArgs(order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.Id, item.Id, item.Quantity, item.UnitPrice).
			WillReturnError(errors.New("something got wrong"))

		mock.ExpectRollback().
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error when insert the new items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("DELETE FROM order_items").
			WithArgs(order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("INSERT INTO order_items").
			WithArgs(order.Id, item.Id, item.Quantity, item.UnitPrice).
			WillReturnError(errors.New("something got wrong"))

		mock.ExpectRollback()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error if rollback fails when deleting old items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("DELETE FROM order_items").
			WithArgs(order.Id).
			WillReturnError(errors.New("something got wrong"))

		mock.ExpectRollback().
			WillReturnError(errors.New("something got wrong"))

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})

	t.Run("Should return error when deleting old items", func(t *testing.T) {
		// Arrange
		db, mock, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ctx := context.Background()

		now := time.Now()

		order := order_entity.NewOrder("customer_id", now)
		item := order_entity.NewItem("product_id", 1, 10.0)

		err = order.AddItem(item, now)
		assert.NoError(t, err)

		mock.ExpectBegin()
		mock.ExpectExec("UPDATE orders").
			WithArgs(order.State, order.StateUpdatedAt, order.UpdatedAt, order.Id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectExec("DELETE FROM order_items").
			WithArgs(order.Id).
			WillReturnError(errors.New("something got wrong"))

		mock.ExpectRollback()

		repo := NewOrderRepository(db)

		// Act
		err = repo.Update(ctx, &order, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, mock.ExpectationsWereMet())
	})
}
