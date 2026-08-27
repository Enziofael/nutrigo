package repository

import (
	"context"

	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type UserRepository interface {
	Create(ctx context.Context, req models.UserCreateRequest, include models.IncludeQuery) (*models.UserCreateResponse, error)
	List(ctx context.Context, req models.UserListRequest, include models.IncludeQuery) (*models.UserListResponse, error)
	Get(ctx context.Context, req models.UserGetRequest, include models.IncludeQuery) (*models.UserGetResponse, error)
	Patch(ctx context.Context, req models.UserPatchRequest, current *models.User, version int64, include models.IncludeQuery) (*models.UserPatchResponse, error)
	Delete(ctx context.Context, req models.UserDeleteRequest) error
}

type ExerciseRepository interface {
	Create(ctx context.Context, req models.ExerciseCreateRequest, include models.IncludeQuery) (*models.ExerciseCreateResponse, error)
	List(ctx context.Context, req models.ExerciseListRequest, include models.IncludeQuery) (*models.ExerciseListResponse, error)
	Get(ctx context.Context, req models.ExerciseGetRequest, include models.IncludeQuery) (*models.ExerciseGetResponse, error)
	Patch(ctx context.Context, req models.ExercisePatchRequest, current *models.Exercise, version int64, include models.IncludeQuery) (*models.ExercisePatchResponse, error)
	Delete(ctx context.Context, req models.ExerciseDeleteRequest) error
}
