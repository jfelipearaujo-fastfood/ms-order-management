package payment

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/payment_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/payment/send_to_pay"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler(t *testing.T) {
	t.Run("Should create a payment", func(t *testing.T) {
		// Arrange
		sendToPayService := mocks.NewMockSendToPayService[send_to_pay.SendToPayDto](t)
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		getOrderService.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{
				Items: []order_entity.Item{
					{
						Id: uuid.NewString(),
					},
				},
			}, nil).
			Once()

		sendToPayService.On("Handle", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		reqBody := send_to_pay.SendToPayDto{
			OrderID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(sendToPayService, getOrderService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
		sendToPayService.AssertExpectations(t)
		getOrderService.AssertExpectations(t)
	})

	t.Run("Should return error when order is not found", func(t *testing.T) {
		// Arrange
		sendToPayService := mocks.NewMockSendToPayService[send_to_pay.SendToPayDto](t)
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		getOrderService.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{}, custom_error.ErrOrderNotFound).
			Once()

		reqBody := send_to_pay.SendToPayDto{
			OrderID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(sendToPayService, getOrderService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusNotFound, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusNotFound,
			Message: "unable to find the order",
			Details: "order not found",
		}, he.Message)

		sendToPayService.AssertExpectations(t)
		getOrderService.AssertExpectations(t)
	})

	t.Run("Should return internal error when trying to find the order", func(t *testing.T) {
		// Arrange
		sendToPayService := mocks.NewMockSendToPayService[send_to_pay.SendToPayDto](t)
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		getOrderService.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{}, assert.AnError).
			Once()

		reqBody := send_to_pay.SendToPayDto{
			OrderID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(sendToPayService, getOrderService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusInternalServerError, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Details: "assert.AnError general error for testing",
		}, he.Message)

		sendToPayService.AssertExpectations(t)
		getOrderService.AssertExpectations(t)
	})

	t.Run("Should return error when order has no items", func(t *testing.T) {
		// Arrange
		sendToPayService := mocks.NewMockSendToPayService[send_to_pay.SendToPayDto](t)
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		getOrderService.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{}, nil).
			Once()

		reqBody := send_to_pay.SendToPayDto{
			OrderID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(sendToPayService, getOrderService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusBadRequest,
			Message: "operation not allowed",
			Details: "order has no items",
		}, he.Message)

		sendToPayService.AssertExpectations(t)
		getOrderService.AssertExpectations(t)
	})

	t.Run("Should return error when order has on going payments", func(t *testing.T) {
		// Arrange
		sendToPayService := mocks.NewMockSendToPayService[send_to_pay.SendToPayDto](t)
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		getOrderService.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{
				Items: []order_entity.Item{
					{
						Id: uuid.NewString(),
					},
				},
				Payments: []payment_entity.Payment{
					{
						State: payment_entity.WaitingForApproval,
					},
				},
			}, nil).
			Once()

		reqBody := send_to_pay.SendToPayDto{
			OrderID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(sendToPayService, getOrderService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusBadRequest,
			Message: "operation not allowed",
			Details: "order has on going payments or is already paid",
		}, he.Message)

		sendToPayService.AssertExpectations(t)
		getOrderService.AssertExpectations(t)
	})

	t.Run("Should return error when payment service return domain error", func(t *testing.T) {
		// Arrange
		sendToPayService := mocks.NewMockSendToPayService[send_to_pay.SendToPayDto](t)
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		getOrderService.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{
				Items: []order_entity.Item{
					{
						Id: uuid.NewString(),
					},
				},
			}, nil).
			Once()

		sendToPayService.On("Handle", mock.Anything, mock.Anything, mock.Anything).
			Return(custom_error.ErrRequestNotValid).
			Once()

		reqBody := send_to_pay.SendToPayDto{
			OrderID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(sendToPayService, getOrderService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusUnprocessableEntity, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusUnprocessableEntity,
			Message: "validation error",
			Details: "request not valid, please check the fields",
		}, he.Message)

		sendToPayService.AssertExpectations(t)
		getOrderService.AssertExpectations(t)
	})

	t.Run("Should return internal error when trying to create the payment", func(t *testing.T) {
		// Arrange
		sendToPayService := mocks.NewMockSendToPayService[send_to_pay.SendToPayDto](t)
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		getOrderService.On("Handle", mock.Anything, mock.Anything).
			Return(order_entity.Order{
				Items: []order_entity.Item{
					{
						Id: uuid.NewString(),
					},
				},
			}, nil).
			Once()

		sendToPayService.On("Handle", mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		reqBody := send_to_pay.SendToPayDto{
			OrderID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(sendToPayService, getOrderService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusInternalServerError, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
			Details: "assert.AnError general error for testing",
		}, he.Message)

		sendToPayService.AssertExpectations(t)
		getOrderService.AssertExpectations(t)
	})

}
