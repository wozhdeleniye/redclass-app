package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"github.com/wozhdeleniye/redclass-app/internal/repositories/postgres"
)

type ProjectService struct {
	projectRepo *postgres.ProjectRepository
	taskRepo    *postgres.TaskRepository
	roleRepo    *postgres.RoleRepository
	problemRepo *postgres.ProblemRepository
}

func NewProjectService(pr *postgres.ProjectRepository, tr *postgres.TaskRepository, rr *postgres.RoleRepository, prr *postgres.ProblemRepository) *ProjectService {
	return &ProjectService{
		projectRepo: pr,
		taskRepo:    tr,
		roleRepo:    rr,
		problemRepo: prr,
	}
}

func randomCode(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b)[:n], nil
}

func (s *ProjectService) CreateProject(ctx context.Context, userID, taskID uuid.UUID, req *models.CreateProjectRequest) (*models.Project, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	role, err := s.roleRepo.GetByUserAndSubject(ctx, userID, task.SubjectID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("user must be a member of subject to create project")
	}

	existingProject, err := s.projectRepo.GetUserProjectByTask(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}
	if existingProject != nil {
		return nil, errors.New("user can only be in one project per task")
	}

	var code string
	for i := 0; i < 5; i++ {
		c, err := randomCode(8)
		if err != nil {
			return nil, err
		}
		existing, err := s.projectRepo.GetByCode(ctx, c)
		if err != nil {
			return nil, err
		}
		if existing == nil {
			code = c
			break
		}
	}
	if code == "" {
		return nil, errors.New("failed to generate unique project code")
	}

	project := &models.Project{
		ID:          uuid.New(),
		TaskID:      taskID,
		CreatorID:   userID,
		Title:       req.Title,
		Description: req.Description,
		Code:        code,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	member := &models.ProjectMember{
		ID:        uuid.New(),
		ProjectID: project.ID,
		UserID:    userID,
		Role:      models.ProjectRoleCreator,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.projectRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}

	task, err = s.taskRepo.GetByID(ctx, project.TaskID)
	if err != nil {
		return nil, err
	}
	project.Task = task

	problemService := NewProblemService(s.problemRepo, s.projectRepo)
	if _, err := problemService.CreateMainProblem(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) JoinProject(ctx context.Context, userID uuid.UUID, code string) (*models.ProjectMember, error) {
	project, err := s.projectRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("project with this code not found")
	}

	isMember, err := s.projectRepo.IsUserMember(ctx, project.ID, userID)
	if err != nil {
		return nil, err
	}
	if isMember {
		return nil, errors.New("user is already a member of this project")
	}

	member := &models.ProjectMember{
		ID:        uuid.New(),
		ProjectID: project.ID,
		UserID:    userID,
		Role:      models.ProjectRoleMember,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.projectRepo.AddMember(ctx, member); err != nil {
		return nil, err
	}
	member.Project = project
	return member, nil
}

func (s *ProjectService) GetTaskProjects(ctx context.Context, taskID uuid.UUID) ([]*models.Project, error) {
	return s.projectRepo.GetByTask(ctx, taskID)
}

func (s *ProjectService) GetUserProjects(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	return s.projectRepo.GetUserProjects(ctx, userID)
}

func (s *ProjectService) GetProjectUsers(ctx context.Context, projectID uuid.UUID) ([]*models.User, error) {
	members, err := s.projectRepo.GetMembersByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	users := make([]*models.User, 0, len(members))
	for _, m := range members {
		if m.User != nil {
			users = append(users, m.User)
		}
	}
	return users, nil
}
