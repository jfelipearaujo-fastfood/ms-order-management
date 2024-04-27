package service

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/common"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
)

type CreateOrderService[T any] interface {
	Handle(ctx context.Context, request T) (*entity.Order, error)
}

type GetOrderDto struct {
	UUID    string `param:"id"`       // api/v1/orders/:id
	TrackID string `param:"track_id"` // api/v1/orders/tracking/:track_id
}

type GetOrderService interface {
	Handle(ctx context.Context, request GetOrderDto) (entity.Order, error)
}

type GetOrdersDto struct {
	State int `query:"state"`

	common.Pagination
}

type GetOrdersService interface {
	Handle(ctx context.Context, request GetOrdersDto) ([]entity.Order, error)
}

type UpdateOrderDto struct {
	UUID string `param:"id"`

	State int           `json:"state"`
	Items []entity.Item `json:"items"`
}

type UpdateOrderService interface {
	Handle(ctx context.Context, request UpdateOrderDto) (*entity.Order, error)
}
