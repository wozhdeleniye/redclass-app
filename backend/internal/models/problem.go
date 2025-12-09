package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Problem struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProjectID   uuid.UUID      `json:"project_id" gorm:"type:uuid;not null;index"`
	ParentID    *uuid.UUID     `json:"parent_id" gorm:"type:uuid;index"`
	CreatorID   uuid.UUID      `json:"creator_id" gorm:"type:uuid;not null;index"`
	Number      int            `json:"number" gorm:"not null"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	StartTime   time.Time      `json:"start_time" gorm:"not null"`
	EndTime     time.Time      `json:"end_time" gorm:"not null"`
	Solved      bool           `json:"solved" gorm:"not null;default:false"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Project   *Project           `json:"project" gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	Parent    *Problem           `json:"parent" gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL"`
	Creator   *User              `json:"creator" gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT"`
	Assignees []*ProblemAssignee `json:"assignees" gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE"`
	Children  []*Problem         `json:"children" gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE"`
}

type ProblemAssignee struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProblemID uuid.UUID      `json:"problem_id" gorm:"type:uuid;not null;index:idx_problem_assignee,unique"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index:idx_problem_assignee,unique"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User    *User    `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Problem *Problem `json:"problem" gorm:"foreignKey:ProblemID;constraint:OnDelete:CASCADE"`
}

type CreateProblemRequest struct {
	Title       string      `json:"title" validate:"required"`
	Description string      `json:"description"`
	StartTime   *time.Time  `json:"start_time"`
	EndTime     *time.Time  `json:"end_time"`
	AssigneeIDs []uuid.UUID `json:"assignee_ids"`
}

type UpdateProblemRequest struct {
	Title       *string      `json:"title"`
	Description *string      `json:"description"`
	StartTime   *time.Time   `json:"start_time"`
	EndTime     *time.Time   `json:"end_time"`
	AssigneeIDs *[]uuid.UUID `json:"assignee_ids"`
}

type ProblemStatistics struct {
	Completed  int `json:"completed"`
	Incomplete int `json:"incomplete"`
	Total      int `json:"total"`
	Percentage int `json:"percentage"`
}

type ChildrenStatistics struct {
	Completed  int `json:"completed"`
	Incomplete int `json:"incomplete"`
	Total      int `json:"total"`
}
