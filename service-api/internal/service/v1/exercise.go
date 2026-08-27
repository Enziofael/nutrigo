package service

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type ExerciseService struct {
	exerciseReposirory repository.ExerciseRepository
	userRepository     repository.UserRepository
}

func NewExerciseService(exerciseR repository.ExerciseRepository, userR repository.UserRepository) *ExerciseService {
	return &ExerciseService{
		exerciseReposirory: exerciseR,
		userRepository:     userR,
	}
}

// Exercises

// Based on [models.ExerciseCreateRequest]
//
// Errors:
//
// - Validation
//   - [models.ErrInvalidUserID]
//   - [models.ErrInvalidName]
//   - [models.ErrInvalidDescription]
//   - [models.ErrInvalidUnitPreferences]
//   - [models.ErrInvalidIncludeQuery]
//
// - Repository
//
// - Database
//   - internal errors (Considered as http.StatusInternalServerError)
//   - constraint violation (Should be pre-checked in service to avoid this)
func (s *ExerciseService) Create(ctx context.Context, req models.ExerciseCreateRequest, include models.IncludeQuery) (*models.ExerciseCreateResponse, error) {

	// Validation

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Checking if user exists (user_id fk constraint)

	existingResp, err := s.exerciseReposirory.Get(ctx, models.ExerciseGetRequest{ID: req.UserID}, models.IncludeQuery{})
	if err == nil && existingResp != nil { // err is nil -> already exists
		return nil, fmt.Errorf("user with id %d %w", req.UserID, repository.ErrAlreadyExists) // AlreadyExists
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) { //err is not ErrNotFound -> something went wrong
		return nil, err // Internal
	}
	// err IS ErrNotFound -> can Create()

	// Creating

	resp, err := s.exerciseReposirory.Create(ctx, req, include)
	if err != nil {
		return nil, err // Internal
	}

	return resp, nil // Created
}

// Based on [models.ExerciseListRequest]
//
// Errors:
//
// - Validation
//   - [models.ErrInvalidUserID] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidLimit] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidOffset] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidSort] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidOrder] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidIncludeQuery] (Considered as http.StatusBadRequest)
//
// - Repository
//
// - Database
//   - internal errors (Considered as http.StatusInternalServerError)
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (s *ExerciseService) List(ctx context.Context, req models.ExerciseListRequest, include models.IncludeQuery) (*models.ExerciseListResponse, error) {

	// Validation

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Listing

	resp, err := s.exerciseReposirory.List(ctx, req, include)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ExerciseService) Get(ctx context.Context, req models.ExerciseGetRequest, include models.IncludeQuery) (*models.ExerciseGetResponse, error) {

	// Validation

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Getting

	resp, err := s.exerciseReposirory.Get(ctx, req, include)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ExerciseService) Patch(ctx context.Context, req models.ExercisePatchRequest, include models.IncludeQuery) (*models.ExercisePatchResponse, error) {

	// Validation

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Getting current

	current, err := s.exerciseReposirory.Get(ctx, models.ExerciseGetRequest{ID: req.ID}, include)
	if err != nil {
		return nil, err
	}

	// Patching

	resp, err := s.exerciseReposirory.Patch(ctx, req, &current.Exercise, current.Version, include)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ExerciseService) Delete(ctx context.Context, req models.ExerciseDeleteRequest) error {

	// Validation

	if err := req.Validate(); err != nil {
		return err
	}

	// Deletion

	err := s.exerciseReposirory.Delete(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

/*

// ExerciseEntries

func (s *ExerciseService) EntryCreate(c *gin.Context, req models.ExerciseEntryCreateRequest) (models.ExerciseEntryCreateResponse, error) {

}

func (s *ExerciseService) EntryList(c *gin.Context, req models.ExerciseEntryListRequest) (models.ExerciseEntryListResponse, error) {

}

/*

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

func (s *ExerciseService) Search(ctx context.Context, req models.ExerciseSearchRequest) (*models.ExerciseSearchResponse, error) {
	req, err := req.Sanitize()
	if err != nil {
		return nil, fmt.Errorf("%w Search at exercise service: %w", repository.ErrInvalidRequest, err)
	}

	return s.repo.Search(ctx, req)
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

*/
