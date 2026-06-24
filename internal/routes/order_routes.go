package routes

import (
	"net/http"

	"laundry-system/internal/handlers"
)

func registerOrderRoutes(mux *http.ServeMux, orderHandler *handlers.OrderHandler) {
	mux.HandleFunc("GET /api/orders", withPermission("orders", "read", orderHandler.List))
	mux.HandleFunc("POST /api/orders", withPermission("orders", "create", orderHandler.Create))
	mux.HandleFunc("GET /api/orders/{id}", withPermission("orders", "read", orderHandler.GetByID))
	mux.HandleFunc("PUT /api/orders/{id}", withPermission("orders", "update", orderHandler.Update))
	mux.HandleFunc("PATCH /api/orders/{id}/status", withPermission("orders", "update_status", orderHandler.UpdateStatus))
	mux.HandleFunc("PATCH /api/orders/{id}/payment", withPermission("orders", "update", orderHandler.UpdatePayment))
	mux.HandleFunc("POST /api/orders/{id}/notify-ready", withPermission("orders", "update", orderHandler.NotifyReady))
	mux.HandleFunc("DELETE /api/orders/{id}", withPermission("orders", "delete", orderHandler.Delete))

	mux.HandleFunc("GET /api/reports/summary", withPermission("reports", "read", orderHandler.Summary))
}
