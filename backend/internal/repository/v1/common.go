package v1

import (
	"context"
	"database/sql"
)

type CommonRepository interface {
	Ping(ctx context.Context) (OK bool, err error)
}

type CommonPostgresRepository struct {
	db *sql.DB
}

func NewCommonPostgresRepository(db *sql.DB) *CommonPostgresRepository {
	return &CommonPostgresRepository{db: db}
}

func (r *CommonPostgresRepository) Ping(ctx context.Context) (bool, error) {
	var OK bool
	err := r.db.Ping()
	if err != nil {
		OK = false
	} else {
		OK = true
	}

	return OK, err
}
