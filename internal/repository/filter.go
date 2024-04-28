package repository

import "github.com/jfelipearaujo-org/ms-order-management/internal/entity"

type GetAllOrdersFilter struct {
	CustomerID string

	StateFrom entity.State
	StateTo   entity.State
}
