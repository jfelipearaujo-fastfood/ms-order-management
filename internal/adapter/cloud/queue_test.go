package cloud

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/testtools"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetQueueName(t *testing.T) {
	t.Run("Should return queue name", func(t *testing.T) {
		// Arrange
		fakeProcessor := mocks.NewMockProcessMessageService[process.ProcessMessageDto](t)

		service := NewQueueService("test-queue", aws.Config{}, fakeProcessor)

		// Act
		queueName := service.GetQueueName()

		// Assert
		assert.Equal(t, "test-queue", queueName)
	})
}

func TestUpdateQueueUrl(t *testing.T) {
	t.Run("Should return nil when queue is found", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		fakeProcessor := mocks.NewMockProcessMessageService[process.ProcessMessageDto](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		// Act
		err := service.UpdateQueueUrl(ctx)

		// Assert
		assert.NoError(t, err)
		testtools.ExitTest(stubber, t)
	})

	t.Run("Should return error when GetQueueUrl operation fails", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		raiseErr := &testtools.StubError{Err: errors.New("ClientError")}

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Error:         raiseErr,
		})

		fakeProcessor := mocks.NewMockProcessMessageService[process.ProcessMessageDto](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		// Act
		err := service.UpdateQueueUrl(ctx)

		// Assert
		testtools.VerifyError(err, raiseErr, t)
		testtools.ExitTest(stubber, t)
	})
}

func TestStartConsuming(t *testing.T) {
	t.Run("Should start consuming messages", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"be6293ff-4ec0-4ed8-95c9-b36ce99aa105\",\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
					},
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629b"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567891"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567891"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := mocks.NewMockProcessMessageService[process.ProcessMessageDto](t)

		fakeProcessor.On("Handle", ctx, mock.Anything).
			Return(nil).
			Times(2)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})

	t.Run("Should log error when cant unmarshal notification", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : false,
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"be6293ff-4ec0-4ed8-95c9-b36ce99aa105\",\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : false,
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := mocks.NewMockProcessMessageService[process.ProcessMessageDto](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})

	t.Run("Should log error when notification type is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "InvalidType",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"be6293ff-4ec0-4ed8-95c9-b36ce99aa105\",\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := mocks.NewMockProcessMessageService[process.ProcessMessageDto](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})

	t.Run("Should log error when message body is invalid", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "this is invalid",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := mocks.NewMockProcessMessageService[process.ProcessMessageDto](t)

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})

	t.Run("Should log error when something goes wrong when processing the message", func(t *testing.T) {
		// Arrange
		ctx := context.Background()
		stubber := testtools.NewStubber()

		stubber.Add(testtools.Stub{
			OperationName: "GetQueueUrl",
			Input: &sqs.GetQueueUrlInput{
				QueueName: aws.String("test-queue"),
			},
			Output: &sqs.GetQueueUrlOutput{
				QueueUrl: aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
			},
		})

		response := `{
			"Type" : "Notification",
			"MessageId" : "fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a",
			"TopicArn" : "arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic",
			"Message" : "{\"order_id\":\"be6293ff-4ec0-4ed8-95c9-b36ce99aa105\",\"payment_id\":\"a5c81ac9-a549-44c5-bb09-c330116b929f\",\"items\":[{\"id\":\"3822eb8e-3da9-416e-a248-3551fc628566\",\"name\":\"Hamburguer\",\"quantity\":1},{\"id\":\"ca685ace-ef25-4aa3-97f5-489394aa6356\",\"name\":\"Refrigerante\",\"quantity\":1}],\"total_items\":2,\"amount\":59.980000000000004}",
			"Timestamp" : "2024-05-19T02:01:36.927Z",
			"SignatureVersion" : "1",
			"Signature" : "e2Jex1vYJslu5gc0YPvaoprA6Vnbus7VuaQOjKVoegQ8i+5yqtWD47Zl7+O5mh/vLOEcNKkXKVNDk++idzRxEg40uZQcWOwDewqaItZvD2XH6b/mqYAnf4QjAjIF3+orXpSZQn/hatp7KzsYvd7bnPmO3YyzuqwD4t4Zz19GvatIuYsjDkcueWXX5/HOJJhAGSQFg/hnETAnllWZuDAgwDOUF6sPfa7zSUGSyj2ymHlSyMPNOLmM5VMpouujU0lFwYlZqHwg3WbEONRHyZ7Fs6JO8wPRG1J3kUvjcZ7qQwo4ARGTIbXZ7xJv9mYjE79Sdl3S5yXkvg4CambuE9Gpig==",
			"SigningCertURL" : "https://sns.us-east-1.amazonaws.com/SimpleNotificationService-60eadc530605d63b8e62a523676ef735.pem",
			"UnsubscribeURL" : "https://sns.us-east-1.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-east-1:000000000000:OrderPaymentTopic:961e369d-aee9-40d8-ab2e-4c6a5e2eab95"
		}`

		stubber.Add(testtools.Stub{
			OperationName: "ReceiveMessage",
			Input: &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     20,
			},
			Output: &sqs.ReceiveMessageOutput{
				Messages: []types.Message{
					{
						MessageId:     aws.String("fc8e9ffd-6122-5c52-8fb9-c13e3ee2629a"),
						Body:          aws.String(response),
						ReceiptHandle: aws.String("1234567890"),
					},
				},
			},
		})

		stubber.Add(testtools.Stub{
			OperationName: "DeleteMessage",
			Input: &sqs.DeleteMessageInput{
				QueueUrl:      aws.String("https://sqs.us-east-1.amazonaws.com/123456789012/test-queue"),
				ReceiptHandle: aws.String("1234567890"),
			},
			Output: &sqs.DeleteMessageOutput{},
		})

		fakeProcessor := mocks.NewMockProcessMessageService[process.ProcessMessageDto](t)

		fakeProcessor.On("Handle", ctx, mock.Anything).
			Return(assert.AnError).
			Once()

		service := NewQueueService("test-queue", *stubber.SdkConfig, fakeProcessor)

		err := service.UpdateQueueUrl(ctx)
		assert.NoError(t, err)

		// Act
		service.ConsumeMessages(ctx)

		// Assert
		testtools.ExitTest(stubber, t)
		fakeProcessor.AssertExpectations(t)
	})
}
