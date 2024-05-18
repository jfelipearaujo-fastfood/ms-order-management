package get

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_validator"
)

type GetOrderDto struct {
	OrderId    string `param:"id" validate:"required_without=TrackId,uuid-when-not-empty"`
	TrackId    string `param:"track_id" validate:"required_without=OrderId,track-id-when-not-empty"`
	CustomerId string `param:"customer_id" validate:"required_without=OrderId TrackId,uuid-when-not-empty"`
}

func (dto *GetOrderDto) FindViaID() bool {
	return dto.OrderId != ""
}

func (dto *GetOrderDto) FindViaCustomerID() bool {
	return dto.CustomerId != ""
}

func (dto *GetOrderDto) Validate() error {
	validator := validator.New()
	err := custom_validator.RegisterCustomValidations(validator)
	if err != nil {
		return err
	}

	if err := validator.Struct(dto); err != nil {
		return custom_error.ErrRequestNotValid
	}

	return nil
}
