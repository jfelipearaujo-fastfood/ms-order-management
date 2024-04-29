package create

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/common"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/custom_error"
)

type Service struct {
	repository   repository.OrderRepository
	timeProvider provider.TimeProvider
}

func NewService(
	repository repository.OrderRepository,
	timeProvider provider.TimeProvider,
) *Service {
	return &Service{
		repository:   repository,
		timeProvider: timeProvider,
	}
}

func (s *Service) Handle(ctx context.Context, request CreateOrderDto) (*order_entity.Order, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	filter := repository.GetAllOrdersFilter{
		CustomerID: request.CustomerID,
		StateFrom:  order_entity.Created,
		StateTo:    order_entity.Delivered,
	}

	count, _, err := s.repository.GetAll(ctx, common.Pagination{}, filter)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, custom_error.ErrOrderAlreadyExists
	}

	order := order_entity.NewOrder(request.CustomerID, s.timeProvider.GetTime())

	if err := s.repository.Create(ctx, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
