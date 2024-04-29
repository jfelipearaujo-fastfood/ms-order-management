package repository

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/common"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *order_entity.Order) error
	GetByID(ctx context.Context, id string) (order_entity.Order, error)
	GetByTrackID(ctx context.Context, trackId string) (order_entity.Order, error)
	GetAll(ctx context.Context, pagination common.Pagination, filter GetAllOrdersFilter) (int, []order_entity.Order, error)
	Update(ctx context.Context, order *order_entity.Order, updateItems bool) error
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *payment_entity.Payment) error
}
