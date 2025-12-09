package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRole string

const (
	ProjectRoleCreator ProjectRole = "creator"
	ProjectRoleMember  ProjectRole = "member"
)

type Project struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TaskID      uuid.UUID      `json:"task_id" gorm:"type:uuid;not null;index"`
	CreatorID   uuid.UUID      `json:"creator_id" gorm:"type:uuid;not null;index"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Code        string         `json:"code" gorm:"uniqueIndex;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Task    *Task            `json:"task" gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
	Creator *User            `json:"creator" gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT"`
	Members []*ProjectMember `json:"members" gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

type ProjectMember struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProjectID uuid.UUID      `json:"project_id" gorm:"type:uuid;not null;index:idx_project_user,unique"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index:idx_project_user,unique"`
	Role      ProjectRole    `json:"role" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User    *User    `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Project *Project `json:"project" gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

type CreateProjectRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type JoinProjectRequest struct {
	Code string `json:"code" validate:"required"`
}
