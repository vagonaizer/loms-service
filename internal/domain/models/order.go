package models

import "time"

type OrderStatus string

const (
	OrderStatusNew             OrderStatus = "new"
	OrderStatusAwaitingPayment OrderStatus = "awaiting payment"
	OrderStatusFailed          OrderStatus = "failed"
	OrderStatusPayed           OrderStatus = "payed"
	OrderStatusCancelled       OrderStatus = "cancelled"
)

type OrderItem struct {
	SKU   uint32
	Count uint32
}

type Order struct {
	ID        int64
	UserID    int64
	Status    OrderStatus
	Items     []OrderItem
	CreatedAt time.Time
	UpdatedAt time.Time
}
