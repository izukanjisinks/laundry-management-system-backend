package services

import (
	"fmt"

	"laundry-system/internal/models"
	"laundry-system/internal/repository"
	"laundry-system/internal/utils"
)

type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *UserService) Create(u *models.User) error {
	if u.FullName == "" {
		return fmt.Errorf("full_name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	if u.Password == "" {
		return fmt.Errorf("password is required")
	}
	if u.RoleID == "" {
		return fmt.Errorf("role_id is required")
	}

	// Verify role exists
	if _, err := s.roleRepo.GetByID(u.RoleID); err != nil {
		return fmt.Errorf("role not found")
	}

	// Check duplicate email
	existing, _ := s.userRepo.GetByEmail(u.Email)
	if existing != nil {
		return fmt.Errorf("a user with email %s already exists", u.Email)
	}

	hashed, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashed
	u.IsActive = true

	return s.userRepo.Create(u)
}

func (s *UserService) GetByID(id string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	role, err := s.roleRepo.GetByID(user.RoleID)
	if err != nil {
		return nil, err
	}
	user.Role = role
	return user, nil
}

func (s *UserService) List() ([]models.User, error) {
	return s.userRepo.List()
}

func (s *UserService) Update(id string, updates *models.User) (*models.User, error) {
	if updates.FullName == "" {
		return nil, fmt.Errorf("full_name is required")
	}
	if updates.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if updates.RoleID == "" {
		return nil, fmt.Errorf("role_id is required")
	}

	// Verify role exists
	if _, err := s.roleRepo.GetByID(updates.RoleID); err != nil {
		return nil, fmt.Errorf("role not found")
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check duplicate email only if it changed
	if user.Email != updates.Email {
		existing, _ := s.userRepo.GetByEmail(updates.Email)
		if existing != nil {
			return nil, fmt.Errorf("a user with email %s already exists", updates.Email)
		}
	}

	user.FullName = updates.FullName
	user.Email = updates.Email
	user.RoleID = updates.RoleID
	user.IsActive = updates.IsActive

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	role, _ := s.roleRepo.GetByID(user.RoleID)
	user.Role = role
	return user, nil
}

func (s *UserService) UpdatePassword(id, newPassword string) error {
	if newPassword == "" {
		return fmt.Errorf("password is required")
	}
	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(id, hashed)
}

func (s *UserService) Delete(id string, requestingUserID string) error {
	if id == requestingUserID {
		return fmt.Errorf("cannot deactivate your own account")
	}
	return s.userRepo.Delete(id)
}

func (s *UserService) ListRoles() ([]models.Role, error) {
	return s.roleRepo.List()
}
