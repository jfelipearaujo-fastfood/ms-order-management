package send_to_pay

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
)

type SendToPayDto struct {
	OrderID string `json:"order_id" validate:"required,uuid4"`
}

func (dto *SendToPayDto) Validate() error {
	validator := validator.New()

	if err := validator.Struct(dto); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
