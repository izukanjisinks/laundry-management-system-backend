package routes

import (
	"net/http"

	"laundry-system/internal/handlers"
)

func registerUserRoutes(mux *http.ServeMux, userHandler *handlers.UserHandler) {
	mux.HandleFunc("GET /api/users", withPermission("users", "read", userHandler.List))
	mux.HandleFunc("POST /api/users", withPermission("users", "create", userHandler.Create))
	mux.HandleFunc("GET /api/users/{id}", withPermission("users", "read", userHandler.GetByID))
	mux.HandleFunc("PUT /api/users/{id}", withPermission("users", "update", userHandler.Update))
	mux.HandleFunc("PATCH /api/users/{id}/password", withPermission("users", "update", userHandler.UpdatePassword))
	mux.HandleFunc("DELETE /api/users/{id}", withPermission("users", "delete", userHandler.Delete))

	mux.HandleFunc("GET /api/roles", withPermission("users", "read", userHandler.ListRoles))
}
