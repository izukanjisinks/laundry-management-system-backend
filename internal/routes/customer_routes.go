package routes

import (
	"net/http"

	"laundry-system/internal/handlers"
)

func registerCustomerRoutes(mux *http.ServeMux, customerHandler *handlers.CustomerHandler, orderHandler *handlers.OrderHandler) {
	mux.HandleFunc("GET /api/customers", withPermission("customers", "read", customerHandler.List))
	mux.HandleFunc("POST /api/customers", withPermission("customers", "create", customerHandler.Create))
	mux.HandleFunc("GET /api/customers/{id}", withPermission("customers", "read", customerHandler.GetByID))
	mux.HandleFunc("PUT /api/customers/{id}", withPermission("customers", "update", customerHandler.Update))
	mux.HandleFunc("DELETE /api/customers/{id}", withPermission("customers", "delete", customerHandler.Delete))
	mux.HandleFunc("GET /api/customers/{id}/orders", withPermission("orders", "read", orderHandler.ListByCustomer))
}
