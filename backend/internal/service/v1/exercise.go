package service

import (
	"context"
	"strings"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type ExerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(repo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

func (s *ExerciseService) GetByID(ctx context.Context, id int64) (*models.Exercise, error) {
	if id == 0 {
		return nil, ErrIdRequired
	}
	return s.repo.GetByID(ctx, id)
}

func (s *ExerciseService) ListByTgID(ctx context.Context, req models.ExerciseListRequest) (*[]models.Exercise, error) {
	return s.repo.ListByTgID(ctx, req)
}

func (s *ExerciseService) CountByTgID(ctx context.Context, tgID int64) (int, error) {
	return s.repo.CountByTgID(ctx, tgID)
}

func (s *ExerciseService) Create(ctx context.Context, req models.ExerciseCreateRequest) (*models.Exercise, error) {
	if req.TgID == 0 {
		return nil, ErrTgIdRequired
	}
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return nil, ErrNameRequired
	}

	return s.repo.Create(ctx, req)
}

func (s *ExerciseService) Patch(ctx context.Context, id int64, req models.ExercisePatchRequest) (*models.Exercise, error) {
	if id == 0 {
		return nil, ErrIdRequired
	}

	if req.Name == nil &&
		req.Description == nil &&
		req.Technique == nil &&
		req.WeightUnit == nil {
		return nil, ErrInvalidRequest
	}

	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		if trimmed == "" {
			return nil, ErrNameRequired
		}
		req.Name = &trimmed
	}

	if req.WeightUnit != nil {
		trimmed := strings.TrimSpace(*req.WeightUnit)
		if trimmed == "" {
			return nil, ErrWeightUnitRequired
		}
		req.WeightUnit = &trimmed
	}

	return s.repo.Patch(ctx, id, req)
}

func (s *ExerciseService) Delete(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrIdRequired
	}
	return s.repo.Delete(ctx, id)
}
