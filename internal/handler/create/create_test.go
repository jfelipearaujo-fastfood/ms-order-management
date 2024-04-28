package create

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/create"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandle(t *testing.T) {
	t.Run("Should create an order", func(t *testing.T) {
		//  Arrange
		service := mocks.NewMockCreateOrderService[create.CreateOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(&entity.Order{}, nil).
			Once()

		reqBody := create.CreateOrderDto{
			CustomerID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Code)
		service.AssertExpectations(t)
	})

	t.Run("Should return bad request when request is invalid", func(t *testing.T) {
		//  Arrange
		service := mocks.NewMockCreateOrderService[create.CreateOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(nil, custom_error.ErrRequestNotValid).
			Once()

		reqBody := create.CreateOrderDto{
			CustomerID: "invalid",
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(service)

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

		service.AssertExpectations(t)
	})

	t.Run("Should return conflict when order already exists", func(t *testing.T) {
		//  Arrange
		service := mocks.NewMockCreateOrderService[create.CreateOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(nil, custom_error.ErrOrderAlreadyExists).
			Once()

		reqBody := create.CreateOrderDto{
			CustomerID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(service)

		// Act
		err = handler.Handle(ctx)

		// Assert
		assert.Error(t, err)

		he, ok := err.(*echo.HTTPError)
		assert.True(t, ok)

		assert.Equal(t, http.StatusConflict, he.Code)
		assert.Equal(t, custom_error.AppError{
			Code:    http.StatusConflict,
			Message: "unable to create the order",
			Details: "order already exists",
		}, he.Message)

		service.AssertExpectations(t)
	})

	t.Run("Should return internal server error when an unexpected error occurs", func(t *testing.T) {
		//  Arrange
		service := mocks.NewMockCreateOrderService[create.CreateOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(nil, assert.AnError).
			Once()

		reqBody := create.CreateOrderDto{
			CustomerID: uuid.NewString(),
		}

		body, err := json.Marshal(reqBody)
		assert.NoError(t, err)

		req := httptest.NewRequest(echo.POST, "/", bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)

		handler := NewHandler(service)

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

		service.AssertExpectations(t)
	})
}
