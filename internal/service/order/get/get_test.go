package get

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHandle(t *testing.T) {
	t.Run("Should return order when find via id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		id := uuid.NewString()

		repository := mocks.NewMockOrderRepository(t)

		repository.On("GetByID", ctx, id).
			Return(entity.Order{
				UUID: id,
			}, nil).
			Once()

		service := NewService(repository)

		req := GetOrderDto{
			UUID: id,
		}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
		repository.AssertExpectations(t)
	})

	t.Run("Should return order when find via track id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		trackId := "ABC123"

		repository := mocks.NewMockOrderRepository(t)

		repository.On("GetByTrackID", ctx, trackId).
			Return(entity.Order{
				TrackID: entity.NewTrackIDFrom(trackId),
			}, nil).
			Once()

		service := NewService(repository)

		req := GetOrderDto{
			TrackId: trackId,
		}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.NotEmpty(t, resp)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderRepository(t)

		service := NewService(repository)

		req := GetOrderDto{}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resp)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when validation returns error when find by id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderRepository(t)

		service := NewService(repository)

		req := GetOrderDto{}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resp)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when validation returns error when find by track id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		repository := mocks.NewMockOrderRepository(t)

		service := NewService(repository)

		req := GetOrderDto{}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resp)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when repository returns error when find by id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		id := uuid.NewString()

		repository := mocks.NewMockOrderRepository(t)

		repository.On("GetByID", ctx, id).
			Return(entity.Order{}, assert.AnError).
			Once()

		service := NewService(repository)

		req := GetOrderDto{
			UUID: id,
		}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resp)
		repository.AssertExpectations(t)
	})

	t.Run("Should return error when repository returns error when find by track id", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		trackId := "ABC123"

		repository := mocks.NewMockOrderRepository(t)

		repository.On("GetByTrackID", ctx, trackId).
			Return(entity.Order{}, assert.AnError).
			Once()

		service := NewService(repository)

		req := GetOrderDto{
			TrackId: trackId,
		}

		// Act
		resp, err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, resp)
		repository.AssertExpectations(t)
	})
}
