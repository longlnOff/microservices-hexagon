package service

import (
	"context"
	"github.com/longlnOff/microservices-hexagon/order/internal/core/domain"
	"github.com/longlnOff/microservices-hexagon/order/internal/port"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

/*
 * OrderService implements port.OrderService interface
 * and provides an access to the category repository & cache service
 */
 type OrderService struct {
	repo  port.OrderRepository
	cache port.CacheRepository
	payment port.PaymentPort
}

// NewOrderService creates a new category service instance
func NewOrderService(repo port.OrderRepository, cache port.CacheRepository, payment port.PaymentPort) *OrderService {
	return &OrderService{
		repo,
		cache,
		payment,
	}
}

func (s *OrderService) PlaceOrder(ctx context.Context, order domain.Order) (domain.Order, error) {
	err := s.repo.Save(ctx, &order)
	if err != nil {
		return domain.Order{}, err
	}
	paymentErr := s.payment.Charge(ctx, &order)
	if paymentErr != nil {
		st, _ := status.FromError(paymentErr)
		fieldErr := errdetails.BadRequest_FieldViolation{
			Field:       "payment",
			Description: st.Message(),
		}
		badReq := &errdetails.BadRequest{}
		badReq.FieldViolations = append(badReq.FieldViolations, &fieldErr)
		orderStatus := status.New(codes.InvalidArgument, "order creation failed")
		statusWithDetails, _:= orderStatus.WithDetails(badReq)
		return domain.Order{}, statusWithDetails.Err()
	}

	return order, nil
}
