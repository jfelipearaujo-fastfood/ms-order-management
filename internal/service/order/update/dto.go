package update

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/errors"
)

type UpdateOrderItemDto struct {
	UUID      string  `json:"id" validate:"required,uuid4"`
	UnitPrice float64 `json:"unit_price" validate:"required,min=0.01,max=1000"`
	Quantity  int     `json:"quantity" validate:"required,min=1,max=100"`
}

type UpdateOrderDto struct {
	UUID string `param:"id" validate:"required,uuid4"`

	State int                  `json:"state" validate:"required"`
	Items []UpdateOrderItemDto `json:"items" validate:"dive"`
}

func (dto *UpdateOrderDto) Validate() error {
	validator := validator.New()

	if err := validator.Struct(dto); err != nil {
		return errors.ErrRequestNotValid
	}

	return nil
}
