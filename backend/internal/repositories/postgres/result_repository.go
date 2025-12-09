package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"gorm.io/gorm"
)

type ResultRepository struct {
	db *gorm.DB
}

func NewResultRepository(db *gorm.DB) *ResultRepository {
	return &ResultRepository{db: db}
}

func (r *ResultRepository) Create(ctx context.Context, res *models.Result) error {
	if res.ID == uuid.Nil {
		res.ID = uuid.New()
	}
	res.CreatedAt = time.Now()
	return r.db.WithContext(ctx).Create(res).Error
}

func (r *ResultRepository) GetByProblemID(ctx context.Context, problemID uuid.UUID) (*models.Result, error) {
	var res models.Result
	err := r.db.WithContext(ctx).
		Preload("Creator").
		Preload("Problem").
		First(&res, "problem_id = ?", problemID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *ResultRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Result{}, "id = ?", id).Error
}

func (r *ResultRepository) Update(ctx context.Context, res *models.Result) error {
	return r.db.WithContext(ctx).Save(res).Error
}
