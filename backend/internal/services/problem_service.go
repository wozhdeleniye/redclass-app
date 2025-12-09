package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"github.com/wozhdeleniye/redclass-app/internal/repositories/postgres"
)

type ProblemService struct {
	problemRepo *postgres.ProblemRepository
	projectRepo *postgres.ProjectRepository
}

func NewProblemService(
	problemRepo *postgres.ProblemRepository,
	projectRepo *postgres.ProjectRepository,
) *ProblemService {
	return &ProblemService{
		problemRepo: problemRepo,
		projectRepo: projectRepo,
	}
}

// CreateMainProblem создает главную проблему при создании проекта
func (s *ProblemService) CreateMainProblem(ctx context.Context, project *models.Project) (*models.Problem, error) {
	if project.Task == nil {
		return nil, errors.New("project task is not loaded")
	}

	now := time.Now()
	endTime := now
	if project.Task.DueDate != nil {
		endTime = *project.Task.DueDate
	}

	mainProblem := &models.Problem{
		ProjectID:   project.ID,
		ParentID:    nil, // главная проблема
		CreatorID:   project.CreatorID,
		Number:      1,
		Title:       fmt.Sprintf("Main: %s", project.Title),
		Description: fmt.Sprintf("Main problem for project: %s", project.Title),
		StartTime:   now,
		EndTime:     endTime,
		Solved:      false,
	}

	if err := s.problemRepo.Create(ctx, mainProblem); err != nil {
		return nil, err
	}
	return mainProblem, nil
}

// CreateProblem создает новую проблему (подпроблему или главную)
func (s *ProblemService) CreateProblem(
	ctx context.Context,
	userID uuid.UUID,
	projectID uuid.UUID,
	parentID *uuid.UUID,
	req *models.CreateProblemRequest,
) (*models.Problem, error) {
	// Проверяем, что пользователь член проекта
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project not found")
	}

	isMember, err := s.projectRepo.IsUserMember(ctx, projectID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a project member")
	}

	// Получаем главную проблему для валидации времени
	mainProblem, err := s.problemRepo.GetMainProblemByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if mainProblem == nil {
		return nil, errors.New("main problem not found")
	}

	// Устанавливаем время по умолчанию
	startTime := mainProblem.StartTime
	endTime := mainProblem.EndTime

	if req.StartTime != nil {
		startTime = *req.StartTime
	}
	if req.EndTime != nil {
		endTime = *req.EndTime
	}

	// Валидируем время
	if startTime.Before(mainProblem.StartTime) {
		return nil, errors.New("start_time cannot be earlier than main problem start_time")
	}
	if endTime.After(mainProblem.EndTime) {
		return nil, errors.New("end_time cannot be later than main problem end_time")
	}
	if startTime.After(endTime) {
		return nil, errors.New("start_time cannot be after end_time")
	}

	// Получаем следующий номер
	number, err := s.problemRepo.GetNextNumber(ctx, projectID)
	if err != nil {
		return nil, err
	}

	problem := &models.Problem{
		ProjectID:   projectID,
		ParentID:    parentID,
		CreatorID:   userID,
		Number:      number,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		Solved:      false,
	}

	if err := s.problemRepo.Create(ctx, problem); err != nil {
		return nil, err
	}

	// Добавляем назначенных пользователей
	if len(req.AssigneeIDs) > 0 {
		for _, assigneeID := range req.AssigneeIDs {
			// Проверяем, что назначаемый пользователь - член проекта
			isMember, err := s.projectRepo.IsUserMember(ctx, projectID, assigneeID)
			if err != nil {
				return nil, err
			}
			if !isMember {
				return nil, fmt.Errorf("user %s is not a project member", assigneeID)
			}

			if _, err := s.problemRepo.AddAssignee(ctx, problem.ID, assigneeID); err != nil {
				return nil, err
			}
		}
	}

	return problem, nil
}

// UpdateProblem обновляет проблему (только создатель или главный проекта)
func (s *ProblemService) UpdateProblem(
	ctx context.Context,
	userID uuid.UUID,
	problemID uuid.UUID,
	req *models.UpdateProblemRequest,
) (*models.Problem, error) {
	problem, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return nil, errors.New("problem not found")
	}

	// Получаем главную проблему для валидации времени
	mainProblem, err := s.problemRepo.GetMainProblemByProject(ctx, problem.ProjectID)
	if err != nil {
		return nil, err
	}

	// Обновляем поля
	if req.Title != nil {
		problem.Title = *req.Title
	}
	if req.Description != nil {
		problem.Description = *req.Description
	}

	// Обновляем время с валидацией
	if req.StartTime != nil {
		if req.StartTime.Before(mainProblem.StartTime) {
			return nil, errors.New("start_time cannot be earlier than main problem start_time")
		}
		problem.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		if req.EndTime.After(mainProblem.EndTime) {
			return nil, errors.New("end_time cannot be later than main problem end_time")
		}
		problem.EndTime = *req.EndTime
	}

	if problem.StartTime.After(problem.EndTime) {
		return nil, errors.New("start_time cannot be after end_time")
	}

	// Обновляем назначенных
	if req.AssigneeIDs != nil {
		// Удаляем старых
		existingAssignees, err := s.problemRepo.GetAssignees(ctx, problemID)
		if err != nil {
			return nil, err
		}

		for _, existing := range existingAssignees {
			// Проверяем, находится ли в новом списке
			found := false
			for _, newID := range *req.AssigneeIDs {
				if existing.UserID == newID {
					found = true
					break
				}
			}
			if !found {
				if err := s.problemRepo.RemoveAssignee(ctx, problemID, existing.UserID); err != nil {
					return nil, err
				}
			}
		}

		// Добавляем новых
		for _, assigneeID := range *req.AssigneeIDs {
			if _, err := s.problemRepo.AddAssignee(ctx, problemID, assigneeID); err != nil {
				return nil, err
			}
		}
	}

	if err := s.problemRepo.Update(ctx, problem); err != nil {
		return nil, err
	}

	return problem, nil
}

// DeleteProblem удаляет проблему
func (s *ProblemService) DeleteProblem(ctx context.Context, userID uuid.UUID, problemID uuid.UUID) error {
	problem, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		return err
	}
	if problem == nil {
		return errors.New("problem not found")
	}

	// Проверяем права: создатель проблемы или создатель проекта
	project, err := s.projectRepo.GetByID(ctx, problem.ProjectID)
	if err != nil {
		return err
	}

	if problem.CreatorID != userID && project.CreatorID != userID {
		return errors.New("only creator or project owner can delete problem")
	}

	return s.problemRepo.Delete(ctx, problemID)
}

// GetProblem получает проблему
func (s *ProblemService) GetProblem(ctx context.Context, userID uuid.UUID, problemID uuid.UUID) (*models.Problem, error) {
	problem, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return nil, errors.New("problem not found")
	}

	return problem, nil
}

// GetProjectProblems получает все проблемы проекта
func (s *ProblemService) GetProjectProblems(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, assignedOnly bool) ([]*models.Problem, error) {

	if assignedOnly {
		return s.problemRepo.GetProjectProblemsAssigned(ctx, projectID, userID)
	}

	return s.problemRepo.GetProjectProblems(ctx, projectID)
}

// GetMainProblem получает главную проблему проекта
func (s *ProblemService) GetMainProblem(ctx context.Context, userID uuid.UUID, projectID uuid.UUID) (*models.Problem, error) {
	return s.problemRepo.GetMainProblemByProject(ctx, projectID)
}

// GetProblemByIDDirect получает проблему по ID без проверки доступа (внутренний метод)
func (s *ProblemService) GetProblemByIDDirect(ctx context.Context, problemID uuid.UUID) (*models.Problem, error) {
	return s.problemRepo.GetByID(ctx, problemID)
}

// GetSubproblems получает все дочерние проблемы для родительской проблемы
func (s *ProblemService) GetSubproblems(ctx context.Context, userID uuid.UUID, parentID uuid.UUID) ([]*models.Problem, error) {
	// Получаем родительскую проблему
	parent, err := s.problemRepo.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent problem not found")
	}

	return s.problemRepo.GetChildProblems(ctx, parentID)
}

// GetProjectStatistics возвращает статистику по всем проблемам проекта
func (s *ProblemService) GetProjectStatistics(ctx context.Context, userID uuid.UUID, projectID uuid.UUID) (*models.ProblemStatistics, error) {
	return s.problemRepo.GetProjectStatistics(ctx, projectID)
}

// GetChildrenStatistics возвращает статистику по дочерним проблемам
func (s *ProblemService) GetChildrenStatistics(ctx context.Context, parentID uuid.UUID) (*models.ChildrenStatistics, error) {
	return s.problemRepo.GetChildrenStatistics(ctx, parentID)
}
