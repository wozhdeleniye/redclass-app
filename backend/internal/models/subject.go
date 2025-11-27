package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subject struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"not null;index"`
	Description string         `json:"description"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Roles []*Role `json:"roles" gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE"`
	Tasks []*Task `json:"tasks" gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE"`
}

type CreateSubjectRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Code        string `json:"code" validate:"required"`
}

type UpdateSubjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}
