package custom_error

import "net/http"

var (
	ErrRequestNotValid BusinessError = New(http.StatusUnprocessableEntity, "validation error", "request not valid, please check the fields")

	ErrOrderInvalidStateTransition BusinessError = New(http.StatusBadRequest, "unable to update order state", "invalid state transition")
	ErrOrderNotFound               BusinessError = New(http.StatusNotFound, "unable to find the order", "order not found")
	ErrOrderAlreadyExists          BusinessError = New(http.StatusConflict, "unable to create the order", "order already exists")
	ErrOrderItemAlreadyExists      BusinessError = New(http.StatusConflict, "unable to add an item", "order item already exists")
	ErrOrderInProgress             BusinessError = New(http.StatusBadRequest, "unable to update/insert information to the order", "order is in progress")
	ErrOrderAlreadyCompleted       BusinessError = New(http.StatusBadRequest, "unable to update/insert information to the order", "order is already completed or cancelled")
)
