package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	if task.ID == uuid.Nil {
		task.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).
		Preload("Subject").
		Preload("CreatedBy").
		First(&task, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetBySubject(ctx context.Context, subjectID uuid.UUID, limit, offset int) ([]*models.Task, int64, error) {
	var tasks []*models.Task
	var total int64

	err := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("subject_id = ?", subjectID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Where("subject_id = ?", subjectID).
		Preload("CreatedBy").
		Limit(limit).
		Offset(offset).
		Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).
		Model(task).
		Updates(task).Error
}

func (r *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, "id = ?", id).Error
}
