package update

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

func (s *Service) Handle(ctx context.Context, request UpdateOrderDto) (*entity.Order, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	order, err := s.repository.GetByID(ctx, request.UUID)
	if err != nil {
		return nil, err
	}

	order.State = entity.State(request.State)
	order.StateUpdatedAt = s.timeProvider.GetTime()
	order.UpdatedAt = s.timeProvider.GetTime()

	shouldUpdateItems := len(request.Items) > 0

	if err := s.repository.Update(ctx, &order, shouldUpdateItems); err != nil {
		return nil, err
	}

	return &order, nil
}
