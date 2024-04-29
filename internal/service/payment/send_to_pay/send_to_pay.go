package send_to_pay

import (
	"context"
	"log/slog"

	"github.com/jfelipearaujo-org/ms-order-management/internal/adapter/cloud"
)

type Service struct {
	topic cloud.TopicService
}

func NewService(topic cloud.TopicService) *Service {
	return &Service{
		topic: topic,
	}
}

func (s *Service) Handle(ctx context.Context, request SendToPayDto) error {
	if err := request.Validate(); err != nil {
		return err
	}

	messageId, err := s.topic.PublishMessage(ctx, request)
	if err != nil {
		return err
	}

	slog.InfoContext(ctx, "message sent to topic %s with id %s", s.topic.GetTopicName(), messageId)

	return nil
}
