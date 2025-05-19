package repository

import (
	"context"
	"time"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/longlnOff/microservices-hexagon/order/internal/adapter/storage/postgres"
	"github.com/longlnOff/microservices-hexagon/order/internal/core/domain"
)

var (
	QueryTimeoutDuration = 5 * time.Second
)

/**
 * OrderRepository implements port.OrderRepository interface
 * and provides an access to the postgres database
 */
type OrderRepository struct {
	*postgres.DB
}

// NewOrderRepository creates a new category repository instance
func NewOrderRepository(db *postgres.DB) *OrderRepository {
	return &OrderRepository{
		db,
	}
}


func (o *OrderRepository) Save(context.Context, *domain.Order) error {

	return nil
}
func (o *OrderRepository) Get(ctx context.Context, id int) (domain.Order, error) {
	query := `
		SELECT id, customer_id, status, order_items, created_at
		FROM orders
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var order domain.Order
	err := o.DB.QueryRow(ctx, query, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Status,
		&order.OrderItems,
		&order.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Order{}, domain.ErrDataNotFound
		}
		return domain.Order{}, err
	}


	return order, nil
}

