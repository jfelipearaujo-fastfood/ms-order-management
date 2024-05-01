package server

import (
	"testing"

	"github.com/jfelipearaujo-org/ms-order-management/internal/environment"
	"github.com/stretchr/testify/assert"
)

func TestNewServer(t *testing.T) {
	t.Run("Should return a new server", func(t *testing.T) {
		// Arrange
		config := &environment.Config{
			ApiConfig: &environment.ApiConfig{
				Port: 8080,
			},
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
			CloudConfig: &environment.CloudConfig{
				OrderPaymentTopicName: "order-payment-topic",
				UpdateOrderQueueName:  "update-order-queue",
			},
		}

		// Act
		server := NewServer(config)

		// Assert
		assert.NotNil(t, server)
	})

	t.Run("Should return a new server with base endpoint", func(t *testing.T) {
		// Arrange
		config := &environment.Config{
			ApiConfig: &environment.ApiConfig{
				Port: 8080,
			},
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
			CloudConfig: &environment.CloudConfig{
				OrderPaymentTopicName: "order-payment-topic",
				UpdateOrderQueueName:  "update-order-queue",
				BaseEndpoint:          "http://localhost:8080",
			},
		}

		// Act
		server := NewServer(config)

		// Assert
		assert.NotNil(t, server)
	})

	t.Run("Should create a http server", func(t *testing.T) {
		// Arrange
		config := &environment.Config{
			ApiConfig: &environment.ApiConfig{
				Port: 8080,
			},
			DbConfig: &environment.DatabaseConfig{
				Url: "postgres://host:1234",
			},
			CloudConfig: &environment.CloudConfig{
				OrderPaymentTopicName: "order-payment-topic",
				UpdateOrderQueueName:  "update-order-queue",
			},
		}

		server := NewServer(config)

		// Act
		httpServer := server.GetHttpServer()

		// Assert
		assert.NotNil(t, httpServer)
		assert.Equal(t, ":8080", httpServer.Addr)
	})
}
