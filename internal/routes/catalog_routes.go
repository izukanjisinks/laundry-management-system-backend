package routes

import (
	"net/http"

	"laundry-system/internal/handlers"
)

func registerCatalogRoutes(mux *http.ServeMux, catalogHandler *handlers.CatalogHandler) {
	mux.HandleFunc("/api/catalog", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			catalogHandler.List(w, r)
			return
		}
		http.NotFound(w, r)
	})
}
