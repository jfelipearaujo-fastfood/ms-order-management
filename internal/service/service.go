package service

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
)

type CreateOrderService[T any] interface {
	Handle(ctx context.Context, request T) (*entity.Order, error)
}

type GetOrderService[T any] interface {
	Handle(ctx context.Context, request T) (entity.Order, error)
}

type GetOrdersService[T any] interface {
	Handle(ctx context.Context, request T) (int, []entity.Order, error)
}

type UpdateOrderService[T any] interface {
	Handle(ctx context.Context, request T) (*entity.Order, error)
}
