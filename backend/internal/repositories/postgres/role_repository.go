package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *RoleRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Subject").
		First(&role, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetByUserAndSubject(ctx context.Context, userID, subjectID uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Subject").
		First(&role, "user_id = ? AND subject_id = ?", userID, subjectID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetUserRoles(ctx context.Context, userID uuid.UUID) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.WithContext(ctx).
		Preload("Subject").
		Where("user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) GetSubjectRoles(ctx context.Context, subjectID uuid.UUID) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("subject_id = ?", subjectID).
		Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) Update(ctx context.Context, role *models.Role) error {
	return r.db.WithContext(ctx).
		Model(role).
		Updates(role).Error
}

func (r *RoleRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Role{}, "id = ?", id).Error
}

func (r *RoleRepository) GetSubjectAdmin(ctx context.Context, subjectID uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.WithContext(ctx).
		Preload("User").
		Where("subject_id = ? AND role_type = ?", subjectID, models.RoleAdmin).
		First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}
