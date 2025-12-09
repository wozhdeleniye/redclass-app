package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleType string

const (
	RoleStudent RoleType = "student"
	RoleTeacher RoleType = "teacher"
	RoleAdmin   RoleType = "admin"
)

type Role struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index:idx_user_subject,unique"`
	SubjectID uuid.UUID      `json:"subject_id" gorm:"type:uuid;not null;index:idx_user_subject,unique"`
	RoleType  RoleType       `json:"role_type" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	User    *User    `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Subject *Subject `json:"subject" gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE"`
}

func (r *Role) TableName() string {
	return "roles"
}

func (r *Role) IsTeacher() bool {
	return r.RoleType == RoleTeacher
}

func (r *Role) IsStudent() bool {
	return r.RoleType == RoleStudent
}

func (r *Role) IsAdmin() bool {
	return r.RoleType == RoleAdmin
}

type CreateRoleRequest struct {
	UserID      uuid.UUID `json:"user_id" validate:"required"`
	SubjectID   uuid.UUID `json:"subject_id" validate:"required"`
	RoleType    RoleType  `json:"role_type" validate:"required"`
	Permissions []string  `json:"permissions"`
}

type UpdateRoleRequest struct {
	RoleType    *RoleType `json:"role_type"`
	Permissions *[]string `json:"permissions"`
}
