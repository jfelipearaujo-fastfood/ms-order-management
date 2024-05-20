package update

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/update"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should update the order", func(t *testing.T) {
		// Arrange
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		order := order_entity.NewOrder(uuid.NewString(), time.Now())

		getOrderService.On("Handle", mock.Anything, get.GetOrderDto{
			OrderId: order.Id,
		}).
			Return(order, nil).
			Once()

		updateService.On("Handle", mock.Anything, mock.Anything, update.UpdateOrderDto{
			OrderId: order.Id,
			State:   2,
		}).
			Return(nil).
			Once()

		reqBody := update.UpdateOrderDto{
			OrderId: order.Id,
			State:   2,
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(getOrderService, updateService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
		getOrderService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return error when the order does not exist", func(t *testing.T) {
		// Arrange
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		order := order_entity.NewOrder(uuid.NewString(), time.Now())

		getOrderService.On("Handle", mock.Anything, get.GetOrderDto{
			OrderId: order.Id,
		}).
			Return(order_entity.Order{}, custom_error.ErrOrderNotFound).
			Once()

		reqBody := update.UpdateOrderDto{
			OrderId: order.Id,
			State:   2,
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(getOrderService, updateService)

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

		getOrderService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return error when try to find the order", func(t *testing.T) {
		// Arrange
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		order := order_entity.NewOrder(uuid.NewString(), time.Now())

		getOrderService.On("Handle", mock.Anything, get.GetOrderDto{
			OrderId: order.Id,
		}).
			Return(order_entity.Order{}, assert.AnError).
			Once()

		reqBody := update.UpdateOrderDto{
			OrderId: order.Id,
			State:   2,
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.PATCH, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(getOrderService, updateService)

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

		getOrderService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return validation error when try to update the order", func(t *testing.T) {
		// Arrange
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		order := order_entity.NewOrder(uuid.NewString(), time.Now())

		getOrderService.On("Handle", mock.Anything, get.GetOrderDto{
			OrderId: order.Id,
		}).
			Return(order, nil).
			Once()

		updateService.On("Handle", mock.Anything, mock.Anything, update.UpdateOrderDto{
			OrderId: order.Id,
			State:   2,
		}).
			Return(custom_error.ErrRequestNotValid).
			Once()

		reqBody := update.UpdateOrderDto{
			OrderId: order.Id,
			State:   2,
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(getOrderService, updateService)

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

		getOrderService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return error when try to update the order", func(t *testing.T) {
		// Arrange
		getOrderService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		order := order_entity.NewOrder(uuid.NewString(), time.Now())

		getOrderService.On("Handle", mock.Anything, get.GetOrderDto{
			OrderId: order.Id,
		}).
			Return(order, nil).
			Once()

		updateService.On("Handle", mock.Anything, mock.Anything, update.UpdateOrderDto{
			OrderId: order.Id,
			State:   2,
		}).
			Return(assert.AnError).
			Once()

		reqBody := update.UpdateOrderDto{
			OrderId: order.Id,
			State:   2,
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(getOrderService, updateService)

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

		getOrderService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})
}
