package get

import (
	"github.com/go-playground/validator/v10"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_validator"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/errors"
)

type GetOrderDto struct {
	UUID    string `param:"id" validate:"required_without=TrackId,uuid-when-not-empty"`        // api/v1/orders/:id
	TrackId string `param:"track_id" validate:"required_without=UUID,track-id-when-not-empty"` // api/v1/orders/tracking/:track_id
}

func (dto *GetOrderDto) FindViaID() bool {
	return dto.UUID != ""
}

func (dto *GetOrderDto) Validate() error {
	validator := validator.New()
	err := custom_validator.RegisterCustomValidations(validator)
	if err != nil {
		return err
	}

	if err := validator.Struct(dto); err != nil {
		return errors.ErrRequestNotValid
	}

	return nil
}
