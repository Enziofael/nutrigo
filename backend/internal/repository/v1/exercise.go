package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type ExerciseRepository interface {
	Create(ctx context.Context, req models.ExerciseCreateRequest) (*models.Exercise, error)
	Delete(ctx context.Context, req models.ExerciseDeleteRequest) error
	Get(ctx context.Context, req models.ExerciseGetRequest) (*models.Exercise, error)
	List(ctx context.Context, req models.ExerciseListRequest) (*[]models.Exercise, error)
	Count(ctx context.Context, req models.ExerciseCountRequest) (int, error)
	Patch(ctx context.Context, req models.ExercisePatchRequest, current *models.Exercise) (*models.Exercise, error)
}

type ExercisePostgresRepository struct {
	db *sql.DB
}

func NewExercisePostgresRepository(db *sql.DB) *ExercisePostgresRepository {
	return &ExercisePostgresRepository{db: db}
}

// =========================================================
// CRUD REQUESTS
// =========================================================

func (r *ExercisePostgresRepository) Create(ctx context.Context, req models.ExerciseCreateRequest) (*models.Exercise, error) {
	var e models.Exercise
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO exercises (tg_id, name)
        VALUES ($1, $2) 
		RETURNING id, tg_id, name, description, technique, weight_unit, rating, created_at`,
		req.TgID, req.Name).Scan(
		&e.ID,
		&e.TgID,
		&e.Name,
		&e.Description,
		&e.Technique,
		&e.WeightUnit,
		&e.Rating,
		&e.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *ExercisePostgresRepository) Delete(ctx context.Context, req models.ExerciseDeleteRequest) (err error) {
	result, err := r.db.ExecContext(ctx, `
		DELETE 
		FROM exercises 
		WHERE id = $1`,
		req.ID,
	)

	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("Can't Delete exercise with id = %d: %w", req.ID, ErrNotFound)
	}

	return nil
}

func (r *ExercisePostgresRepository) Get(ctx context.Context, req models.ExerciseGetRequest) (*models.Exercise, error) {
	var e models.Exercise

	err := r.db.QueryRowContext(ctx, `
		SELECT id, tg_id, name, description, technique, weight_unit, rating, created_at
		FROM exercises 
		WHERE id = $1`,
		req.ID).Scan(
		&e.ID,
		&e.TgID,
		&e.Name,
		&e.Description,
		&e.Technique,
		&e.WeightUnit,
		&e.Rating,
		&e.CreatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("Can't Get exercise with id = %d: %w", req.ID, ErrNotFound)
	}
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (r *ExercisePostgresRepository) List(ctx context.Context, req models.ExerciseListRequest) (res *[]models.Exercise, err error) {
	query := `
		SELECT id, tg_id, name, description, technique, weight_unit, rating, created_at
		FROM exercises 
		WHERE 
		tg_id = $1`

	query += fmt.Sprintf(" ORDER BY %s %s", req.SortBy, req.Order)
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", req.Limit, req.Offset)

	rows, err := r.db.QueryContext(ctx, query, req.TgID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var ex models.Exercise
		if err := rows.Scan(
			&ex.ID,
			&ex.TgID,
			&ex.Name,
			&ex.Description,
			&ex.Technique,
			&ex.WeightUnit,
			&ex.Rating,
			&ex.CreatedAt); err != nil {
			return nil, err
		}
		exercises = append(exercises, ex)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &exercises, nil
}

func (r *ExercisePostgresRepository) Count(ctx context.Context, req models.ExerciseCountRequest) (count int, err error) {
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) 
		FROM exercises 
		WHERE tg_id = $1`,
		req.TgID).Scan(
		&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

// При масштабировании может быть гонка данных и потерянные обновления (изменение между получаением current и обновлением)
func (r *ExercisePostgresRepository) Patch(ctx context.Context, req models.ExercisePatchRequest, current *models.Exercise) (patched *models.Exercise, err error) {
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
		current.WeightUnit = models.WeightUnit(*req.WeightUnit)
	}
	if req.Rating != nil {
		current.Rating = *req.Rating
	}

	result, err := r.db.ExecContext(ctx, `
		UPDATE exercises 
		SET 
		name = $1, 
		description = $2, 
		technique = $3, 
		weight_unit = $4,
		rating = $5
		WHERE id = $6`,
		current.Name,
		current.Description,
		current.Technique,
		current.WeightUnit,
		current.Rating,
		req.ID,
	)
	if err != nil {
		return nil, err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return nil, fmt.Errorf("Can't Patch exercise with id = %d: %w", req.ID, ErrNotFound)
	}

	return current, nil
}
