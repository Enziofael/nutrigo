package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	models "github.com/Enziofael/nutrigo/shared/models/v1"
	"github.com/lib/pq"
)

type UserRepository interface {
	Create(ctx context.Context, req models.UserCreateRequest) (*models.User, error)
	Delete(ctx context.Context, req models.UserDeleteRequest) error
	Get(ctx context.Context, req models.UserGetRequest) (*models.User, error)
	Patch(ctx context.Context, req models.UserPatchRequest, current *models.User) (*models.User, error)
}

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

// =========================================================
// CRUD REQUESTS
// =========================================================

func (r *UserPostgresRepository) Create(ctx context.Context, req models.UserCreateRequest) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO users (tg_id, tg_tag)
        VALUES ($1, $2)
        RETURNING tg_id, tg_tag, created_at, last_messaged_at, status, context, context_data`,
		req.TgID, req.TgTag).Scan(
		&u.TgID,
		&u.TgTag,
		&u.CreatedAt,
		&u.LastMessagedAt,
		&u.Status,
		&u.Context,
		&u.ContextData,
	)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return nil, fmt.Errorf("Can't create user with tg_id = %d: %w", req.TgID, ErrAlreadyExists)
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserPostgresRepository) Delete(ctx context.Context, req models.UserDeleteRequest) error {
	result, err := r.db.ExecContext(ctx, `
        DELETE 
		FROM users
        WHERE tg_id = $1`,
		req.TgID,
	)

	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("Can't Delete user with tg_id = %d: %w", req.TgID, ErrNotFound)
	}

	return nil
}

func (r *UserPostgresRepository) Get(ctx context.Context, req models.UserGetRequest) (*models.User, error) {
	var u models.User
	var context sql.NullString
	var contextDataJSON sql.NullString

	err := r.db.QueryRowContext(ctx, `
        SELECT tg_id, tg_tag, created_at, last_messaged_at, status, context, context_data
        FROM users
        WHERE tg_id = $1`,
		req.TgID).Scan(
		&u.TgID,
		&u.TgTag,
		&u.CreatedAt,
		&u.LastMessagedAt,
		&u.Status,
		&context,
		&contextDataJSON,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("Can't Get user with tg_id = %d: %w", req.TgID, ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	u.Context = context.String

	if !contextDataJSON.Valid {
		// Null context data
		u.ContextData = models.ContextData{MessageID: 0, Focus: 0, Values: make(map[string]string)}
	} else {
		var cd models.ContextData
		if err := json.Unmarshal([]byte(contextDataJSON.String), &cd); err == nil {
			u.ContextData = cd
		} else {
			return nil, fmt.Errorf("Can't Get user with tg_id = %d: %w", req.TgID, ErrUnmarshalJSON)
		}
	}

	return &u, nil
}

// При масштабировании может быть гонка данных и потерянные обновления (изменение между получаением current и обновлением)
func (r *UserPostgresRepository) Patch(ctx context.Context, req models.UserPatchRequest, current *models.User) (*models.User, error) {
	if req.TgTag != nil {
		current.TgTag = *req.TgTag
	}
	if req.Status != nil {
		current.Status = *req.Status
	}
	if req.Context != nil {
		current.Context = *req.Context
	}
	if req.ContextData != nil {
		current.ContextData = *req.ContextData
	}

	contextDataJSON, err := json.Marshal(current.ContextData)
	if err != nil {
		return nil, fmt.Errorf("Can't Patch user with tg_id = %d: %w", req.TgID, ErrMarshalJSON)
	}

	result, err := r.db.ExecContext(ctx, `
		UPDATE users 
		SET 
		tg_tag = $1, 
		status = $2, 
		context = $3,
		context_data = $4
		WHERE tg_id = $5`,
		current.TgTag,
		current.Status,
		current.Context,
		contextDataJSON,
		req.TgID,
	)
	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("Can't Patch user with tg_id = %d: %w", req.TgID, ErrNotFound)
	}

	return current, nil
}
