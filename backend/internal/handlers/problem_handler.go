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

type ProblemHandler struct {
	problemService *services.ProblemService
	resultService  *services.ResultService
	validate       *validator.Validate
}

func NewProblemHandler(ps *services.ProblemService, rs *services.ResultService) *ProblemHandler {
	return &ProblemHandler{problemService: ps, resultService: rs, validate: validator.New()}
}

// CreateProblem создает проблему в проекте (POST /api/projects/{projectId}/problems)
func (h *ProblemHandler) CreateProblem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	projectIDStr := vars["projectId"]
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var req models.CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	problem, err := h.problemService.CreateProblem(r.Context(), userID, projectID, nil, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(problem)
}

// CreateSubproblem создает подпроблему (POST /api/problems/{parentId}/subproblems)
func (h *ProblemHandler) CreateSubproblem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	parentIDStr := vars["parentId"]
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}

	var req models.CreateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем родительскую проблему чтобы узнать projectID
	parentProblem, err := h.problemService.GetProblemByIDDirect(r.Context(), parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	problem, err := h.problemService.CreateProblem(r.Context(), userID, parentProblem.ProjectID, &parentID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(problem)
}

// UpdateProblem обновляет проблему (PUT /api/problems/{problemId})
func (h *ProblemHandler) UpdateProblem(w http.ResponseWriter, r *http.Request) {
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

	var req models.UpdateProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	problem, err := h.problemService.UpdateProblem(r.Context(), userID, problemID, &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problem)
}

// DeleteProblem удаляет проблему (DELETE /api/problems/{problemId})
func (h *ProblemHandler) DeleteProblem(w http.ResponseWriter, r *http.Request) {
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

	if err := h.problemService.DeleteProblem(r.Context(), userID, problemID); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetProblem получает проблему по ID (GET /api/problems/{problemId})
func (h *ProblemHandler) GetProblem(w http.ResponseWriter, r *http.Request) {
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

	problem, err := h.problemService.GetProblem(r.Context(), userID, problemID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var result *models.Result
	if problem.Solved {
		if h.resultService != nil {
			res, err := h.resultService.GetResult(r.Context(), userID, problemID)
			if err == nil {
				result = res
			}
		}
	}

	childStats, err := h.problemService.GetChildrenStatistics(r.Context(), problemID)
	if err != nil {
		childStats = &models.ChildrenStatistics{}
	}

	resp := struct {
		Problem            *models.Problem            `json:"problem"`
		Result             *models.Result             `json:"result,omitempty"`
		ChildrenStatistics *models.ChildrenStatistics `json:"children_statistics,omitempty"`
	}{
		Problem:            problem,
		Result:             result,
		ChildrenStatistics: childStats,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GetProjectProblems получает все проблемы проекта (GET /api/projects/{projectId}/problems)
func (h *ProblemHandler) GetProjectProblems(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	projectIDStr := vars["projectId"]
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}
	q := r.URL.Query()
	assignedOnly := false
	if v := q.Get("assigned_only"); v == "1" || v == "true" || v == "True" {
		assignedOnly = true
	}

	problems, err := h.problemService.GetProjectProblems(r.Context(), userID, projectID, assignedOnly)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problems)
}

// GetMainProblem получает главную проблему проекта (GET /api/projects/{projectId}/problems/main)
func (h *ProblemHandler) GetMainProblem(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	projectIDStr := vars["projectId"]
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	problem, err := h.problemService.GetMainProblem(r.Context(), userID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problem)
}

// GetSubproblems получает все дочерние проблемы (GET /api/problems/{parentId}/subproblems)
func (h *ProblemHandler) GetSubproblems(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	parentIDStr := vars["parentId"]
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}

	problems, err := h.problemService.GetSubproblems(r.Context(), userID, parentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(problems)
}

// GetProjectStatistics получает статистику по всем проблемам проекта (GET /api/projects/{projectId}/statistics)
func (h *ProblemHandler) GetProjectStatistics(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(uuid.UUID)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	projectIDStr := vars["projectId"]
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	stats, err := h.problemService.GetProjectStatistics(r.Context(), userID, projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
