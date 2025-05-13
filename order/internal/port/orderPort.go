package port

import (
	"context"

	"github.com/longlnOff/microservices-hexagon/order/internal/core/domain"
)

//go:generate mockgen -source=category.go -destination=mock/category.go -package=mock

// CategoryRepository is an interface for interacting with category-related data
type OrderRepository interface {
	Get(ctx context.Context, id string) (domain.Order, error)
	Save(context.Context, *domain.Order) error
}
// OrderService is an interface for interacting with category-related business logic
type OrderService interface {
	PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error)
}
