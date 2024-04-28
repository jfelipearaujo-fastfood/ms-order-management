package create

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/common"
	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
	"github.com/jfelipearaujo-org/ms-order-management/internal/shared/errors"
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

func (s *Service) Handle(ctx context.Context, request CreateOrderDto) (*entity.Order, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	filter := repository.GetAllOrdersFilter{
		CustomerID: request.CustomerID,
		StateFrom:  entity.Created,
		StateTo:    entity.Delivered,
	}

	count, _, err := s.repository.GetAll(ctx, common.Pagination{}, filter)
	if err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.ErrOrderAlreadyExists
	}

	order := entity.NewOrder(request.CustomerID, s.timeProvider.GetTime())

	if err := s.repository.Create(ctx, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
