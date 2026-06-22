package handlers

import (
	"net/http"

	"laundry-system/internal/models"
	"laundry-system/internal/services"
	"laundry-system/internal/utils"
)

type CustomerHandler struct {
	customerService *services.CustomerService
}

func NewCustomerHandler(customerService *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (h *CustomerHandler) List(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	customers, err := h.customerService.List(search)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if customers == nil {
		customers = []models.Customer{}
	}
	utils.RespondJSON(w, http.StatusOK, customers)
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	if err := utils.DecodeJSON(r, &customer); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.customerService.Create(&customer); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusCreated, customer)
}

func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	customer, err := h.customerService.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updates models.Customer
	if err := utils.DecodeJSON(r, &updates); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	customer, err := h.customerService.Update(id, &updates)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, customer)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.customerService.Delete(id); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Customer deleted successfully"})
}
