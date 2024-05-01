package get_all

import (
	"github.com/jfelipearaujo-org/ms-order-management/internal/common"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
)

type GetOrdersDto struct {
	CustomerID string `query:"customer_id"`
	State      int    `query:"state"`

	common.Pagination
}

func (dto *GetOrdersDto) Validate() error {
	if order_entity.IsValidState(order_entity.OrderState(dto.State)) {
		return nil
	}

	return custom_error.ErrRequestNotValid
}
