package add_item

import (
	"net/http"

	"github.com/jfelipearaujo-org/ms-order-management/internal/service"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/update"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	getService    service.GetOrderService[get.GetOrderDto]
	updateService service.UpdateOrderService[update.UpdateOrderDto]
}

func NewHandler(
	getService service.GetOrderService[get.GetOrderDto],
	updateService service.UpdateOrderService[update.UpdateOrderDto],
) *Handler {
	return &Handler{
		getService:    getService,
		updateService: updateService,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request update.UpdateOrderDto

	if err := ctx.Bind(&request); err != nil {
		return custom_error.NewHttpAppError(http.StatusBadRequest, "invalid request", err)
	}

	context := ctx.Request().Context()

	getOrderRequest := get.GetOrderDto{
		OrderId: request.OrderId,
	}

	order, err := h.getService.Handle(context, getOrderRequest)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	if order.IsCompleted() {
		return custom_error.NewHttpAppErrorFromBusinessError(custom_error.ErrOrderAlreadyCompleted)
	}

	if len(request.Items) > 0 && !order.CanAddItems() {
		return custom_error.NewHttpAppErrorFromBusinessError(custom_error.ErrOrderInProgress)
	}

	if order.HasOnGoingPayments() {
		return custom_error.NewHttpAppErrorFromBusinessError(custom_error.ErrOrderHasOnGoingPayments)
	}

	request.State = int(order.State)

	if err := h.updateService.Handle(context, &order, request); err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	order.RefreshStateTitle()

	return ctx.JSON(http.StatusOK, order)
}
