package get_all

import (
	"context"

	"github.com/jfelipearaujo-org/ms-order-management/internal/entity/order_entity"
	"github.com/jfelipearaujo-org/ms-order-management/internal/repository"
)

type Service struct {
	repository repository.OrderRepository
}

func NewService(repository repository.OrderRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Handle(ctx context.Context, request GetOrdersDto) (int, []order_entity.Order, error) {
	request.SetDefaults()

	if err := request.Validate(); err != nil {
		return 0, nil, err
	}

	filter := repository.GetAllOrdersFilter{
		CustomerID: request.CustomerID,
		StateFrom:  order_entity.OrderState(request.State),
		StateTo:    order_entity.OrderState(request.State),
	}

	count, orders, err := s.repository.GetAll(ctx, request.Pagination, filter)
	if err != nil {
		return 0, nil, err
	}

	return count, orders, nil
}
