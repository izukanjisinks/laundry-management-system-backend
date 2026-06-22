package middleware

import (
	"net/http"

	"laundry-system/internal/models"
	"laundry-system/internal/utils"
)

func RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(ContextKeyUser).(*models.User)
			if !ok || user == nil {
				utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			if !user.HasPermission(resource, action) {
				utils.RespondError(w, http.StatusForbidden, "Forbidden: insufficient permissions")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequireAnyPermission(permissions ...[2]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(ContextKeyUser).(*models.User)
			if !ok || user == nil {
				utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			for _, p := range permissions {
				if user.HasPermission(p[0], p[1]) {
					next.ServeHTTP(w, r)
					return
				}
			}
			utils.RespondError(w, http.StatusForbidden, "Forbidden: insufficient permissions")
		})
	}
}
