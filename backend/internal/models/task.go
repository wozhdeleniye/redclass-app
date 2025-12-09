package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SubjectID   uuid.UUID      `json:"subject_id" gorm:"type:uuid;not null;index"`
	CreatedByID uuid.UUID      `json:"created_by_id" gorm:"type:uuid;not null;index"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Subject   *Subject `json:"subject" gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE"`
	CreatedBy *User    `json:"created_by" gorm:"foreignKey:CreatedByID;constraint:OnDelete:RESTRICT"`
}

type CreateTaskRequest struct {
	SubjectID   uuid.UUID  `json:"subject_id" validate:"required"`
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description"`
	DueDate     *time.Time `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
}
