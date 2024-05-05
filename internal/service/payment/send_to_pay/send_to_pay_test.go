package send_to_pay

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/adapter/cloud/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	provider_mock "github.com/jfelipearaujo-org/ms-order-management/internal/provider/mocks"
	repository_mock "github.com/jfelipearaujo-org/ms-order-management/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return nil when message is sent", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		messageId := "message-id"

		topicService := mocks.NewMockTopicService(t)
		repository := repository_mock.NewMockPaymentRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		topicService.On("GetTopicName").
			Return("topic-name").
			Once()

		topicService.On("PublishMessage", ctx, mock.Anything).
			Return(&messageId, nil).
			Once()

		now := time.Now()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		repository.On("Create", ctx, mock.Anything).
			Return(nil).
			Once()

		service := NewService(topicService, repository, timeProvider)

		order := &order_entity.Order{
			Id: uuid.NewString(),
			Items: []order_entity.Item{
				{
					Id:        uuid.NewString(),
					Name:      "name",
					UnitPrice: 10.5,
					Quantity:  1,
				},
			},
		}

		req := SendToPayDto{
			OrderID:   uuid.NewString(),
			PaymentId: uuid.NewString(),
			Items: []SendToPayItemDto{
				{
					Id:       uuid.NewString(),
					Name:     "name",
					Quantity: 1,
				},
			},
			TotalItems: 1,
			Amount:     10.5,
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.NoError(t, err)
		topicService.AssertExpectations(t)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when request is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		topicService := mocks.NewMockTopicService(t)
		repository := repository_mock.NewMockPaymentRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		service := NewService(topicService, repository, timeProvider)

		order := &order_entity.Order{
			Id: uuid.NewString(),
			Items: []order_entity.Item{
				{
					UnitPrice: 10.5,
					Quantity:  1,
				},
			},
		}

		req := SendToPayDto{
			OrderID:   "",
			PaymentId: uuid.NewString(),
			Items: []SendToPayItemDto{
				{
					Id:       uuid.NewString(),
					Name:     "name",
					Quantity: 1,
				},
			},
			TotalItems: 1,
			Amount:     10.5,
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.Error(t, err)
		topicService.AssertExpectations(t)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when message is not sent", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		topicService := mocks.NewMockTopicService(t)
		repository := repository_mock.NewMockPaymentRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		topicService.On("PublishMessage", ctx, mock.Anything).
			Return(nil, assert.AnError).
			Once()

		service := NewService(topicService, repository, timeProvider)

		order := &order_entity.Order{
			Id: uuid.NewString(),
			Items: []order_entity.Item{
				{
					UnitPrice: 10.5,
					Quantity:  1,
				},
			},
		}

		req := SendToPayDto{
			OrderID:   uuid.NewString(),
			PaymentId: uuid.NewString(),
			Items: []SendToPayItemDto{
				{
					Id:       uuid.NewString(),
					Name:     "name",
					Quantity: 1,
				},
			},
			TotalItems: 1,
			Amount:     10.5,
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.Error(t, err)
		topicService.AssertExpectations(t)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return error when payment is not created", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		messageId := "message-id"

		topicService := mocks.NewMockTopicService(t)
		repository := repository_mock.NewMockPaymentRepository(t)
		timeProvider := provider_mock.NewMockTimeProvider(t)

		topicService.On("GetTopicName").
			Return("topic-name").
			Once()

		topicService.On("PublishMessage", ctx, mock.Anything).
			Return(&messageId, nil).
			Once()

		now := time.Now()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		repository.On("Create", ctx, mock.Anything).
			Return(assert.AnError).
			Once()

		service := NewService(topicService, repository, timeProvider)

		order := &order_entity.Order{
			Id: uuid.NewString(),
			Items: []order_entity.Item{
				{
					UnitPrice: 10.5,
					Quantity:  1,
				},
			},
		}

		req := SendToPayDto{
			OrderID:   uuid.NewString(),
			PaymentId: uuid.NewString(),
			Items: []SendToPayItemDto{
				{
					Id:       uuid.NewString(),
					Name:     "name",
					Quantity: 1,
				},
			},
			TotalItems: 1,
			Amount:     10.5,
		}

		// Act
		err := service.Handle(ctx, order, req)

		// Assert
		assert.Error(t, err)
		topicService.AssertExpectations(t)
		repository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})
}
