package main

import (
	"fmt"
	"log"

	"laundry-system/internal/config"
	"laundry-system/internal/database"
	"laundry-system/internal/models"
	"laundry-system/internal/repository"
	"laundry-system/internal/utils"
)

func main() {
	cfg := config.Load()

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	if err := database.Connect(connStr); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	roleRepo := repository.NewRoleRepository()
	userRepo := repository.NewUserRepository()

	// Check if admin user already exists
	existing, _ := userRepo.GetByEmail(cfg.AdminEmail)
	if existing != nil {
		log.Printf("Admin user already exists: %s\n", cfg.AdminEmail)
		return
	}

	// Get admin role
	adminRole, err := roleRepo.GetByName(models.RoleAdmin)
	if err != nil {
		log.Fatalf("Admin role not found — run migrations first: %v", err)
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(cfg.AdminPassword)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Create admin user
	admin := &models.User{
		FullName: "System Admin",
		Email:    cfg.AdminEmail,
		Password: hashedPassword,
		RoleID:   adminRole.ID,
		IsActive: true,
	}

	if err := userRepo.Create(admin); err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	log.Printf("✓ Admin user created: %s\n", cfg.AdminEmail)
	log.Printf("  Password: %s\n", cfg.AdminPassword)
	log.Println("  (Change this password after first login)")
}
