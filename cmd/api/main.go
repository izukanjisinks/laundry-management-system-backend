package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"laundry-system/internal/config"
	"laundry-system/internal/database"
	"laundry-system/internal/handlers"
	"laundry-system/internal/repository"
	"laundry-system/internal/routes"
	"laundry-system/internal/services"
	"laundry-system/internal/utils/email"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting Laundry Management System in %s mode\n", cfg.Environment)

	// Connect to database
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	if err := database.Connect(connStr); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("✓ Database connected")

	// Run migrations
	if err := database.RunMigrations(database.MigrationsDir()); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("✓ Migrations completed")

	// Exit after migrations if --migrate-only flag is passed
	if len(os.Args) > 1 && os.Args[1] == "--migrate-only" {
		log.Println("Migration-only mode — exiting")
		return
	}

	// Email service
	emailSvc := email.NewEmailService(&cfg.Email)

	// Repositories
	userRepo := repository.NewUserRepository()
	roleRepo := repository.NewRoleRepository()
	customerRepo := repository.NewCustomerRepository()
	orderRepo := repository.NewOrderRepository()
	catalogRepo := repository.NewCatalogRepository()

	// Services
	authService := services.NewAuthService(userRepo, roleRepo)
	catalogService := services.NewCatalogService(catalogRepo)
	customerService := services.NewCustomerService(customerRepo, emailSvc)
	orderService := services.NewOrderService(orderRepo, customerRepo, emailSvc)
	userService := services.NewUserService(userRepo, roleRepo)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	catalogHandler := handlers.NewCatalogHandler(catalogService)
	customerHandler := handlers.NewCustomerHandler(customerService)
	orderHandler := handlers.NewOrderHandler(orderService)
	userHandler := handlers.NewUserHandler(userService)

	// Router
	router := routes.Setup(authHandler, catalogHandler, customerHandler, orderHandler, userHandler, userRepo, roleRepo)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("✓ Server listening on http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
