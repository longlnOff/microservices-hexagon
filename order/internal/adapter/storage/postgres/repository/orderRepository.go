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


// Flow
// 1. Create transaction 
// 2. Create order in orders table
// 3. Create order items in order_items table
// 4. Commit transaction
// 5. rollback when error
func (o *OrderRepository) Save(ctx context.Context, order *domain.Order) error {
	querySaveOrder := `
		INSERT INTO orders (customer_id, status, created_at)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	querySaveOrderItem := `
		INSERT INTO order_items (order_id, product_code, quantity, unit_price, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()
	// start transaction
	err := pgx.BeginFunc(ctx, o.DB.DB, func(tx pgx.Tx) error {
		// 2. Create order in orders table
		err := tx.QueryRow(
							ctx, 
							querySaveOrder,
							order.CustomerID,
							order.Status,
							order.CreatedAt,	
						).Scan(
							&order.ID,
						)
		if err != nil {
			return err
		}
		
		// 3. Create order items in order_items table
		for _, orderItem := range order.OrderItems {
			_, err := tx.Exec(ctx, 
							querySaveOrderItem, 
							order.ID, 
							orderItem.ProductCode, 
							orderItem.Quantity, 
							orderItem.UnitPrice, 
							time.Now(),
						)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		if errCode := o.DB.ErrorCode(err); errCode == pgerrcode.UniqueViolation {
			return domain.ErrConflictingData
		}
		return err
	}

	return nil
}

// Flow
// 1. Get order from orders table with id
// 2. Get all order items from order_items table with order id
// 3. transform data to store in domain.Order
func (o *OrderRepository) Get(ctx context.Context, id int) (domain.Order, error) {
	queryGetOrder := `
					SELECT id, customer_id, status, created_at
					FROM orders
					WHERE id = $1
				`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var order domain.Order
	err := o.DB.DB.QueryRow(ctx, queryGetOrder, id).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Order{}, domain.ErrDataNotFound
		}
		return domain.Order{}, err
	}

	var orderItems []domain.OrderItem
	queryGetOrderItems := `
					SELECT product_code, unit_price, quantity FROM order_items
					WHERE order_id = $1
					ORDER BY created_at
				`
	rows, err := o.DB.DB.Query(ctx, queryGetOrderItems, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Order{}, domain.ErrDataNotFound
		}
		return domain.Order{}, err
	}
	for rows.Next() {
		var orderItem domain.OrderItem
		err := rows.Scan(
			&orderItem.ProductCode,
			&orderItem.UnitPrice,
			&orderItem.Quantity,
		)
		if err != nil {
			return domain.Order{}, err
		}
		orderItems = append(orderItems, orderItem)
	}

	order.OrderItems = orderItems
	return order, nil
}
