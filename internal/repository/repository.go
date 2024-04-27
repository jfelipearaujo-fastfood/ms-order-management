package repository

import (
	"context"
	"errors"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
)

var (
	ErrOrderNotFound error = errors.New("order not found")
)

type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) error
	GetByID(ctx context.Context, id string) (entity.Order, error)
}
