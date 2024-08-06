package process

import (
	"context"
	"testing"
	"time"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
	provider_mocks "github.com/jfelipearaujo-org/ms-order-management/internal/provider/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should return an error when the message is not valid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.Error(t, err)
	})
}

func TestHandleOrderResponse(t *testing.T) {
	t.Run("Should process a message with order response", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		orderRepository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{
				State: order_entity.Created,
			}, nil).
			Once()

		orderRepository.On("Update", ctx, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{
			OrderId: "order-id",
			OrderResponse: &OrderResponse{
				State: "Received",
			},
		}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.NoError(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error when the order is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		orderRepository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{}, custom_error.ErrOrderNotFound).
			Once()

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{
			OrderId: "order-id",
			OrderResponse: &OrderResponse{
				State: "Received",
			},
		}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.Error(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error when the order state transition is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		orderRepository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{
				State: order_entity.Received,
			}, nil).
			Once()

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{
			OrderId: "order-id",
			OrderResponse: &OrderResponse{
				State: "Received",
			},
		}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.Error(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error when the order update fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		orderRepository.On("GetByID", ctx, mock.Anything).
			Return(order_entity.Order{
				State: order_entity.Created,
			}, nil).
			Once()

		orderRepository.On("Update", ctx, mock.Anything, mock.Anything).
			Return(custom_error.ErrOrderInProgress).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{
			OrderId: "order-id",
			OrderResponse: &OrderResponse{
				State: "Received",
			},
		}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.Error(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})
}

func TestHandlePaymentResponse(t *testing.T) {
	t.Run("Should process a message with payment response", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		order := order_entity.Order{
			Payments: []payment_entity.Payment{
				{
					PaymentId: "payment-id",
					State:     payment_entity.WaitingForApproval,
				},
			},
		}

		orderRepository.On("GetByID", ctx, mock.Anything).
			Return(order, nil).
			Times(2)

		paymentRepository.On("Update", ctx, mock.Anything).
			Return(nil).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{
			OrderId: "order-id",
			PaymentResponse: &PaymentResponse{
				PaymentId: "payment-id",
				State:     "Approved",
			},
		}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.NoError(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error when the payment is not found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		order := order_entity.Order{
			Payments: []payment_entity.Payment{
				{
					PaymentId: "payment-id",
					State:     payment_entity.WaitingForApproval,
				},
			},
		}

		orderRepository.On("GetByID", ctx, mock.Anything).
			Return(order, nil).
			Once()

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{
			OrderId: "order-id",
			PaymentResponse: &PaymentResponse{
				PaymentId: "unknown-payment-id",
				State:     "Approved",
			},
		}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.Error(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error when the payment state transition is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		order := order_entity.Order{
			Payments: []payment_entity.Payment{
				{
					PaymentId: "payment-id",
					State:     payment_entity.Approved,
				},
			},
		}

		orderRepository.On("GetByID", ctx, mock.Anything).
			Return(order, nil).
			Once()

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{
			OrderId: "order-id",
			PaymentResponse: &PaymentResponse{
				PaymentId: "payment-id",
				State:     "Approved",
			},
		}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.Error(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should return an error when the payment update fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		now := time.Now()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		order := order_entity.Order{
			Payments: []payment_entity.Payment{
				{
					PaymentId: "payment-id",
					State:     payment_entity.WaitingForApproval,
				},
			},
		}

		orderRepository.On("GetByID", ctx, mock.Anything).
			Return(order, nil).
			Once()

		paymentRepository.On("Update", ctx, mock.Anything).
			Return(assert.AnError).
			Once()

		timeProvider.On("GetTime").
			Return(now).
			Once()

		service := NewService(orderRepository, paymentRepository, timeProvider)

		message := ProcessMessageDto{
			OrderId: "order-id",
			PaymentResponse: &PaymentResponse{
				PaymentId: "payment-id",
				State:     "Approved",
			},
		}

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.Error(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})

	t.Run("Should update the order state to cancelled when all payments are rejected", func(t *testing.T) {
		// Arrange
		ctx := context.Background()

		orderRepository := mocks.NewMockOrderRepository(t)
		paymentRepository := mocks.NewMockPaymentRepository(t)
		timeProvider := provider_mocks.NewMockTimeProvider(t)

		message := ProcessMessageDto{
			OrderId: "order_id",
			PaymentResponse: &PaymentResponse{
				PaymentId: "3",
				State:     "Rejected",
			},
		}

		order := order_entity.NewOrder("customer_id", time.Now())
		order.Payments = []payment_entity.Payment{
			{PaymentId: "1", State: payment_entity.Rejected},
			{PaymentId: "2", State: payment_entity.Rejected},
			{PaymentId: "3", State: payment_entity.WaitingForApproval},
		}

		orderRepository.On("GetByID", ctx, "order_id").Return(order, nil)

		paymentRepository.On("Update", ctx, mock.Anything).Return(nil)
		orderRepository.On("Update", ctx, mock.Anything, mock.Anything).Return(nil)

		timeProvider.On("GetTime").Return(time.Now())

		service := NewService(orderRepository, paymentRepository, timeProvider)

		// Act
		err := service.Handle(ctx, message)

		// Assert
		assert.NoError(t, err)
		orderRepository.AssertExpectations(t)
		paymentRepository.AssertExpectations(t)
		timeProvider.AssertExpectations(t)
	})
}
