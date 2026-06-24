package services

import (
	"fmt"
	"log"
	"math"
	"strings"

	"laundry-system/internal/models"
	"laundry-system/internal/repository"
	"laundry-system/internal/utils/email"
)

type OrderService struct {
	orderRepo    *repository.OrderRepository
	customerRepo *repository.CustomerRepository
	emailSvc     *email.EmailService
}

func NewOrderService(orderRepo *repository.OrderRepository, customerRepo *repository.CustomerRepository, emailSvc *email.EmailService) *OrderService {
	return &OrderService{orderRepo: orderRepo, customerRepo: customerRepo, emailSvc: emailSvc}
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

	// Calculate subtotal, tax, total
	o.Subtotal = calculateSubtotal(o.Items)
	if o.TaxRate == 0 {
		o.TaxRate = 0.075
	}
	o.TaxAmount = roundTo2(o.Subtotal * o.TaxRate)
	o.TotalPrice = roundTo2(o.Subtotal + o.TaxAmount)

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
			string(models.OrderStatusReady):     true,
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
	order.ServiceType = updates.ServiceType
	order.DueAt = updates.DueAt
	order.Subtotal = calculateSubtotal(updates.Items)
	if order.TaxRate == 0 {
		order.TaxRate = 0.075
	}
	order.TaxAmount = roundTo2(order.Subtotal * order.TaxRate)
	order.TotalPrice = roundTo2(order.Subtotal + order.TaxAmount)

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

func (s *OrderService) UpdatePayment(id string, paymentStatus models.PaymentStatus, paymentMethod models.PaymentMethod) error {
	validStatus := map[models.PaymentStatus]bool{
		models.PaymentStatusUnpaid:  true,
		models.PaymentStatusPartial: true,
		models.PaymentStatusPaid:    true,
	}
	if !validStatus[paymentStatus] {
		return fmt.Errorf("invalid payment_status: %s", paymentStatus)
	}
	if paymentStatus != models.PaymentStatusUnpaid && paymentMethod == "" {
		return fmt.Errorf("payment_method is required when payment_status is %s", paymentStatus)
	}
	if _, err := s.orderRepo.GetByID(id); err != nil {
		return err
	}
	return s.orderRepo.UpdatePayment(id, paymentStatus, paymentMethod)
}

func (s *OrderService) Summary() (map[string]interface{}, error) {
	return s.orderRepo.Summary()
}

func calculateSubtotal(items []models.OrderItem) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Qty) * item.Price
	}
	return roundTo2(total)
}

func roundTo2(v float64) float64 {
	return math.Round(v*100) / 100
}

func mustNextStatus(s models.OrderStatus) models.OrderStatus {
	next, _ := s.NextStatus()
	return next
}

// NotifyReady sends a pickup-ready email to the customer for the given order.
func (s *OrderService) NotifyReady(orderID string) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return fmt.Errorf("order not found")
	}
	if order.Status != models.OrderStatusReady {
		return fmt.Errorf("order is not ready for pickup")
	}
	if order.Customer == nil || order.Customer.Email == "" {
		return fmt.Errorf("customer has no email address on file")
	}

	serviceLabels := map[models.ServiceType]string{
		models.ServiceTypeWashFold: "Wash & Fold",
		models.ServiceTypeDryClean: "Dry Clean",
		models.ServiceTypeIroning:  "Ironing",
		models.ServiceTypeWashIron: "Wash & Iron",
	}
	serviceLabel := serviceLabels[order.ServiceType]
	if serviceLabel == "" {
		serviceLabel = strings.ToUpper(string(order.ServiceType))
	}

	type item = struct {
		Name string
		Qty  int
	}
	items := make([]item, len(order.Items))
	for i, it := range order.Items {
		items[i] = item{Name: it.Name, Qty: it.Qty}
	}

	orderNumber := fmt.Sprintf("WP-%d", order.OrderNumber)
	body := email.OrderReadyTemplate(order.Customer.Name, orderNumber, serviceLabel, items)

	go func() {
		if err := s.emailSvc.SendEmail(
			[]string{order.Customer.Email},
			fmt.Sprintf("Your laundry is ready for pickup — %s", orderNumber),
			body,
		); err != nil {
			log.Printf("[EMAIL] Failed to send pickup-ready email for order %s: %v", orderID, err)
		}
	}()

	return nil
}
