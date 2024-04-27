package create

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/errors"
)

type CreateOrderDto struct {
	CustomerID string `json:"customer_id" validate:"required,uuid4"`
}

func (dto *CreateOrderDto) Validate() error {
	validator := validator.New()

	if err := validator.Struct(dto); err != nil {
		return errors.ErrRequestNotValid
	}

	return nil
}
