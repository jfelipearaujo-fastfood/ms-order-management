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
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider/time_provider"
	order_repository "github.com/jfelipearaujo-org/ms-order-management/internal/repository/order"
	order_create_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/create"
	order_get_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	order_update_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/update"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config   *environment.Config
	Database database.DatabaseService
	Topic    cloud.TopicService
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

	return &Server{
		Config:   config,
		Database: database.NewDatabase(config),
		Topic:    cloud.NewService(config.CloudConfig.OrderPaymentTopicName, cloudConfig),
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
	e.Use(middleware.Recover())

	s.registerHealthCheck(e)

	group := e.Group(fmt.Sprintf("/api/%s", s.Config.ApiConfig.ApiVersion))

	s.registerOrderHandlers(group)

	return e
}

func (server *Server) registerHealthCheck(e *echo.Echo) {
	healthHandler := health.NewHandler(server.Database)

	e.GET("/health", healthHandler.Handle)
}

func (s *Server) registerOrderHandlers(e *echo.Group) {
	// providers
	timeProvider := time_provider.NewTimeProvider(time.Now)

	// repositories
	repository := order_repository.NewOrderRepository(s.Database.GetInstance())

	// services
	createOrderService := order_create_service.NewService(repository, timeProvider)
	getOrderService := order_get_service.NewService(repository)
	updateOrderService := order_update_service.NewService(repository, timeProvider)

	// handlers
	createOrderHandler := create.NewHandler(createOrderService)
	addOrderItemHandler := add_item.NewHandler(getOrderService, updateOrderService)
	getOrderByIdOrTrackIdHandler := get_by_id.NewHandler(getOrderService)

	e.POST("/orders", createOrderHandler.Handle)
	e.POST("/orders/:id/items", addOrderItemHandler.Handle)
	e.GET("/orders/:id", getOrderByIdOrTrackIdHandler.Handle)
	e.GET("/orders/tracking/:track_id", getOrderByIdOrTrackIdHandler.Handle)
}
