package service

import (
	"context"
	"fmt"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// =========================================================
// CRUD REQUESTS
// =========================================================

func (s *UserService) Create(ctx context.Context, req models.UserCreateRequest) (*models.User, error) {
	req, err := req.Sanitize()
	if err != nil {
		return nil, fmt.Errorf("%w Create at user service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.Create(ctx, req)
}

func (s *UserService) Delete(ctx context.Context, req models.UserDeleteRequest) error {
	req, err := req.Sanitize()
	if err != nil {
		return fmt.Errorf("%w Delete at user service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.Delete(ctx, req)
}

func (s *UserService) Get(ctx context.Context, req models.UserGetRequest) (*models.User, error) {
	req, err := req.Sanitize()
	if err != nil {
		return nil, fmt.Errorf("%w Get at user service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.Get(ctx, req)
}

func (s *UserService) Patch(ctx context.Context, req models.UserPatchRequest) (*models.User, error) {
	req, err := req.Sanitize()
	if err != nil {
		return nil, fmt.Errorf("%w Patch at user service: %w", repository.ErrInvalidRequest, err)
	}

	existing, err := s.repo.Get(ctx, models.UserGetRequest{TgID: req.TgID})
	if err != nil {
		return nil, err
	}

	return s.repo.Patch(ctx, req, existing)
}
