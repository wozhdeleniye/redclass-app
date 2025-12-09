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

type ResultHandler struct {
	resultService *services.ResultService
	validate      *validator.Validate
}

func NewResultHandler(rs *services.ResultService) *ResultHandler {
	return &ResultHandler{resultService: rs, validate: validator.New()}
}

// CreateResult (POST /api/problems/{problemId}/result)
func (h *ResultHandler) CreateResult(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	problemIDStr := vars["problemId"]
	problemID, err := uuid.Parse(problemIDStr)
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	var req models.CreateResultRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.resultService.CreateResult(r.Context(), userID, problemID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

// GetResult (GET /api/problems/{problemId}/result)
func (h *ResultHandler) GetResult(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	problemIDStr := vars["problemId"]
	problemID, err := uuid.Parse(problemIDStr)
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	res, err := h.resultService.GetResult(r.Context(), userID, problemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
