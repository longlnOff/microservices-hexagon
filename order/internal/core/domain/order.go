package domain

import (
	"time"
)

type OrderItem struct {
	ProductCode string 		`json:"product_code"`
	UnitPrice   float32 	`json:"unit_price"`
	Quantity    int32 		`json:"quantity"`
}

type Order struct {
	ID					int64		`json:"id"`
	CustomerID			int64		`json:"customer_id"`
	Status				string		`json:"status"`
	OrderItems			[]OrderItem `json:"order_items"`
	CreatedAt			time.Time	`json:"created_at"`
}

func NewOrder(customerID int64, orderItems []OrderItem) Order {
	return Order{
		CustomerID: customerID,
		Status: "PENDING",
		OrderItems: orderItems,
		CreatedAt: time.Now(),
	}
}


func (o *Order) TotalPrice() float32 {
	var totalPrice float32
	for _, orderItem := range o.OrderItems {
		totalPrice += orderItem.UnitPrice * float32(orderItem.Quantity)
	}
	return totalPrice
}
