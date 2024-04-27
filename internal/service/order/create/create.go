package create

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/provider"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
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

	order := entity.NewOrder(request.CustomerID, s.timeProvider.GetTime())

	if err := s.repository.Create(ctx, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
