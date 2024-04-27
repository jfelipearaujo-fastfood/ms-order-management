package create

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	provider_mock "github.com/jfelipearaujo-org/ms-order-management/internal/provider/mocks"
	repository_mock "github.com/jfelipearaujo-org/ms-order-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return nil when request is valid", func(t *testing.T) {
		// Arrange
		repository := repository_mock.NewMockOrderRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		now := time.Now()

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

		repository.On("Create", mock.Anything, mock.Anything).
			Return(errors.New("something got wrong")).
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
}
