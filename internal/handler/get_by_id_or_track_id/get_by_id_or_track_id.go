package get_by_id

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-order-management/internal/service"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service service.GetOrderService[get.GetOrderDto]
}

func NewHandler(service service.GetOrderService[get.GetOrderDto]) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request get.GetOrderDto

	if err := ctx.Bind(&request); err != nil {
		return custom_error.NewHttpAppError(http.StatusBadRequest, "invalid request", err)
	}

	context := ctx.Request().Context()

	order, err := h.service.Handle(context, request)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	order.RefreshStateTitle()
	order.CalculateTotals()

	return ctx.JSON(http.StatusOK, order)
}
