package services

import (
	"fmt"

	"laundry-system/internal/models"
	"laundry-system/internal/repository"
)

type CustomerService struct {
	customerRepo *repository.CustomerRepository
}

func NewCustomerService(customerRepo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

func (s *CustomerService) Create(c *models.Customer) error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Phone == "" {
		return fmt.Errorf("phone is required")
	}

	// Check for duplicate phone
	existing, _ := s.customerRepo.GetByPhone(c.Phone)
	if existing != nil {
		return fmt.Errorf("a customer with phone %s already exists", c.Phone)
	}

	return s.customerRepo.Create(c)
}

func (s *CustomerService) GetByID(id string) (*models.Customer, error) {
	return s.customerRepo.GetByID(id)
}

func (s *CustomerService) List(search string) ([]models.Customer, error) {
	return s.customerRepo.List(search)
}

func (s *CustomerService) Update(id string, updates *models.Customer) (*models.Customer, error) {
	if updates.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if updates.Phone == "" {
		return nil, fmt.Errorf("phone is required")
	}

	customer, err := s.customerRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check duplicate phone only if it changed
	if customer.Phone != updates.Phone {
		existing, _ := s.customerRepo.GetByPhone(updates.Phone)
		if existing != nil {
			return nil, fmt.Errorf("a customer with phone %s already exists", updates.Phone)
		}
	}

	customer.Name = updates.Name
	customer.Phone = updates.Phone
	customer.Email = updates.Email
	customer.Address = updates.Address
	customer.Notes = updates.Notes

	if err := s.customerRepo.Update(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) Delete(id string) error {
	return s.customerRepo.Delete(id)
}
