package handlers

import (
	"net/http"

	"laundry-system/internal/middleware"
	"laundry-system/internal/models"
	"laundry-system/internal/services"
	"laundry-system/internal/utils"
)

type OrderHandler struct {
	orderService *services.OrderService
}

func NewOrderHandler(orderService *services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	orders, err := h.orderService.List(status)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if orders == nil {
		orders = []models.Order{}
	}
	utils.RespondJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(middleware.ContextKeyUser).(*models.User)
	if !ok || user == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var order models.Order
	if err := utils.DecodeJSON(r, &order); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	order.CreatedBy = user.ID
	if err := h.orderService.Create(&order); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	order, err := h.orderService.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updates models.Order
	if err := utils.DecodeJSON(r, &updates); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	order, err := h.orderService.Update(id, &updates)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var body struct {
		Status string `json:"status"`
	}
	if err := utils.DecodeJSON(r, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if body.Status == "" {
		utils.RespondError(w, http.StatusBadRequest, "status is required")
		return
	}

	order, err := h.orderService.UpdateStatus(id, models.OrderStatus(body.Status))
	if err != nil {
		utils.RespondError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) ListByCustomer(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("id")
	orders, err := h.orderService.ListByCustomer(customerID)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if orders == nil {
		orders = []models.Order{}
	}
	utils.RespondJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.orderService.Delete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Order deleted successfully"})
}

func (h *OrderHandler) Summary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.orderService.Summary()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, summary)
}
