package routes

import (
	"net/http"

	"laundry-system/internal/handlers"
)

func registerAuthRoutes(mux *http.ServeMux, authHandler *handlers.AuthHandler) {
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
}
