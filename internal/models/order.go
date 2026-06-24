package models

import "time"

type OrderStatus string
type ServiceType string
type PaymentStatus string
type PaymentMethod string

const (
	OrderStatusReceived OrderStatus = "received"
	OrderStatusWashing  OrderStatus = "washing"
	OrderStatusReady    OrderStatus = "ready"
	OrderStatusPickedUp OrderStatus = "picked_up"

	ServiceTypeWashFold  ServiceType = "wash_fold"
	ServiceTypeDryClean  ServiceType = "dry_clean"
	ServiceTypeIroning   ServiceType = "ironing"
	ServiceTypeWashIron  ServiceType = "wash_iron"

	PaymentStatusUnpaid  PaymentStatus = "unpaid"
	PaymentStatusPartial PaymentStatus = "partial"
	PaymentStatusPaid    PaymentStatus = "paid"

	PaymentMethodCash     PaymentMethod = "cash"
	PaymentMethodCard     PaymentMethod = "card"
	PaymentMethodTransfer PaymentMethod = "transfer"
)

// CanTransitionTo checks whether transitioning from the current status to target is valid.
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	next, ok := s.NextStatus()
	return ok && next == target
}

// NextStatus returns the valid next step in the progression.
func (s OrderStatus) NextStatus() (OrderStatus, bool) {
	switch s {
	case OrderStatusReceived:
		return OrderStatusWashing, true
	case OrderStatusWashing:
		return OrderStatusReady, true
	case OrderStatusReady:
		return OrderStatusPickedUp, true
	default:
		return "", false
	}
}

type OrderItem struct {
	Name  string  `json:"name"`
	Qty   int     `json:"qty"`
	Price float64 `json:"price"`
}

type Order struct {
	ID            string        `json:"id"`
	OrderNumber   int           `json:"order_number"`
	CustomerID    string        `json:"customer_id"`
	Customer      *Customer     `json:"customer,omitempty"`
	CreatedBy     string        `json:"created_by"`
	Creator       *User         `json:"creator,omitempty"`
	Status        OrderStatus   `json:"status"`
	ServiceType   ServiceType   `json:"service_type"`
	Items         []OrderItem   `json:"items"`
	Subtotal      float64       `json:"subtotal"`
	TaxRate       float64       `json:"tax_rate"`
	TaxAmount     float64       `json:"tax_amount"`
	TotalPrice    float64       `json:"total_price"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
	Notes         string        `json:"notes,omitempty"`
	DueAt         *time.Time    `json:"due_at,omitempty"`
	ReceivedAt    time.Time     `json:"received_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	PickedUpAt    *time.Time    `json:"picked_up_at,omitempty"`
}
