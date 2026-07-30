package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "github.com/Enziofael/nutrigo/shared/models/v1"
	"github.com/lib/pq"
)

type ExerciseRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Exercise, error)
	ListByTgID(ctx context.Context, req models.ExerciseListRequest) (*[]models.Exercise, error)
	CountByTgID(ctx context.Context, tgID int64) (int, error)
	Create(ctx context.Context, req models.ExerciseCreateRequest) (*models.Exercise, error)
	Patch(ctx context.Context, ID int64, req models.ExercisePatchRequest) (*models.Exercise, error)
	Delete(ctx context.Context, ID int64) error
}

type ExercisePostgresRepository struct {
	db *sql.DB
}

func NewExercisePostgresRepository(db *sql.DB) *ExercisePostgresRepository {
	return &ExercisePostgresRepository{db: db}
}

func (r *ExercisePostgresRepository) GetByID(ctx context.Context, id int64) (*models.Exercise, error) {
	var ex models.Exercise
	err := r.db.QueryRowContext(ctx, `
		SELECT id, tg_id, name, description, technique, weight_unit 
		FROM exercises 
		WHERE id = $1
	`, id).Scan(&ex.ID, &ex.TgID, &ex.Name, &ex.Description, &ex.Technique, &ex.WeightUnit)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &ex, nil
}

func (r *ExercisePostgresRepository) ListByTgID(ctx context.Context, req models.ExerciseListRequest) (*[]models.Exercise, error) {
	query := `
		SELECT id, tg_id, name, description, technique, weight_unit 
		FROM exercises 
		WHERE 
		tg_id = $1`

	query += fmt.Sprintf(" ORDER BY %s %s", req.SortBy, req.Order)
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", req.Limit, req.Offset)

	rows, err := r.db.QueryContext(ctx, query,
		req.TgID)
	if err != nil {
		return nil, fmt.Errorf("list exercises: %w", err)
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var ex models.Exercise
		if err := rows.Scan(&ex.ID, &ex.TgID, &ex.Name, &ex.Description, &ex.Technique, &ex.WeightUnit); err != nil {
			return nil, fmt.Errorf("scan exercise: %w", err)
		}
		exercises = append(exercises, ex)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return &exercises, nil
}

func (r *ExercisePostgresRepository) CountByTgID(ctx context.Context, tgID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM exercises 
		WHERE tg_id = $1
	`, tgID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count exercises: %w", err)
	}
	return count, nil
}

func (r *ExercisePostgresRepository) Create(ctx context.Context, req models.ExerciseCreateRequest) (*models.Exercise, error) {
	var e models.Exercise
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO exercises (tg_id, name)
        VALUES ($1, $2) 
		RETURNING id, tg_id, name, description, technique, weight_unit`,
		req.TgID, req.Name).Scan(
		&e.ID,
		&e.TgID,
		&e.Name,
		&e.Description,
		&e.Technique,
		&e.WeightUnit,
	)

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		return nil, ErrExerciseAlreadyExists
	}
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *ExercisePostgresRepository) Patch(ctx context.Context, id int64, req models.ExercisePatchRequest) (*models.Exercise, error) {
	current, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil
	}

	if req.Name != nil {
		current.Name = *req.Name
	}
	if req.Description != nil {
		current.Description = req.Description
	}
	if req.Technique != nil {
		current.Technique = req.Technique
	}
	if req.WeightUnit != nil {
		current.WeightUnit = *req.WeightUnit
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE exercises 
		SET 
		name = $1, 
		description = $2, 
		technique = $3, 
		weight_unit = $4 
		WHERE id = $5`,
		current.Name,
		current.Description,
		current.Technique,
		current.WeightUnit,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("patch exercise: %w", err)
	}

	return current, nil
}

func (r *ExercisePostgresRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE 
		FROM exercises 
		WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("delete exercise: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrExerciseNotFound
	}
	return nil
}
