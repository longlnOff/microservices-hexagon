package grpc

import (
	"context"

	"github.com/longlnOff/microservices-hexagon/order/internal/core/domain"
	"github.com/longlnOff/microservices-proto/golang/order"
)


func (g GRPCServer) Create(ctx context.Context, request *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	var orderItems []domain.OrderItem
	for _, orderItem := range request.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	newOrder := domain.NewOrder(request.UserId, orderItems)
	result, err := g.service.PlaceOrder(ctx, newOrder)
	if err != nil {
		return nil, err
	}
	return &order.CreateOrderResponse{OrderId: result.ID}, nil
}


