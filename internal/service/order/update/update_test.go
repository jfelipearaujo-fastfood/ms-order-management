package update

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	provider_mock "github.com/jfelipearaujo-org/ms-order-management/internal/provider/mocks"
	repository_mock "github.com/jfelipearaujo-org/ms-order-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should update order without items", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		repository.On("Update", ctx, mock.Anything, false).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		order := &order_entity.Order{}

		req := UpdateOrderDto{
			OrderId: uuid.NewString(),
			State:   1,
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.NoError(t, err)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should update order with items", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		repository.On("Update", ctx, mock.Anything, true).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Times(2)

		service := NewService(repository, timeProvider)

		order := &order_entity.Order{}

		req := UpdateOrderDto{
			OrderId: uuid.NewString(),
			State:   1,
			Items: []UpdateOrderItemDto{
				{
					ItemId:    uuid.NewString(),
					Name:      "name",
					UnitPrice: 10.0,
					Quantity:  1,
				},
			},
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.NoError(t, err)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		service := NewService(repository, timeProvider)

		order := &order_entity.Order{}

		req := UpdateOrderDto{}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.Error(t, err)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when try to update the order", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		repository.On("Update", ctx, mock.Anything, false).
			Return(assert.AnError).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		order := &order_entity.Order{}

		req := UpdateOrderDto{
			OrderId: uuid.NewString(),
			State:   1,
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.Error(t, err)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when cannot update the state", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		order := &order_entity.Order{
			State: order_entity.Received,
		}

		req := UpdateOrderDto{
			OrderId: uuid.NewString(),
			State:   1,
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.Error(t, err)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when cannot add an item", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		timeProvider.On("GetTime").
			Return(now).
			Times(2)

		service := NewService(repository, timeProvider)

		itemId := uuid.NewString()

		order := &order_entity.Order{
			State: order_entity.Created,
			Items: []order_entity.Item{
				{
					Id: itemId,
				},
			},
		}

		req := UpdateOrderDto{
			OrderId: uuid.NewString(),
			State:   1,
			Items: []UpdateOrderItemDto{
				{
					ItemId:    itemId,
					Name:      "name",
					UnitPrice: 10.0,
					Quantity:  1,
				},
			},
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.Error(t, err)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})
}
