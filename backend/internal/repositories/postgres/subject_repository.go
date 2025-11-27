package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) Create(ctx context.Context, subject *models.Subject) error {
	if subject.ID == uuid.Nil {
		subject.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(subject).Error
}

func (r *SubjectRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.WithContext(ctx).
		Preload("Roles").
		Preload("Tasks").
		First(&subject, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, err
	}
	return &subject, nil
}

func (r *SubjectRepository) GetByCode(ctx context.Context, code string) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.WithContext(ctx).
		First(&subject, "code = ?", code).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, err
	}
	return &subject, nil
}

func (r *SubjectRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Subject, int64, error) {
	var subjects []*models.Subject
	var total int64

	err := r.db.WithContext(ctx).Model(&models.Subject{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&subjects).Error
	if err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}

func (r *SubjectRepository) Update(ctx context.Context, subject *models.Subject) error {
	return r.db.WithContext(ctx).
		Model(subject).
		Updates(subject).Error
}

func (r *SubjectRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Subject{}, "id = ?", id).Error
}

func (r *SubjectRepository) GetUserSubjects(ctx context.Context, userID uuid.UUID) ([]*models.Subject, error) {
	var subjects []*models.Subject
	err := r.db.WithContext(ctx).
		Joins("JOIN roles ON roles.subject_id = subjects.id").
		Where("roles.user_id = ? AND roles.deleted_at IS NULL", userID).
		Find(&subjects).Error
	return subjects, err
}
