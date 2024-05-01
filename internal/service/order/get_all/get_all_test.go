package get_all

import (
	"context"
	"testing"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return the orders", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderRepository(t)

		repository.On("GetAll", ctx, mock.Anything, mock.Anything).
			Return(1, []order_entity.Order{{}}, nil).
			Once()

		service := NewService(repository)

		req := GetOrdersDto{
			State: 1,
		}

		// Act
		count, res, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 1, count)
		assert.Len(t, res, 1)
		repository.AssertExpectations(t)
	})

	t.Run("Should return an error when request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderRepository(t)

		service := NewService(repository)

		req := GetOrdersDto{}

		// Act
		count, res, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Zero(t, count)
		assert.Nil(t, res)
		repository.AssertExpectations(t)
	})

	t.Run("Should return an error when repository returns an error", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderRepository(t)

		repository.On("GetAll", ctx, mock.Anything, mock.Anything).
			Return(0, nil, assert.AnError).
			Once()

		service := NewService(repository)

		req := GetOrdersDto{
			State: 1,
		}

		// Act
		count, res, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Zero(t, count)
		assert.Nil(t, res)
		repository.AssertExpectations(t)
	})
}
