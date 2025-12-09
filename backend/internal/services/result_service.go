package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/wozhdeleniye/redclass-app/internal/models"
	"github.com/wozhdeleniye/redclass-app/internal/repositories/postgres"
)

type ResultService struct {
	resultRepo  *postgres.ResultRepository
	problemRepo *postgres.ProblemRepository
	projectRepo *postgres.ProjectRepository
}

func NewResultService(rr *postgres.ResultRepository, pr *postgres.ProblemRepository, pjr *postgres.ProjectRepository) *ResultService {
	return &ResultService{
		resultRepo:  rr,
		problemRepo: pr,
		projectRepo: pjr,
	}
}

func (s *ResultService) CreateResult(ctx context.Context, userID uuid.UUID, problemID uuid.UUID, req *models.CreateResultRequest) (*models.Result, error) {

	problem, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return nil, errors.New("problem not found")
	}

	isMember, err := s.projectRepo.IsUserMember(ctx, problem.ProjectID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a project member")
	}

	existing, err := s.resultRepo.GetByProblemID(ctx, problemID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("result for this problem already exists")
	}

	res := &models.Result{
		ProblemID: problemID,
		CreatorID: userID,
		Done:      true,
		Comment:   req.Comment,
	}

	if err := s.resultRepo.Create(ctx, res); err != nil {
		return nil, err
	}

	problem.Solved = true
	if err := s.problemRepo.Update(ctx, problem); err != nil {

		_ = s.resultRepo.Delete(ctx, res.ID)
		return nil, err
	}

	return res, nil
}

func (s *ResultService) GetResult(ctx context.Context, userID uuid.UUID, problemID uuid.UUID) (*models.Result, error) {
	problem, err := s.problemRepo.GetByID(ctx, problemID)
	if err != nil {
		return nil, err
	}
	if problem == nil {
		return nil, errors.New("problem not found")
	}

	isMember, err := s.projectRepo.IsUserMember(ctx, problem.ProjectID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("user is not a project member")
	}

	return s.resultRepo.GetByProblemID(ctx, problemID)
}
