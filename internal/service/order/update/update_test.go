package update

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
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

		repository.On("GetByID", ctx, mock.Anything).
			Return(entity.Order{}, nil).
			Once()

		repository.On("Update", ctx, mock.Anything, false).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Times(2)

		service := NewService(repository, timeProvider)

		req := UpdateOrderDto{
			UUID:  uuid.NewString(),
			State: 1,
		}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should update order with items", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(entity.Order{}, nil).
			Once()

		repository.On("Update", ctx, mock.Anything, true).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Times(2)

		service := NewService(repository, timeProvider)

		req := UpdateOrderDto{
			UUID:  uuid.NewString(),
			State: 1,
			Items: []UpdateOrderItemDto{
				{
					UUID:      uuid.NewString(),
					UnitPrice: 10.0,
					Quantity:  1,
				},
			},
		}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		service := NewService(repository, timeProvider)

		req := UpdateOrderDto{}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when try to find the order", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(entity.Order{}, assert.AnError).
			Once()

		service := NewService(repository, timeProvider)

		req := UpdateOrderDto{
			UUID:  uuid.NewString(),
			State: 1,
		}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when try to update the order", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		repository.On("GetByID", ctx, mock.Anything).
			Return(entity.Order{}, nil).
			Once()

		repository.On("Update", ctx, mock.Anything, false).
			Return(assert.AnError).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Times(2)

		service := NewService(repository, timeProvider)

		req := UpdateOrderDto{
			UUID:  uuid.NewString(),
			State: 1,
		}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})
}
