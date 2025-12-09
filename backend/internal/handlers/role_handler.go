package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"github.com/wozhdeleniye/redclass-app/internal/services"
)

type RoleHandler struct {
	roleService *services.RoleService
	validate    *validator.Validate
}

func NewRoleHandler(roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		roleService: roleService,
		validate:    validator.New(),
	}
}

func (h *RoleHandler) ChangeRole(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roleIDStr := vars["roleId"]

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var req struct {
		RoleType models.RoleType `json:"role_type" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	role, err := h.roleService.ChangeRole(r.Context(), userID, roleID, req.RoleType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (h *RoleHandler) RemoveFromSubject(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roleIDStr := vars["roleId"]

	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	if err := h.roleService.RemoveFromSubject(r.Context(), userID, roleID); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *RoleHandler) GetSubjectMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	subjectIDStr := vars["id"]

	subjectID, err := uuid.Parse(subjectIDStr)
	if err != nil {
		http.Error(w, "Invalid subject ID", http.StatusBadRequest)
		return
	}

	roles, err := h.roleService.GetSubjectRoles(r.Context(), subjectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}
