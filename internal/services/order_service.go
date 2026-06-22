package services

import (
	"fmt"

	"laundry-system/internal/models"
	"laundry-system/internal/repository"
)

type OrderService struct {
	orderRepo    *repository.OrderRepository
	customerRepo *repository.CustomerRepository
}

func NewOrderService(orderRepo *repository.OrderRepository, customerRepo *repository.CustomerRepository) *OrderService {
	return &OrderService{orderRepo: orderRepo, customerRepo: customerRepo}
}

func (s *OrderService) Create(o *models.Order) error {
	if o.CustomerID == "" {
		return fmt.Errorf("customer_id is required")
	}
	if len(o.Items) == 0 {
		return fmt.Errorf("order must have at least one item")
	}
	for i, item := range o.Items {
		if item.Name == "" {
			return fmt.Errorf("item %d: name is required", i+1)
		}
		if item.Qty <= 0 {
			return fmt.Errorf("item %d: qty must be greater than zero", i+1)
		}
		if item.Price < 0 {
			return fmt.Errorf("item %d: price cannot be negative", i+1)
		}
	}

	// Verify customer exists
	if _, err := s.customerRepo.GetByID(o.CustomerID); err != nil {
		return fmt.Errorf("customer not found")
	}

	// Calculate total if not set
	if o.TotalPrice == 0 {
		o.TotalPrice = calculateTotal(o.Items)
	}

	o.Status = models.OrderStatusReceived
	return s.orderRepo.Create(o)
}

func (s *OrderService) GetByID(id string) (*models.Order, error) {
	return s.orderRepo.GetByID(id)
}

func (s *OrderService) List(status string) ([]models.Order, error) {
	// Validate status filter if provided
	if status != "" {
		valid := map[string]bool{
			string(models.OrderStatusReceived): true,
			string(models.OrderStatusWashing):  true,
			string(models.OrderStatusDone):     true,
			string(models.OrderStatusPickedUp): true,
		}
		if !valid[status] {
			return nil, fmt.Errorf("invalid status filter: %s", status)
		}
	}
	return s.orderRepo.List(status)
}

func (s *OrderService) ListByCustomer(customerID string) ([]models.Order, error) {
	if _, err := s.customerRepo.GetByID(customerID); err != nil {
		return nil, fmt.Errorf("customer not found")
	}
	return s.orderRepo.ListByCustomer(customerID)
}

func (s *OrderService) UpdateStatus(id string, newStatus models.OrderStatus) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if !order.Status.CanTransitionTo(newStatus) {
		return nil, fmt.Errorf(
			"invalid status transition: %s → %s (expected: %s → %s)",
			order.Status, newStatus,
			order.Status, mustNextStatus(order.Status),
		)
	}

	if err := s.orderRepo.UpdateStatus(id, newStatus); err != nil {
		return nil, err
	}

	order.Status = newStatus
	return order, nil
}

func (s *OrderService) Update(id string, updates *models.Order) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Only allow editing orders that haven't been picked up
	if order.Status == models.OrderStatusPickedUp {
		return nil, fmt.Errorf("cannot edit a completed order")
	}

	if len(updates.Items) == 0 {
		return nil, fmt.Errorf("order must have at least one item")
	}

	order.Items = updates.Items
	order.Notes = updates.Notes
	order.TotalPrice = calculateTotal(updates.Items)

	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) Delete(id string) error {
	order, err := s.orderRepo.GetByID(id)
	if err != nil {
		return err
	}
	if order.Status != models.OrderStatusReceived {
		return fmt.Errorf("only orders with status 'received' can be deleted")
	}
	return s.orderRepo.Delete(id)
}

func (s *OrderService) Summary() (map[string]interface{}, error) {
	return s.orderRepo.Summary()
}

func calculateTotal(items []models.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Qty) * item.Price
	}
	return total
}

func mustNextStatus(s models.OrderStatus) models.OrderStatus {
	next, _ := s.NextStatus()
	return next
}
