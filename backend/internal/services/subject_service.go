package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"github.com/wozhdeleniye/redclass-app/internal/repositories/postgres"
)

type SubjectService struct {
	subjectRepo *postgres.SubjectRepository
	roleRepo    *postgres.RoleRepository
	userRepo    *postgres.UserRepository
}

func NewSubjectService(
	subjectRepo *postgres.SubjectRepository,
	roleRepo *postgres.RoleRepository,
	userRepo *postgres.UserRepository,
) *SubjectService {
	return &SubjectService{
		subjectRepo: subjectRepo,
		roleRepo:    roleRepo,
		userRepo:    userRepo,
	}
}

func (s *SubjectService) CreateSubject(ctx context.Context, creatorID uuid.UUID, req *models.CreateSubjectRequest) (*models.Subject, error) {
	if req.Name == "" {
		return nil, errors.New("subject name is required")
	}
	if req.Code == "" {
		return nil, errors.New("subject code is required")
	}

	subject := &models.Subject{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.subjectRepo.Create(ctx, subject); err != nil {
		return nil, fmt.Errorf("failed to create subject: %w", err)
	}

	adminRole := &models.Role{
		ID:        uuid.New(),
		UserID:    creatorID,
		SubjectID: subject.ID,
		RoleType:  models.RoleAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roleRepo.Create(ctx, adminRole); err != nil {
		return nil, fmt.Errorf("failed to create admin role: %w", err)
	}

	return subject, nil
}

func (s *SubjectService) GetSubjectByID(ctx context.Context, id uuid.UUID) (*models.Subject, error) {
	return s.subjectRepo.GetByID(ctx, id)
}

func (s *SubjectService) GetAllSubjects(ctx context.Context, limit, offset int) ([]*models.Subject, int64, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.subjectRepo.GetAll(ctx, limit, offset)
}

func (s *SubjectService) UpdateSubject(ctx context.Context, id uuid.UUID, req *models.UpdateSubjectRequest) (*models.Subject, error) {
	subject, err := s.subjectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		subject.Name = *req.Name
	}
	if req.Description != nil {
		subject.Description = *req.Description
	}
	subject.UpdatedAt = time.Now()

	if err := s.subjectRepo.Update(ctx, subject); err != nil {
		return nil, err
	}

	return subject, nil
}

func (s *SubjectService) DeleteSubject(ctx context.Context, id uuid.UUID) error {
	return s.subjectRepo.Delete(ctx, id)
}

func (s *SubjectService) JoinSubject(ctx context.Context, userID uuid.UUID, code string) (*models.Role, error) {

	subject, err := s.subjectRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, errors.New("subject with this code not found")
	}

	existingRole, err := s.roleRepo.GetByUserAndSubject(ctx, userID, subject.ID)
	if err != nil {
		return nil, err
	}
	if existingRole != nil {
		return nil, errors.New("user is already a member of this subject")
	}

	role := &models.Role{
		ID:        uuid.New(),
		UserID:    userID,
		SubjectID: subject.ID,
		RoleType:  models.RoleStudent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	role.Subject = subject
	return role, nil
}

func (s *SubjectService) GetUserSubjects(ctx context.Context, userID uuid.UUID) ([]*models.Subject, error) {
	return s.subjectRepo.GetUserSubjects(ctx, userID)
}
