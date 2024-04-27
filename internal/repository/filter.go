package repository

import "github.com/jfelipearaujo-org/ms-order-management/internal/entity"

type GetAllOrdersFilter struct {
	State entity.State
}
