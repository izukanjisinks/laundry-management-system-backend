package middleware

import (
	"context"
	"net/http"
	"strings"

	"laundry-system/internal/repository"
	"laundry-system/internal/utils"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyUser   contextKey = "user"
	ContextKeyRole   contextKey = "role"
)

func JWTAuth(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.RespondError(w, http.StatusUnauthorized, "Missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			claims, err := utils.ValidateToken(parts[1])
			if err != nil {
				utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired token")
				return
			}

			// Fetch full user to check current active status
			user, err := userRepo.GetByID(claims.UserID)
			if err != nil {
				utils.RespondError(w, http.StatusUnauthorized, "User not found")
				return
			}

			if !user.IsActive {
				utils.RespondError(w, http.StatusUnauthorized, "Account is inactive")
				return
			}

			// Load role and permissions
			role, err := roleRepo.GetByID(user.RoleID)
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to load user role")
				return
			}
			user.Role = role

			permissions, err := roleRepo.GetPermissionsByRoleID(user.RoleID)
			if err != nil {
				utils.RespondError(w, http.StatusInternalServerError, "Failed to load user permissions")
				return
			}
			user.Permissions = permissions

			ctx := context.WithValue(r.Context(), ContextKeyUserID, user.ID)
			ctx = context.WithValue(ctx, ContextKeyUser, user)
			ctx = context.WithValue(ctx, ContextKeyRole, user.Role.Name)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
