package process

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
)

type Service struct {
	orderRepository   repository.OrderRepository
	paymentRepository repository.PaymentRepository
	timeProvider      provider.TimeProvider
}

func NewService(
	orderRepository repository.OrderRepository,
	paymentRepository repository.PaymentRepository,
	timeProvider provider.TimeProvider,
) *Service {
	return &Service{
		orderRepository:   orderRepository,
		paymentRepository: paymentRepository,
		timeProvider:      timeProvider,
	}
}

func (s *Service) Handle(ctx context.Context, message ProcessMessageDto) error {
	if message.OrderResponse == nil &&
		message.PaymentResponse == nil {
		return custom_error.ErrQueueMessageNotValid
	}

	order, err := s.orderRepository.GetByID(ctx, message.OrderId)
	if err != nil {
		return err
	}

	if message.OrderResponse != nil {
		newState := order_entity.NewOrderState(message.OrderResponse.State)

		if !order.State.CanTransitionTo(newState) {
			return custom_error.ErrOrderInvalidStateTransition
		}

		if err := order.UpdateState(newState, s.timeProvider.GetTime()); err != nil {
			return err
		}

		if err := s.orderRepository.Update(ctx, &order, false); err != nil {
			return err
		}
	}

	if message.PaymentResponse != nil {
		payment := order.GetPaymentByID(message.PaymentResponse.PaymentId)

		if payment == nil {
			return custom_error.ErrPaymentNotFound
		}

		newState := payment_entity.NewPaymentState(message.PaymentResponse.State)

		if !payment.State.CanTransitionTo(newState) {
			return custom_error.ErrPaymentInvalidStateTransition
		}

		payment.UpdateState(newState, s.timeProvider.GetTime())

		if err := s.paymentRepository.Update(ctx, payment); err != nil {
			return err
		}

		order, err = s.orderRepository.GetByID(ctx, message.OrderId)
		if err != nil {
			return err
		}

		if order.ShouldCancel() {
			if err := order.UpdateState(order_entity.Cancelled, s.timeProvider.GetTime()); err != nil {
				return err
			}

			if err := s.orderRepository.Update(ctx, &order, false); err != nil {
				return err
			}
		}
	}

	return nil
}
