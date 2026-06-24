package routes

import (
	"net/http"

	"laundry-system/internal/handlers"
	"laundry-system/internal/middleware"
	"laundry-system/internal/repository"
)

func Setup(
	authHandler *handlers.AuthHandler,
	catalogHandler *handlers.CatalogHandler,
	customerHandler *handlers.CustomerHandler,
	orderHandler *handlers.OrderHandler,
	userHandler *handlers.UserHandler,
	userRepo *repository.UserRepository,
	roleRepo *repository.RoleRepository,
) http.Handler {
	mux := http.NewServeMux()

	// Public routes
	registerAuthRoutes(mux, authHandler)
	registerCatalogRoutes(mux, catalogHandler)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Protected routes
	protected := http.NewServeMux()
	registerCustomerRoutes(protected, customerHandler, orderHandler)
	registerOrderRoutes(protected, orderHandler)
	registerUserRoutes(protected, userHandler)

	jwtAuth := middleware.JWTAuth(userRepo, roleRepo)
	mux.Handle("/api/", jwtAuth(protected))

	return middleware.CORS(mux)
}

// withPermission wraps a handler with RequirePermission middleware.
func withPermission(resource, action string, fn http.HandlerFunc) http.HandlerFunc {
	return middleware.RequirePermission(resource, action)(http.HandlerFunc(fn)).ServeHTTP
}
