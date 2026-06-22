package services

import (
	"fmt"

	"laundry-system/internal/models"
	"laundry-system/internal/repository"
	"laundry-system/internal/utils"
)

type AuthService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewAuthService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *AuthService {
	return &AuthService{userRepo: userRepo, roleRepo: roleRepo}
}

func (s *AuthService) Login(email, password string) (string, *models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return "", nil, fmt.Errorf("account is inactive")
	}

	if !utils.CheckPassword(user.Password, password) {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Load role
	role, err := s.roleRepo.GetByID(user.RoleID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to load user role: %w", err)
	}
	user.Role = role

	// Load permissions
	permissions, err := s.roleRepo.GetPermissionsByRoleID(user.RoleID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to load user permissions: %w", err)
	}
	user.Permissions = permissions

	token, err := utils.GenerateToken(user.ID, user.Email, user.RoleID, user.Role.Name)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Update last login (non-fatal)
	_ = s.userRepo.UpdateLastLogin(user.ID)

	return token, user, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*utils.Claims, error) {
	return utils.ValidateToken(tokenString)
}
