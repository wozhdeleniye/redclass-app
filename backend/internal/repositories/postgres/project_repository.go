package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *models.Project) error {
	if project.ID == uuid.Nil {
		project.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).
		Preload("Members.User").
		First(&project, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetByCode(ctx context.Context, code string) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).First(&project, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) GetByTask(ctx context.Context, taskID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.WithContext(ctx).
		Where("task_id = ?", taskID).
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) GetUserProjects(ctx context.Context, userID uuid.UUID) ([]*models.Project, error) {
	var projects []*models.Project
	err := r.db.WithContext(ctx).
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("project_members.user_id = ? AND project_members.deleted_at IS NULL", userID).
		Find(&projects).Error
	return projects, err
}

func (r *ProjectRepository) GetUserProjectByTask(ctx context.Context, userID, taskID uuid.UUID) (*models.Project, error) {
	var project models.Project
	err := r.db.WithContext(ctx).
		Joins("JOIN project_members ON project_members.project_id = projects.id").
		Where("projects.task_id = ? AND project_members.user_id = ? AND project_members.deleted_at IS NULL", taskID, userID).
		First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) AddMember(ctx context.Context, member *models.ProjectMember) error {
	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(member).Error
}

func (r *ProjectRepository) IsUserMember(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ProjectMember{}).
		Where("project_id = ? AND user_id = ?", projectID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetMembersByProject получает участников проекта с предзагрузкой пользователей
func (r *ProjectRepository) GetMembersByProject(ctx context.Context, projectID uuid.UUID) ([]*models.ProjectMember, error) {
	var members []*models.ProjectMember
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("project_id = ? AND deleted_at IS NULL", projectID).
		Find(&members).Error
	return members, err
}
