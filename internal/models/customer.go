package models

import "time"

type Customer struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email,omitempty"`
	Address     string    `json:"address,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	TotalOrders int       `json:"total_orders,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
