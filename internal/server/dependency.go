package server

import (
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider/time_provider"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service"
	order_create_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/create"
	order_get_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/process"
	order_update_service "github.com/jfelipearaujo-org/ms-order-management/internal/service/order/update"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/payment/send_to_pay"
)

type Dependency struct {
	TimeProvider *time_provider.TimeProvider

	OrderRepository   repository.OrderRepository
	PaymentRepository repository.PaymentRepository

	CreateOrderService service.CreateOrderService[order_create_service.CreateOrderDto]
	GetOrderService    service.GetOrderService[order_get_service.GetOrderDto]
	UpdateOrderService service.UpdateOrderService[order_update_service.UpdateOrderDto]
	SendToPayService   service.SendToPayService[send_to_pay.SendToPayDto]

	ProcessMessageService service.ProcessMessageService[process.ProcessMessageDto]
}
