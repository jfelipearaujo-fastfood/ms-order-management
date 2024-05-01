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
	GetTopicName() string
	UpdateTopicArn(ctx context.Context) error
	PublishMessage(ctx context.Context, message interface{}) (*string, error)
}

type AwsSnsService struct {
	TopicName string
	TopicArn  string
	Client    *sns.Client
}

func NewTopicService(topicName string, config aws.Config) TopicService {
	client := sns.NewFromConfig(config)

	return &AwsSnsService{
		TopicName: topicName,
		Client:    client,
	}
}

func (s *AwsSnsService) GetTopicName() string {
	return s.TopicName
}

func (s *AwsSnsService) UpdateTopicArn(ctx context.Context) error {
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

func (s *AwsSnsService) PublishMessage(ctx context.Context, message interface{}) (*string, error) {
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
