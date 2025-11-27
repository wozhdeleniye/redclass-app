package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"github.com/wozhdeleniye/redclass-app/internal/repositories/postgres"
)

type TaskService struct {
	taskRepo    *postgres.TaskRepository
	roleRepo    *postgres.RoleRepository
	subjectRepo *postgres.SubjectRepository
}

func NewTaskService(
	taskRepo *postgres.TaskRepository,
	roleRepo *postgres.RoleRepository,
	subjectRepo *postgres.SubjectRepository,
) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		roleRepo:    roleRepo,
		subjectRepo: subjectRepo,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, userID uuid.UUID, req *models.CreateTaskRequest) (*models.Task, error) {

	userRole, err := s.roleRepo.GetByUserAndSubject(ctx, userID, req.SubjectID)
	if err != nil {
		return nil, err
	}
	if userRole == nil || (!userRole.IsTeacher() && !userRole.IsAdmin()) {
		return nil, errors.New("only teacher or admin can create tasks")
	}

	if req.Title == "" {
		return nil, errors.New("task title is required")
	}

	task := &models.Task{
		ID:          uuid.New(),
		SubjectID:   req.SubjectID,
		CreatedByID: userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) GetTaskByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	return s.taskRepo.GetByID(ctx, id)
}

func (s *TaskService) GetTasksBySubject(ctx context.Context, subjectID uuid.UUID, limit, offset int) ([]*models.Task, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.taskRepo.GetBySubject(ctx, subjectID, limit, offset)
}

func (s *TaskService) UpdateTask(ctx context.Context, userID, taskID uuid.UUID, req *models.UpdateTaskRequest) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	userRole, err := s.roleRepo.GetByUserAndSubject(ctx, userID, task.SubjectID)
	if err != nil {
		return nil, err
	}
	if userRole == nil || (task.CreatedByID != userID && !userRole.IsAdmin()) {
		return nil, errors.New("only task creator or admin can update task")
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *TaskService) DeleteTask(ctx context.Context, userID, taskID uuid.UUID) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return err
	}

	userRole, err := s.roleRepo.GetByUserAndSubject(ctx, userID, task.SubjectID)
	if err != nil {
		return err
	}
	if userRole == nil || (task.CreatedByID != userID && !userRole.IsAdmin()) {
		return errors.New("only task creator or admin can delete task")
	}

	return s.taskRepo.Delete(ctx, taskID)
}
