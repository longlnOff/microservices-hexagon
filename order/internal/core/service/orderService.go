package service

import (
	"context"
	"github.com/longlnOff/microservices-hexagon/order/internal/core/domain"
	"github.com/longlnOff/microservices-hexagon/order/internal/port"
)

/*
 * OrderService implements port.OrderService interface
 * and provides an access to the category repository & cache service
 */
 type OrderService struct {
	repo  port.OrderRepository
	cache port.CacheRepository
}

// NewOrderService creates a new category service instance
func NewOrderService(repo port.OrderRepository, cache port.CacheRepository) *OrderService {
	return &OrderService{
		repo,
		cache,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	err := s.repo.Save(ctx, &order)
	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}
