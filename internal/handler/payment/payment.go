package payment

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/payment/send_to_pay"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	sendToPayService service.SendToPayService[send_to_pay.SendToPayDto]
	getOrderService  service.GetOrderService[get.GetOrderDto]
}

func NewHandler(
	sendToPayService service.SendToPayService[send_to_pay.SendToPayDto],
	getOrderService service.GetOrderService[get.GetOrderDto],
) *Handler {
	return &Handler{
		sendToPayService: sendToPayService,
		getOrderService:  getOrderService,
	}
}

func (h *Handler) Handle(ctx echo.Context) error {
	var request send_to_pay.SendToPayDto

	if err := ctx.Bind(&request); err != nil {
		return custom_error.NewHttpAppError(http.StatusBadRequest, "invalid request", err)
	}

	if err := (&echo.DefaultBinder{}).BindQueryParams(ctx, &request); err != nil {
		return custom_error.NewHttpAppError(http.StatusBadRequest, "invalid request", err)
	}

	context := ctx.Request().Context()

	getOrderRequest := get.GetOrderDto{
		OrderId: request.OrderID,
	}

	order, err := h.getOrderService.Handle(context, getOrderRequest)
	if err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	if !order.HasItems() {
		return custom_error.NewHttpAppErrorFromBusinessError(custom_error.ErrOrderHasNoItems)
	}

	if !request.Resend && order.HasOnGoingPayments() {
		return custom_error.NewHttpAppErrorFromBusinessError(custom_error.ErrOrderHasOnGoingPayments)
	}

	request.PaymentId = uuid.NewString()

	request.Items = []send_to_pay.SendToPayItemDto{}

	for _, item := range order.Items {
		request.Items = append(request.Items, send_to_pay.SendToPayItemDto{
			Id:       item.Id,
			Name:     item.Name,
			Quantity: item.Quantity,
		})
	}

	request.TotalItems = order.TotalItems
	request.Amount = order.TotalPrice

	if err := h.sendToPayService.Handle(context, &order, request); err != nil {
		if custom_error.IsBusinessErr(err) {
			return custom_error.NewHttpAppErrorFromBusinessError(err)
		}

		return custom_error.NewHttpAppError(http.StatusInternalServerError, "internal server error", err)
	}

	ok := map[string]string{
		"message": "payment sent to be paid",
	}

	return ctx.JSON(http.StatusCreated, ok)
}
