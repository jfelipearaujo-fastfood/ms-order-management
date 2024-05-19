package update

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-order-management/internal/service"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/update"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	getOrderService service.GetOrderService[get.GetOrderDto]
	updateService   service.UpdateOrderService[update.UpdateOrderDto]
}

func NewHandler(
	getOrderService service.GetOrderService[get.GetOrderDto],
	updateService service.UpdateOrderService[update.UpdateOrderDto],
) *Handler {
	return &Handler{
		getOrderService: getOrderService,
		updateService:   updateService,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request update.UpdateOrderDto

	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	context := ctx.Request().Context()

	order, err := h.getOrderService.Handle(context, get.GetOrderDto{
		OrderId: request.OrderId,
	})
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	if err := h.updateService.Handle(context, &order, request); err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	return ctx.JSON(http.StatusCreated, order)
}
