package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"

	"github.com/jfelipearaujo-org/ms-order-management/internal/adapter/cloud"
	"github.com/jfelipearaujo-org/ms-order-management/internal/adapter/database"
	"github.com/jfelipearaujo-org/ms-order-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-order-management/internal/handler/add_item"
	"github.com/jfelipearaujo-org/ms-order-management/internal/handler/create"
	get_by_id "github.com/jfelipearaujo-org/ms-order-management/internal/handler/get_by_id_or_track_id"
	"github.com/jfelipearaujo-org/ms-order-management/internal/handler/health"
	"github.com/jfelipearaujo-org/ms-order-management/internal/handler/payment"
	"github.com/jfelipearaujo-org/ms-order-management/internal/handler/update"
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider/time_provider"
	order_repository "github.com/jfelipearaujo-org/ms-order-management/internal/repository/order"
	payment_repository "github.com/jfelipearaujo-org/ms-order-management/internal/repository/payment"
	order_create_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/create"
	order_get_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/process"
	order_update_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/update"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/payment/send_to_pay"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config          *environment.Config
	DatabaseService database.DatabaseService
	TopicService    cloud.TopicService
	QueueService    cloud.QueueService

	Dependency Dependency
}

func NewServer(config *environment.Config) *Server {
	ctx := context.Background()

	cloudConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	if config.CloudConfig.IsBaseEndpointSet() {
		cloudConfig.BaseEndpoint = aws.String(config.CloudConfig.BaseEndpoint)
	}

	databaseService := database.NewDatabase(config)

	timeProvider := time_provider.NewTimeProvider(time.Now)
	orderRepository := order_repository.NewOrderRepository(databaseService.GetInstance())
	paymentRepository := payment_repository.NewPaymentRepository(databaseService.GetInstance())

	topicService := cloud.NewTopicService(config.CloudConfig.OrderPaymentTopicName, cloudConfig)

	messageProcessor := process.NewService(orderRepository, paymentRepository, timeProvider)

	return &Server{
		Config:          config,
		DatabaseService: databaseService,
		TopicService:    topicService,
		QueueService:    cloud.NewQueueService(config.CloudConfig.UpdateOrderQueueName, cloudConfig, messageProcessor),

		Dependency: Dependency{
			TimeProvider: timeProvider,

			OrderRepository:   orderRepository,
			PaymentRepository: paymentRepository,

			CreateOrderService: order_create_service.NewService(orderRepository, timeProvider),
			GetOrderService:    order_get_service.NewService(orderRepository),
			UpdateOrderService: order_update_service.NewService(orderRepository, timeProvider),
			SendToPayService:   send_to_pay.NewService(topicService, paymentRepository, timeProvider),

			ProcessMessageService: messageProcessor,
		},
	}
}

func (s *Server) GetHttpServer() *http.Server {
	return &http.Server{
		Addr:         fmt.Sprintf(":%d", s.Config.ApiConfig.Port),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(logger.Middleware())
	e.Use(middleware.Recover())

	s.registerHealthCheck(e)

	group := e.Group(fmt.Sprintf("/api/%s", s.Config.ApiConfig.ApiVersion))

	s.registerOrderHandlers(group)

	return e
}

func (server *Server) registerHealthCheck(e *echo.Echo) {
	healthHandler := health.NewHandler(server.DatabaseService)

	e.GET("/health", healthHandler.Handle)
}

func (s *Server) registerOrderHandlers(e *echo.Group) {
	createOrderHandler := create.NewHandler(s.Dependency.CreateOrderService)
	addOrderItemHandler := add_item.NewHandler(s.Dependency.GetOrderService, s.Dependency.UpdateOrderService)
	getOrderByIdOrTrackIdHandler := get_by_id.NewHandler(s.Dependency.GetOrderService)
	sendToPaymentHandler := payment.NewHandler(s.Dependency.SendToPayService, s.Dependency.GetOrderService)
	updateOrderHandler := update.NewHandler(s.Dependency.GetOrderService, s.Dependency.UpdateOrderService)

	e.POST("/orders", createOrderHandler.Handle)
	e.POST("/orders/:id/items", addOrderItemHandler.Handle)
	e.GET("/orders/:id", getOrderByIdOrTrackIdHandler.Handle)
	e.GET("/orders/tracking/:track_id", getOrderByIdOrTrackIdHandler.Handle)
	e.GET("/orders/customer/:customer_id", getOrderByIdOrTrackIdHandler.Handle)
	e.POST("/orders/:order_id/payment", sendToPaymentHandler.Handle)
	e.PATCH("/orders/:id", updateOrderHandler.Handle)
}
