package create

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	provider_mock "github.com/jfelipearaujo-org/ms-order-management/internal/provider/mocks"
	repository_mock "github.com/jfelipearaujo-org/ms-order-management/internal/repository/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return nil when request is valid", func(t *testing.T) {
		// Arrange
		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		now := time.Now()

		repository.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
			Return(0, []order_entity.Order{}, nil).
			Once()

		repository.On("Create", mock.Anything, mock.Anything).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		ctx := context.Background()

		request := CreateOrderDto{
			CustomerID: uuid.NewString(),
		}

		// Act
		resp, err := service.Handle(ctx, request)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when request is invalid", func(t *testing.T) {
		// Arrange
		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		service := NewService(repository, timeProvider)

		ctx := context.Background()

		request := CreateOrderDto{
			CustomerID: "abc",
		}

		// Act
		resp, err := service.Handle(ctx, request)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when repository returns error", func(t *testing.T) {
		// Arrange
		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		now := time.Now()

		repository.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
			Return(0, []order_entity.Order{}, nil).
			Once()

		repository.On("Create", mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(repository, timeProvider)

		ctx := context.Background()

		request := CreateOrderDto{
			CustomerID: uuid.NewString(),
		}

		// Act
		resp, err := service.Handle(ctx, request)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when already exists an active order for the customer", func(t *testing.T) {
		// Arrange
		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		repository.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
			Return(1, []order_entity.Order{}, nil).
			Once()

		service := NewService(repository, timeProvider)

		ctx := context.Background()

		request := CreateOrderDto{
			CustomerID: uuid.NewString(),
		}

		// Act
		resp, err := service.Handle(ctx, request)

		// Assert
		assert.Error(t, err)
		assert.ErrorIs(t, err, custom_error.ErrOrderAlreadyExists)
		assert.Nil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when repository returns error when search for on-going order", func(t *testing.T) {
		// Arrange
		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		repository.On("GetAll", mock.Anything, mock.Anything, mock.Anything).
			Return(0, []order_entity.Order{}, assert.AnError).
			Once()

		service := NewService(repository, timeProvider)

		ctx := context.Background()

		request := CreateOrderDto{
			CustomerID: uuid.NewString(),
		}

		// Act
		resp, err := service.Handle(ctx, request)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, resp)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})
}
