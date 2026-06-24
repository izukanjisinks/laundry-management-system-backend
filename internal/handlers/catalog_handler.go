package handlers

import (
	"net/http"

	"laundry-system/internal/services"
	"laundry-system/internal/utils"
)

type CatalogHandler struct {
	catalogService *services.CatalogService
}

func NewCatalogHandler(catalogService *services.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogService: catalogService}
}

func (h *CatalogHandler) List(w http.ResponseWriter, r *http.Request) {
	items, err := h.catalogService.List()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch catalog")
		return
	}
	utils.RespondJSON(w, http.StatusOK, items)
}
