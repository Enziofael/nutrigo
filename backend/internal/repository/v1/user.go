package v1

import (
	"context"
	"database/sql"
	"errors"

	models "github.com/Enziofael/nutrigo/shared/models/v1"
	"github.com/lib/pq"
)

type UserRepository interface {
	GetByTgId(ctx context.Context, tgID int64) (*models.User, error)
	Create(ctx context.Context, req models.UserCreateRequest) (*models.User, error)
	PatchStatus(ctx context.Context, tgID int64, req models.UserStatusUpdateRequest) (*models.User, error)
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
	err := r.db.QueryRowContext(ctx, `
        SELECT id, tg_id, tg_tag, created_at, last_messaged_at, status
        FROM users
        WHERE tg_id = $1
    `, tgID).Scan(&u.ID, &u.TgID, &u.TgTag, &u.CreatedAt, &u.LastMessagedAt, &u.Status)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
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

func (r *UserPostgresRepository) PatchStatus(ctx context.Context, tgID int64, req models.UserStatusUpdateRequest) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx, `
		UPDATE users
		SET status = $2
		WHERE tg_id = $1
		RETURNING id, tg_id, tg_tag, created_at, last_messaged_at, status
	`, tgID, req.Status).Scan(
		&u.ID,
		&u.TgID,
		&u.TgTag,
		&u.CreatedAt,
		&u.LastMessagedAt,
		&u.Status,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *UserPostgresRepository) Delete(ctx context.Context, tgID int64) error {
	var u models.User
	err := r.db.QueryRowContext(ctx, `
        DELETE 
		FROM users
        WHERE tg_id = $1
    `, tgID).Scan(&u.ID, &u.TgID, &u.TgTag, &u.CreatedAt, &u.LastMessagedAt, &u.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}
