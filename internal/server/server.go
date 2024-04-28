package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jfelipearaujo-org/ms-order-management/internal/adapter/database"
	"github.com/jfelipearaujo-org/ms-order-management/internal/environment"
	"github.com/jfelipearaujo-org/ms-order-management/internal/handler/create"
	"github.com/jfelipearaujo-org/ms-order-management/internal/handler/health"
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider/time_provider"
	order_repository "github.com/jfelipearaujo-org/ms-order-management/internal/repository/order"
	order_create_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/create"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Config *environment.Config
	db     database.DatabaseService
}

func NewServer(config *environment.Config) *http.Server {
	server := &Server{
		Config: config,
		db:     database.NewDatabase(config),
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%d", server.Config.ApiConfig.Port),
		Handler:      server.RegisterRoutes(),
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
	healthHandler := health.NewHandler(server.db)

	e.GET("/health", healthHandler.Handle)
}

func (s *Server) registerOrderHandlers(e *echo.Group) {
	// providers
	timeProvider := time_provider.NewTimeProvider(time.Now)

	// repositories
	repository := order_repository.NewOrderRepository(s.db.GetInstance())

	// services
	createOrderService := order_create_service.NewService(repository, timeProvider)

	// handlers
	createOrderHandler := create.NewHandler(createOrderService)

	e.POST("/orders", createOrderHandler.Handle)
}
