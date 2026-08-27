package service

import (
	"context"
	"errors"
	"fmt"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type UserService struct {
	userRepository     repository.UserRepository
	exerciseRepository repository.ExerciseRepository
}

func NewUserService(userR repository.UserRepository, exerciseR repository.ExerciseRepository) *UserService {
	return &UserService{
		userRepository:     userR,
		exerciseRepository: exerciseR,
	}
}

// Users

// Based on [models.UserCreateRequest]
//
// Errors:
//
// - Validation
//   - [models.ErrInvalidIncludeQuery] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidID] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidTgTag] (Considered as http.StatusBadRequest)
//
// - Service
//   - [repository.ErrAlreadyExists] (Considered as http.StatusConflict)
//
// - Repository
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this)
func (s *UserService) Create(ctx context.Context, req models.UserCreateRequest, include models.IncludeQuery) (*models.UserCreateResponse, error) {

	// Validation

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Checking if exists (ID primary key unique constraint)

	existingResp, err := s.userRepository.Get(ctx, models.UserGetRequest{ID: req.ID}, models.IncludeQuery{})
	if err == nil && existingResp != nil { // err is nil -> already exists
		return nil, fmt.Errorf("user with id %d %w", req.ID, repository.ErrAlreadyExists) // AlreadyExists
	}
	if err != nil && !errors.Is(err, repository.ErrNotFound) { //err is not ErrNotFound -> something went wrong
		return nil, err // Internal
	}
	// err IS ErrNotFound -> can Create()

	// Creating

	resp, err := s.userRepository.Create(ctx, req, include)
	if err != nil {
		return nil, err // Internal
	}

	return resp, nil // Created
}

// Based on [models.UserListRequest]
//
// Errors:
//
// - Validation
//   - [models.ErrInvalidLimit] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidOffset] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidSort] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidOrder] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidIncludeQuery] (Considered as http.StatusBadRequest)
//
// - Service
//
// - Repository
//   - [repository.ErrUnmarshalJSON] (Considered as http.StatusInternalServerError)
//     when database stores invalid JSON that can't be unmarshaled.
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (s *UserService) List(ctx context.Context, req models.UserListRequest, include models.IncludeQuery) (*models.UserListResponse, error) {

	// Validation

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Listing

	resp, err := s.userRepository.List(ctx, req, include)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Based on [models.UserGetRequest]
//
// Errors:
//
// - Validation
//   - [models.ErrInvalidID] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidIncludeQuery] (Considered as http.StatusBadRequest)
//
// - Service
//
// - Repository
//   - [repository.ErrNotFound] (Considered as http.StatusNotFound)
//   - [repository.ErrUnmarshalJSON] (Considered as http.StatusInternalServerError)
//     BUT JSON MUST BE VALIDATED BEFORE WRITING context_data TO DATABASE! DATABASE SHOULD ALWAYS STORE A VALID JSON!
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (s *UserService) Get(ctx context.Context, req models.UserGetRequest, include models.IncludeQuery) (*models.UserGetResponse, error) {

	// Validation

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Getting

	resp, err := s.userRepository.Get(ctx, req, include)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Based on [models.UserPatchRequest]
//
// Errors:
//
// - Validation
//   - [models.ErrInvalidID] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidTgTag] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidUserStatus] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidContextData] (Considered as http.StatusBadRequest)
//   - [models.ErrInvalidIncludeQuery] (Considered as http.StatusBadRequest)
//
// - Service
//
// - Repository
//   - [repository.ErrNotFound] (Considered as http.StatusNotFound)
//   - [reposiroty.ErrRequestConflict] (Considered as http.StatusConflict)
//   - [repository.ErrMarshalJSON] (Considered as http.StatusInternalServerError) FATAL!
//   - [repository.ErrUnmarshalJSON] (Considered as http.StatusInternalServerError) FATAL!
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (s *UserService) Patch(ctx context.Context, req models.UserPatchRequest, include models.IncludeQuery) (*models.UserPatchResponse, error) {

	// Validation

	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Getting current

	current, err := s.userRepository.Get(ctx, models.UserGetRequest{ID: req.ID}, include)
	if err != nil {
		return nil, err
	}

	// Patching

	resp, err := s.userRepository.Patch(ctx, req, &current.User, current.Version, include)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Based on [models.UserDeleteRequest]
//
// Errors:
//
// - Validation
//   - [models.ErrInvalidID] (Considered as http.StatusBadRequest)
//
// - Service
//
// - Repository
//   - [repository.ErrNotFound] (Considered as http.StatusNotFound)
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (s *UserService) Delete(ctx context.Context, req models.UserDeleteRequest) error {

	// Validation

	if err := req.Validate(); err != nil {
		return err
	}

	// Deletion

	err := s.userRepository.Delete(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
