package repository

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
}
