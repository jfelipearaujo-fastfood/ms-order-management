package cloud

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
)

type TopicService interface {
	UpdateTopicArn(ctx context.Context) error
	PublishMessage(ctx context.Context, message interface{}) (*string, error)
}

type Service struct {
	TopicName string
	TopicArn  string
	Client    *sns.Client
}

func NewService(topicName string, config aws.Config) TopicService {
	client := sns.NewFromConfig(config)

	return &Service{
		TopicName: topicName,
		Client:    client,
	}
}

func (s *Service) UpdateTopicArn(ctx context.Context) error {
	output, err := s.Client.ListTopics(ctx, &sns.ListTopicsInput{})
	if err != nil {
		return err
	}

	for _, topic := range output.Topics {
		if strings.Contains(*topic.TopicArn, s.TopicName) {
			s.TopicArn = *topic.TopicArn
			return nil
		}
	}

	return custom_error.ErrTopicNotFound
}

func (s *Service) PublishMessage(ctx context.Context, message interface{}) (*string, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	req := &sns.PublishInput{
		TopicArn: aws.String(s.TopicArn),
		Message:  aws.String(string(body)),
	}

	out, err := s.Client.Publish(ctx, req)
	if err != nil {
		return nil, err
	}

	return out.MessageId, nil
}
