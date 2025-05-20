package port

import (
	"context"

	"github.com/longlnOff/microservices-hexagon/order/internal/core/domain"
)

type PaymentPort interface {
	Charge(context.Context, *domain.Order) error
}
