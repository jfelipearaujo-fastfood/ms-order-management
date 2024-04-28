package add_item

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/update"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should add an item to the order", func(t *testing.T) {
		// Arrange
		getService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		getService.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{
				State: entity.Created,
			}, nil).
			Once()

		updateService.On("Handle", mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Once()

		reqBody := update.UpdateOrderDto{
			Items: []update.UpdateOrderItemDto{
				{
					ItemId:    uuid.NewString(),
					UnitPrice: 1.23,
					Quantity:  1,
				},
			},
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(getService, updateService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		getService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return error when order is not found", func(t *testing.T) {
		// Arrange
		getService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		getService.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{}, custom_error.ErrOrderNotFound).
			Once()

		reqBody := update.UpdateOrderDto{
			Items: []update.UpdateOrderItemDto{
				{
					ItemId:    uuid.NewString(),
					UnitPrice: 1.23,
					Quantity:  1,
				},
			},
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(getService, updateService)

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

		getService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return internal error when try to find the order", func(t *testing.T) {
		// Arrange
		getService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		getService.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{}, assert.AnError).
			Once()

		reqBody := update.UpdateOrderDto{
			Items: []update.UpdateOrderItemDto{
				{
					ItemId:    uuid.NewString(),
					UnitPrice: 1.23,
					Quantity:  1,
				},
			},
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(getService, updateService)

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

		getService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return error when order is already completed", func(t *testing.T) {
		// Arrange
		getService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		getService.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{
				State: entity.Delivered,
			}, nil).
			Once()

		reqBody := update.UpdateOrderDto{
			Items: []update.UpdateOrderItemDto{
				{
					ItemId:    uuid.NewString(),
					UnitPrice: 1.23,
					Quantity:  1,
				},
			},
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(getService, updateService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusBadRequest,
			Message: "unable to update/insert information to the order",
			Details: "order is already completed or cancelled",
		}, he.Message)

		getService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return error when order is in progress", func(t *testing.T) {
		// Arrange
		getService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		getService.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{
				State: entity.Received,
			}, nil).
			Once()

		reqBody := update.UpdateOrderDto{
			Items: []update.UpdateOrderItemDto{
				{
					ItemId:    uuid.NewString(),
					UnitPrice: 1.23,
					Quantity:  1,
				},
			},
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(getService, updateService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusBadRequest, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusBadRequest,
			Message: "unable to update/insert information to the order",
			Details: "order is in progress",
		}, he.Message)

		getService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return business error when try to update the order", func(t *testing.T) {
		// Arrange
		getService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		getService.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{
				State: entity.Created,
			}, nil).
			Once()

		updateService.On("Handle", mock.Anything, mock.Anything, mock.Anything).
			Return(custom_error.ErrOrderItemAlreadyExists).
			Once()

		reqBody := update.UpdateOrderDto{
			Items: []update.UpdateOrderItemDto{
				{
					ItemId:    uuid.NewString(),
					UnitPrice: 1.23,
					Quantity:  1,
				},
			},
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(getService, updateService)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusConflict, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusConflict,
			Message: "unable to add an item",
			Details: "order item already exists",
		}, he.Message)

		getService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})

	t.Run("Should return internal error when try to update the order", func(t *testing.T) {
		// Arrange
		getService := mocks.NewMockGetOrderService[get.GetOrderDto](t)
		updateService := mocks.NewMockUpdateOrderService[update.UpdateOrderDto](t)

		getService.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{
				State: entity.Created,
			}, nil).
			Once()

		updateService.On("Handle", mock.Anything, mock.Anything, mock.Anything).
			Return(assert.AnError).
			Once()

		reqBody := update.UpdateOrderDto{
			Items: []update.UpdateOrderItemDto{
				{
					ItemId:    uuid.NewString(),
					UnitPrice: 1.23,
					Quantity:  1,
				},
			},
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/orders/:id/items")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(getService, updateService)

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

		getService.AssertExpectations(t)
		updateService.AssertExpectations(t)
	})
}
