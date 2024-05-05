package update

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
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

func (s *Service) Handle(ctx context.Context, order *order_entity.Order, request UpdateOrderDto) error {
	if err := request.Validate(); err != nil {
		return err
	}

	if err := order.UpdateState(order_entity.OrderState(request.State), s.timeProvider.GetTime()); err != nil {
		return err
	}

	shouldUpdateItems := len(request.Items) > 0

	if shouldUpdateItems {
		for _, item := range request.Items {
			itemToAdd := order_entity.NewItem(item.ItemId, item.Name, item.UnitPrice, item.Quantity)

			if err := order.AddItem(itemToAdd, s.timeProvider.GetTime()); err != nil {
				return err
			}
		}
	}

	if err := s.repository.Update(ctx, order, shouldUpdateItems); err != nil {
		return err
	}

	order.RefreshStateTitle()

	return nil
}
