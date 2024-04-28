package create

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
)

type CreateOrderDto struct {
	CustomerID string `json:"customer_id" validate:"required,uuid4"`
}

func (dto *CreateOrderDto) Validate() error {
	validator := validator.New()

	if err := validator.Struct(dto); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
