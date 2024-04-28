package get_by_id

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/mocks"
	"github.com/jfelipearaujo-org/ms-order-management/internal/service/order/get"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleGetById(t *testing.T) {
	t.Run("Should return an order", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{}, nil).
			Once()

		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/api/v1/orders/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(service)

		// Act
		err := handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		service.AssertExpectations(t)
	})

	t.Run("Should return business error", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{}, custom_error.ErrOrderNotFound).
			Once()

		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/api/v1/orders/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(service)

		// Act
		err := handler.Handle(ctx)

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
		service.AssertExpectations(t)
	})

	t.Run("Should return internal server error", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{}, assert.AnError).
			Once()

		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/api/v1/orders/:id")
		ctx.SetParamNames("id")
		ctx.SetParamValues(uuid.NewString())

		handler := NewHandler(service)

		// Act
		err := handler.Handle(ctx)

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

func TestHandleGetByTrackId(t *testing.T) {
	t.Run("Should return an order", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{}, nil).
			Once()

		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/api/v1/orders/tracking/:track_id")
		ctx.SetParamNames("track_id")
		ctx.SetParamValues("ABC-123")

		handler := NewHandler(service)

		// Act
		err := handler.Handle(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)
		service.AssertExpectations(t)
	})

	t.Run("Should return business error", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{}, custom_error.ErrOrderNotFound).
			Once()

		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/api/v1/orders/tracking/:track_id")
		ctx.SetParamNames("track_id")
		ctx.SetParamValues("ABC-123")

		handler := NewHandler(service)

		// Act
		err := handler.Handle(ctx)

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
		service.AssertExpectations(t)
	})

	t.Run("Should return internal server error", func(t *testing.T) {
		// Arrange
		service := mocks.NewMockGetOrderService[get.GetOrderDto](t)

		service.On("Handle", mock.Anything, mock.Anything).
			Return(entity.Order{}, assert.AnError).
			Once()

		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		resp := httptest.NewRecorder()

		e := echo.New()
		ctx := e.NewContext(req, resp)
		ctx.SetPath("/api/v1/orders/tracking/:track_id")
		ctx.SetParamNames("track_id")
		ctx.SetParamValues("ABC-123")

		handler := NewHandler(service)

		// Act
		err := handler.Handle(ctx)

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
