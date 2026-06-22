package models

import "time"

type OrderStatus string

const (
	OrderStatusReceived  OrderStatus = "received"
	OrderStatusWashing   OrderStatus = "washing"
	OrderStatusDone      OrderStatus = "done"
	OrderStatusPickedUp  OrderStatus = "picked_up"
)

// NextStatus returns the valid next status in the progression, and whether a next step exists.
func (s OrderStatus) NextStatus() (OrderStatus, bool) {
	switch s {
	case OrderStatusReceived:
		return OrderStatusWashing, true
	case OrderStatusWashing:
		return OrderStatusDone, true
	case OrderStatusDone:
		return OrderStatusPickedUp, true
	default:
		return "", false
	}
}

// CanTransitionTo checks whether transitioning from the current status to the target is valid.
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	next, ok := s.NextStatus()
	return ok && next == target
}

type OrderItem struct {
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

type Order struct {
	ID         string      `json:"id"`
	CustomerID string      `json:"customer_id"`
	Customer   *Customer   `json:"customer,omitempty"`
	CreatedBy  string      `json:"created_by"`
	Creator    *User       `json:"creator,omitempty"`
	Status     OrderStatus `json:"status"`
	Items      []OrderItem `json:"items"`
	TotalPrice float64     `json:"total_price"`
	Notes      string      `json:"notes,omitempty"`
	ReceivedAt time.Time   `json:"received_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	PickedUpAt *time.Time  `json:"picked_up_at,omitempty"`
}
