package payment

import (
	"context"

	"github.com/longlnOff/microservices-proto/golang/payment"
	"github.com/longlnOff/microservices-hexagon/order/internal/core/domain"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentAdapter struct {
	payment payment.PaymentClient
}

func NewAdapter(paymentServiceUrl string) (*PaymentAdapter, error) {
	var opts []grpc.DialOption
	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	conn, err := grpc.NewClient(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	client := payment.NewPaymentClient(conn)
	return &PaymentAdapter{payment: client}, nil
}


func (a *PaymentAdapter) Charge(ctx context.Context, order *domain.Order) error {
	_, err := a.payment.Create(ctx,
		&payment.CreatePaymentRequest{
			UserId:     order.CustomerID,
			OrderId:    order.ID,
			TotalPrice: order.TotalPrice(),
		})
	return err
}
