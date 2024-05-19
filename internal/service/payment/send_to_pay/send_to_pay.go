package send_to_pay

import (
	"context"
	"log/slog"

	"github.com/jfelipearaujo-org/ms-order-management/internal/adapter/cloud"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
)

type Service struct {
	topic        cloud.TopicService
	repository   repository.PaymentRepository
	timeProvider provider.TimeProvider
}

func NewService(
	topic cloud.TopicService,
	repository repository.PaymentRepository,
	timeProvider provider.TimeProvider,
) *Service {
	return &Service{
		topic:        topic,
		repository:   repository,
		timeProvider: timeProvider,
	}
}

func (s *Service) Handle(ctx context.Context, order *order_entity.Order, request SendToPayDto) error {
	if err := request.Validate(); err != nil {
		return err
	}

	order.CalculateTotals()

	request.TotalItems = order.TotalItems
	request.Amount = order.TotalPrice

	messageId, err := s.topic.PublishMessage(ctx, request)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "message sent to topic", "topic", s.topic.GetTopicName(), "message_id", *messageId)

	payment := payment_entity.NewPayment(
		order.Id,
		request.PaymentId,
		order.TotalItems,
		order.TotalPrice,
		s.timeProvider.GetTime(),
	)

	if !request.Resend {
		if err := s.repository.Create(ctx, &payment); err != nil {
			return err
		}
	}
	return nil
}
