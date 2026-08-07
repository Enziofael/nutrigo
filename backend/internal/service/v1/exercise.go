package service

import (
	"context"
	"fmt"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type ExerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(repo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

// =========================================================
// CRUD REQUESTS
// =========================================================

func (s *ExerciseService) Create(ctx context.Context, req models.ExerciseCreateRequest) (*models.Exercise, error) {
	req, err := req.Sanitize()
	if err != nil {
		return nil, fmt.Errorf("%w Create at exercise service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.Create(ctx, req)
}

func (s *ExerciseService) Delete(ctx context.Context, req models.ExerciseDeleteRequest) error {
	req, err := req.Sanitize()
	if err != nil {
		return fmt.Errorf("%w Delete at exercise service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.Delete(ctx, req)
}

func (s *ExerciseService) Get(ctx context.Context, req models.ExerciseGetRequest) (*models.Exercise, error) {
	req, err := req.Sanitize()
	if err != nil {
		return nil, fmt.Errorf("%w Get at exercise service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.Get(ctx, req)
}

func (s *ExerciseService) List(ctx context.Context, req models.ExerciseListRequest) (*[]models.Exercise, error) {
	req, err := req.Sanitize()
	if err != nil {
		return nil, fmt.Errorf("%w List at exercise service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.List(ctx, req)
}

func (s *ExerciseService) Count(ctx context.Context, req models.ExerciseCountRequest) (int, error) {
	req, err := req.Sanitize()
	if err != nil {
		return 0, fmt.Errorf("%w Count at exercise service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.Count(ctx, req)
}

func (s *ExerciseService) Patch(ctx context.Context, req models.ExercisePatchRequest) (*models.Exercise, error) {
	req, err := req.Sanitize()
	if err != nil {
		return nil, fmt.Errorf("%w Patch at exercise service: %w", repository.ErrInvalidRequest, err)
	}

	existing, err := s.repo.Get(ctx, models.ExerciseGetRequest{ID: req.ID})
	if err != nil {
		return nil, err
	}

	return s.repo.Patch(ctx, req, existing)
}
