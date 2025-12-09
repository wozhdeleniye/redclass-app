package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Result struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProblemID uuid.UUID      `json:"problem_id" gorm:"type:uuid;not null;uniqueIndex"`
	CreatorID uuid.UUID      `json:"creator_id" gorm:"type:uuid;not null;index"`
	Done      bool           `json:"done" gorm:"not null;default:true"`
	Comment   string         `json:"comment"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	Problem *Problem `json:"problem" gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE"`
	Creator *User    `json:"creator" gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT"`
}

type CreateResultRequest struct {
	Comment string `json:"comment" validate:"required"`
}

type ResultResponse struct {
	ID        uuid.UUID `json:"id"`
	ProblemID uuid.UUID `json:"problem_id"`
	CreatorID uuid.UUID `json:"creator_id"`
	Done      bool      `json:"done"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}
