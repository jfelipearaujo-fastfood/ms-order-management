package send_to_pay

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/adapter/cloud/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return nil when message is sent", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		messageId := "message-id"

		topicService := mocks.NewMockTopicService(t)

		topicService.On("GetTopicName").
			Return("topic-name").
			Once()

		topicService.On("PublishMessage", ctx, mock.Anything).
			Return(&messageId, nil).
			Once()

		service := NewService(topicService)

		req := SendToPayDto{
			OrderID: uuid.NewString(),
		}

		// Act
		err := service.Handle(ctx, req)

		// Assert
		assert.NoError(t, err)
		topicService.AssertExpectations(t)
	})

	t.Run("Should return error when request is not valid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		topicService := mocks.NewMockTopicService(t)

		service := NewService(topicService)

		req := SendToPayDto{}

		// Act
		err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		topicService.AssertExpectations(t)
	})

	t.Run("Should return error when message is not sent", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		topicService := mocks.NewMockTopicService(t)

		topicService.On("PublishMessage", ctx, mock.Anything).
			Return(nil, assert.AnError).
			Once()

		service := NewService(topicService)

		req := SendToPayDto{
			OrderID: uuid.NewString(),
		}

		// Act
		err := service.Handle(ctx, req)

		// Assert
		assert.Error(t, err)
		topicService.AssertExpectations(t)
	})
}
