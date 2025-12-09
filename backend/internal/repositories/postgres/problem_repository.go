package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"gorm.io/gorm"
)

type ProblemRepository struct {
	db *gorm.DB
}

func NewProblemRepository(db *gorm.DB) *ProblemRepository {
	return &ProblemRepository{db: db}
}

// Create создает новую проблему
func (r *ProblemRepository) Create(ctx context.Context, problem *models.Problem) error {
	return r.db.WithContext(ctx).Create(problem).Error
}

// GetByID получает проблему по ID с предзагруженными данными
func (r *ProblemRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Problem, error) {
	var problem models.Problem
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Project").
		Preload("Parent").
		Preload("Assignees.User").
		First(&problem, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &problem, nil
}

// GetByProjectID получает все проблемы проекта (только главные, без подпроблем)
func (r *ProblemRepository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]*models.Problem, error) {
	var problems []*models.Problem
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Assignees.User").
		Preload("Children").
		Where("project_id = ? AND parent_id IS NULL", projectID).
		Order("number ASC").
		Find(&problems).Error
	return problems, err
}

// GetProjectProblems получает все проблемы проекта включая подпроблемы
func (r *ProblemRepository) GetProjectProblems(ctx context.Context, projectID uuid.UUID) ([]*models.Problem, error) {
	var problems []*models.Problem
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Assignees.User").
		Where("project_id = ?", projectID).
		Order("number ASC").
		Find(&problems).Error
	return problems, err
}

// GetProjectProblemsAssigned получает проблемы проекта, назначенные на конкретного пользователя
func (r *ProblemRepository) GetProjectProblemsAssigned(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) ([]*models.Problem, error) {
	var problems []*models.Problem
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Assignees.User").
		Joins("JOIN problem_assignees ON problem_assignees.problem_id = problems.id").
		Where("problems.project_id = ? AND problem_assignees.user_id = ?", projectID, userID).
		Order("number ASC").
		Find(&problems).Error
	return problems, err
}

// GetChildProblems получает все подпроблемы для родительской проблемы
func (r *ProblemRepository) GetChildProblems(ctx context.Context, parentID uuid.UUID) ([]*models.Problem, error) {
	var problems []*models.Problem
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Assignees.User").
		Where("parent_id = ?", parentID).
		Order("number ASC").
		Find(&problems).Error
	return problems, err
}

// GetNextNumber получает следующий номер проблемы в проекте
func (r *ProblemRepository) GetNextNumber(ctx context.Context, projectID uuid.UUID) (int, error) {
	var maxNumber int
	err := r.db.WithContext(ctx).
		Model(&models.Problem{}).
		Where("project_id = ?", projectID).
		Select("COALESCE(MAX(number), 0)").
		Row().
		Scan(&maxNumber)
	if err != nil {
		return 0, err
	}
	return maxNumber + 1, nil
}

// Update обновляет проблему
func (r *ProblemRepository) Update(ctx context.Context, problem *models.Problem) error {
	return r.db.WithContext(ctx).Save(problem).Error
}

// Delete удаляет проблему
func (r *ProblemRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Problem{}, "id = ?", id).Error
}

// AddAssignee добавляет пользователя к проблеме
func (r *ProblemRepository) AddAssignee(ctx context.Context, problemID uuid.UUID, userID uuid.UUID) (*models.ProblemAssignee, error) {
	// Проверяем, что пользователь уже не назначен
	var existing models.ProblemAssignee
	err := r.db.WithContext(ctx).First(&existing, "problem_id = ? AND user_id = ?", problemID, userID).Error
	if err == nil {
		return &existing, nil // Уже назначен
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	assignee := &models.ProblemAssignee{
		ProblemID: problemID,
		UserID:    userID,
	}
	return assignee, r.db.WithContext(ctx).Create(assignee).Error
}

// RemoveAssignee удаляет пользователя из проблемы
func (r *ProblemRepository) RemoveAssignee(ctx context.Context, problemID uuid.UUID, userID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.ProblemAssignee{}, "problem_id = ? AND user_id = ?", problemID, userID).Error
}

// GetAssignees получает список назначенных пользователей
func (r *ProblemRepository) GetAssignees(ctx context.Context, problemID uuid.UUID) ([]*models.ProblemAssignee, error) {
	var assignees []*models.ProblemAssignee
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("problem_id = ?", problemID).
		Find(&assignees).Error
	return assignees, err
}

// IsUserAssignedToProblem проверяет, назначена ли задача пользователю
func (r *ProblemRepository) IsUserAssignedToProblem(ctx context.Context, problemID uuid.UUID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.ProblemAssignee{}).
		Where("problem_id = ? AND user_id = ?", problemID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetMainProblemByProject получает главную проблему проекта
func (r *ProblemRepository) GetMainProblemByProject(ctx context.Context, projectID uuid.UUID) (*models.Problem, error) {
	var problem models.Problem
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Assignees.User").
		Where("project_id = ? AND parent_id IS NULL", projectID).
		First(&problem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &problem, nil
}
