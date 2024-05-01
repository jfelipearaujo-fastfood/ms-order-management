package repository

import "github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"

type GetAllOrdersFilter struct {
	CustomerID string

	StateFrom order_entity.OrderState
	StateTo   order_entity.OrderState
}
