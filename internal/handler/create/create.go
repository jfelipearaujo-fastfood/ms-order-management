package create

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-order-management/internal/service"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/create"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/errors"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service service.CreateOrderService[create.CreateOrderDto]
}

func NewHandler(
	service service.CreateOrderService[create.CreateOrderDto],
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request create.CreateOrderDto

	if err := ctx.Bind(&request); err != nil {
		return errors.NewHttpAppError(http.StatusBadRequest, "invalid request", err)
	}

	context := ctx.Request().Context()

	order, err := h.service.Handle(context, request)
	if err != nil {
		if err == errors.ErrRequestNotValid {
			return errors.NewHttpAppError(http.StatusUnprocessableEntity, "validation error", err)
		}
		if err == errors.ErrOrderAlreadyExists {
			return errors.NewHttpAppError(http.StatusConflict, "order cannot be created", err)
		}

		return errors.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	return ctx.JSON(http.StatusCreated, order)
}
