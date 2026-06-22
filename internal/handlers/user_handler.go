package handlers

import (
	"net/http"

	"laundry-system/internal/middleware"
	"laundry-system/internal/models"
	"laundry-system/internal/services"
	"laundry-system/internal/utils"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.userService.List()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if users == nil {
		users = []models.User{}
	}
	utils.RespondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := utils.DecodeJSON(r, &user); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.userService.Create(&user); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	user.Password = ""
	utils.RespondJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	user, err := h.userService.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var updates models.User
	if err := utils.DecodeJSON(r, &updates); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	user, err := h.userService.Update(id, &updates)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, user)
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var body struct {
		Password string `json:"password"`
	}
	if err := utils.DecodeJSON(r, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if err := h.userService.UpdatePassword(id, body.Password); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Password updated successfully"})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	requestingUser, ok := r.Context().Value(middleware.ContextKeyUser).(*models.User)
	if !ok || requestingUser == nil {
		utils.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if err := h.userService.Delete(id, requestingUser.ID); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "User deactivated successfully"})
}

func (h *UserHandler) ListRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.userService.ListRoles()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	utils.RespondJSON(w, http.StatusOK, roles)
}
