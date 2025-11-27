package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"github.com/wozhdeleniye/redclass-app/internal/repositories/postgres"
)

type RoleService struct {
	roleRepo *postgres.RoleRepository
}

func NewRoleService(roleRepo *postgres.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) ChangeRole(ctx context.Context, requesterID, roleID uuid.UUID, newRoleType models.RoleType) (*models.Role, error) {

	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return nil, err
	}

	requesterRole, err := s.roleRepo.GetByUserAndSubject(ctx, requesterID, role.SubjectID)
	if err != nil {
		return nil, err
	}
	if requesterRole == nil || !requesterRole.IsAdmin() {
		return nil, errors.New("only admin can change roles")
	}

	if role.IsAdmin() {
		return nil, errors.New("cannot change admin role")
	}

	role.RoleType = newRoleType
	role.UpdatedAt = time.Now()

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) RemoveFromSubject(ctx context.Context, requesterID, roleID uuid.UUID) error {
	role, err := s.roleRepo.GetByID(ctx, roleID)
	if err != nil {
		return err
	}

	requesterRole, err := s.roleRepo.GetByUserAndSubject(ctx, requesterID, role.SubjectID)
	if err != nil {
		return err
	}
	if requesterRole == nil || !requesterRole.IsAdmin() {
		return errors.New("only admin can remove members")
	}

	if role.IsAdmin() {
		return errors.New("cannot remove admin from subject")
	}

	return s.roleRepo.Delete(ctx, roleID)
}

func (s *RoleService) GetSubjectRoles(ctx context.Context, subjectID uuid.UUID) ([]*models.Role, error) {
	return s.roleRepo.GetSubjectRoles(ctx, subjectID)
}
