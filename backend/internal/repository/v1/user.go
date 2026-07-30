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
	GetByTgId(ctx context.Context, tgID int64) (*models.User, error)
	Create(ctx context.Context, req models.UserCreateRequest) (*models.User, error)
	Patch(ctx context.Context, tgID int64, req models.UserPatchRequest) (*models.User, error)
	Delete(ctx context.Context, tgID int64) error
}

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

func (r *UserPostgresRepository) GetByTgId(ctx context.Context, tgID int64) (*models.User, error) {
	var u models.User
	var context sql.NullString
	var contextDataJSON sql.NullString

	err := r.db.QueryRowContext(ctx, `
        SELECT id, tg_id, tg_tag, created_at, last_messaged_at, status, context, context_data
        FROM users
        WHERE tg_id = $1
    `, tgID).Scan(&u.ID, &u.TgID, &u.TgTag, &u.CreatedAt, &u.LastMessagedAt, &u.Status, &context, &contextDataJSON)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	u.Context = context.String // если NULL → пустая строка

	if contextDataJSON.Valid {
		var cd models.ContextData
		if err := json.Unmarshal([]byte(contextDataJSON.String), &cd); err == nil {
			u.ContextData = cd
		} else {
			u.ContextData = models.ContextData{MessageID: 0, Focus: 0, Values: make(map[string]string)}
		}
	} else {
		u.ContextData = models.ContextData{MessageID: 0, Focus: 0, Values: make(map[string]string)}
	}

	return &u, nil
}

func (r *UserPostgresRepository) Create(ctx context.Context, req models.UserCreateRequest) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `
        INSERT INTO users (tg_id, tg_tag)
        VALUES ($1, $2)
        RETURNING id, tg_id, tg_tag, created_at, last_messaged_at, status
    `, req.TgID, req.TgTag).Scan(
		&u.ID,
		&u.TgID,
		&u.TgTag,
		&u.CreatedAt,
		&u.LastMessagedAt,
		&u.Status,
	)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return nil, ErrUserAlreadyExists
	}
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserPostgresRepository) Patch(ctx context.Context, tgID int64, req models.UserPatchRequest) (*models.User, error) {
	current, err := r.GetByTgId(ctx, tgID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, ErrUserNotFound
	}

	if req.TgTag != nil {
		current.TgTag = *req.TgTag
	}
	if req.LastMessagedAt != nil {
		current.LastMessagedAt = *req.LastMessagedAt
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
		return nil, fmt.Errorf("marchal context data: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE users 
		SET 
		tg_tag = $1, 
		last_messaged_at = $2, 
		status = $3, 
		context = $4,
		context_data = $5
		WHERE tg_id = $6`,
		current.TgTag,
		current.LastMessagedAt,
		current.Status,
		current.Context,
		contextDataJSON,
		tgID,
	)
	if err != nil {
		return nil, fmt.Errorf("patch exercise: %w", err)
	}

	return current, nil
}

func (r *UserPostgresRepository) Delete(ctx context.Context, tgID int64) error {
	result, err := r.db.ExecContext(ctx, `
        DELETE 
		FROM users
        WHERE tg_id = $1
    `, tgID)
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrUserNotFound
	}
	return nil
}
