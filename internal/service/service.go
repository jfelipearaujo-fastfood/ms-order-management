package service

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
)

type CreateOrderService[T any] interface {
	Handle(ctx context.Context, request T) (*order_entity.Order, error)
}

type GetOrderService[T any] interface {
	Handle(ctx context.Context, request T) (order_entity.Order, error)
}

type GetOrdersService[T any] interface {
	Handle(ctx context.Context, request T) (int, []order_entity.Order, error)
}

type UpdateOrderService[T any] interface {
	Handle(ctx context.Context, order *order_entity.Order, request T) error
}

// ---

type SendToPayService[T any] interface {
	Handle(ctx context.Context, order *order_entity.Order, request T) error
}

// ---

type ProcessMessageService[T any] interface {
	Handle(ctx context.Context, message T) error
}
