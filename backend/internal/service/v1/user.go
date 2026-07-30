package service

import (
	"context"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetByTgID(ctx context.Context, tgID int64) (*models.User, error) {
	if tgID == 0 {
		return nil, ErrTgIdRequired
	}
	return s.repo.GetByTgId(ctx, tgID)
}

func (s *UserService) Create(ctx context.Context, req models.UserCreateRequest) (*models.User, error) {
	if req.TgID == 0 {
		return nil, ErrTgIdRequired
	}
	if req.TgTag == "" {
		return nil, ErrTgTagRequired
	}

	// Проверка существования
	existing, err := s.repo.GetByTgId(ctx, req.TgID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, repository.ErrUserAlreadyExists
	}

	return s.repo.Create(ctx, req)
}

func (s *UserService) Patch(ctx context.Context, tgID int64, req models.UserPatchRequest) (*models.User, error) {
	if tgID == 0 {
		return nil, ErrTgIdRequired
	}

	if req.TgTag == nil &&
		req.LastMessagedAt == nil &&
		req.Status == nil &&
		req.Context == nil &&
		req.ContextData == nil {
		return nil, ErrInvalidRequest
	}

	if req.Status != nil && !models.ValidateStatus(*req.Status) {
		return nil, ErrInvalidStatus
	}

	return s.repo.Patch(ctx, tgID, req)
}

func (s *UserService) Delete(ctx context.Context, tgID int64) error {
	if tgID == 0 {
		return ErrTgIdRequired
	}

	existing, err := s.repo.GetByTgId(ctx, tgID)
	if err != nil {
		return err
	}
	if existing == nil {
		return repository.ErrUserNotFound
	}

	return s.repo.Delete(ctx, tgID)
}
