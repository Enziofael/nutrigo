package repository

import (
	"database/sql"
)

type ExerciseEntryRepository interface {
}

type ExerciseEntryPostgresRepository struct {
	db *sql.DB
}

func NewExerciseEntryPostgresRepository(db *sql.DB) *ExerciseEntryPostgresRepository {
	return &ExerciseEntryPostgresRepository{db: db}
}
