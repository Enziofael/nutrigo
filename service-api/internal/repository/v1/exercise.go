package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	models "github.com/Enziofael/nutrigo/shared/models/v1"
)

type ExercisePostgresRepository struct {
	db *sql.DB
}

func NewExercisePostgresRepository(db *sql.DB) *ExercisePostgresRepository {
	return &ExercisePostgresRepository{db: db}
}

// Based on [models.ExerciseCreateRequest]
//
// Errors:
//
// - Repository
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this)
func (r *ExercisePostgresRepository) Create(ctx context.Context, req models.ExerciseCreateRequest, include models.IncludeQuery) (*models.ExerciseCreateResponse, error) {

	// ====== VARS ======

	//Responce

	resp := models.ExerciseCreateResponse{}

	//Nullable and parsable fields

	//Query & scan args

	args := []any{}
	dest := []any{}

	// ====== Query building ======

	//Logic part

	QINSERT := `INSERT INTO exercises (user_id, name`
	QVALUES := `VALUES ($1, $2`
	args = []any{req.UserID, req.Name} //$1 $2
	{                                  //Optional fields
		argsCounter := len(args)
		if req.Description != nil {
			argsCounter++
			QINSERT += `, description`
			QVALUES += `, $` + strconv.Itoa(argsCounter) //n=3 $n
			args = append(args, req.Description)
		}
		if req.URL != nil {
			argsCounter++
			QINSERT += `, url`
			QVALUES += `, $` + strconv.Itoa(argsCounter) //n=3-4 $n
			args = append(args, req.URL)
		}
		if req.UnitPreferences != nil {
			argsCounter++
			QINSERT += `, unit_preferences`
			QVALUES += `, $` + strconv.Itoa(argsCounter) //n=3-5 $n
			args = append(args, req.UnitPreferences)
		}
		QINSERT += `) `
		QVALUES += `) `
	}

	//Data part

	QRETURNING := `RETURNING id, user_id, version`
	dest = []any{&resp.Exercise.ID, &resp.Exercise.UserID, &resp.Version} //mandatory returns
	{                                                                     //optional returns
		//include query param based returns
		//"created_at,name,description,url,unit_preferences,rating"
		if include.ContainsParam("created_at") {
			QRETURNING += `, created_at`
			dest = append(dest, &resp.Exercise.CreatedAt)
		}
		if include.ContainsParam("name") {
			QRETURNING += `, name`
			dest = append(dest, &resp.Exercise.Name)
		}
		if include.ContainsParam("description") {
			QRETURNING += `, description`
			dest = append(dest, &resp.Exercise.Description)
		}
		if include.ContainsParam("url") {
			QRETURNING += `, url`
			dest = append(dest, &resp.Exercise.URL)
		}
		if include.ContainsParam("unit_preferences") {
			QRETURNING += `, unit_preferences`
			dest = append(dest, &resp.Exercise.UnitPreferences)
		}
		if include.ContainsParam("rating") {
			QRETURNING += `, rating`
			dest = append(dest, &resp.Exercise.Rating)
		}
	}

	//Concat queries

	QUERY := QINSERT + QVALUES + QRETURNING + `;`

	// ====== Executing ======

	//Db querying

	row := r.db.QueryRowContext(ctx, QUERY, args...)

	//Scanning

	err := row.Scan(dest...)
	if err != nil {
		return nil, err // Internal or constraint violation (should be pre-checked in service to avoid constraing violations)
	}

	return &resp, nil
}

// Based on [models.ExerciseListRequest]
//
// Errors:
//
// - Repository
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (r *ExercisePostgresRepository) List(ctx context.Context, req models.ExerciseListRequest, include models.IncludeQuery) (*models.ExerciseListResponse, error) {

	// ====== VARS ======

	//Responce

	resp := models.ExerciseListResponse{
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   req.Sort,
		Order:  req.Order,
		Search: req.Search,
	}
	var ebuff models.Exercise //scan buffer

	//Nullable and parsable fields

	//Query & scan args

	args := []any{}
	dest := []any{}

	// ====== Query building ======

	//Data part

	QSELECT := `SELECT id, user_id`
	dest = []any{&ebuff.ID, &ebuff.UserID}
	{
		//"created_at,name,description,url,unit_preferences,rating"
		if include.ContainsParam("created_at") {
			QSELECT += `, created_at`
			dest = append(dest, &ebuff.CreatedAt)
		}
		if include.ContainsParam("name") {
			QSELECT += `, name`
			dest = append(dest, &ebuff.Name)
		}
		if include.ContainsParam("description") {
			QSELECT += `, description`
			dest = append(dest, &ebuff.Description)
		}
		if include.ContainsParam("url") {
			QSELECT += `, url`
			dest = append(dest, &ebuff.URL)
		}
		if include.ContainsParam("unit_preferences") {
			QSELECT += `, unit_preferences`
			dest = append(dest, &ebuff.UnitPreferences)
		}
		if include.ContainsParam("rating") {
			QSELECT += `, rating`
			dest = append(dest, &ebuff.Rating)
		}
	}
	QFROM := `FROM exercises`

	//Logic part

	QWHERE := `WHERE `
	{
		if req.Search != "" {
			QWHERE += `GREATEST(similarity(name, $1),similarity(COALESCE(description, ''), $1)) > $2`
			args = append(args, req.Search)             //$1
			args = append(args, models.SearchThreshold) //$2
		} else {
			QWHERE = ``
		}
	}

	QORDERBY := `ORDER BY `
	var orderBy string
	if req.Sort == "relevance" {
		if req.Search != "" {
			orderBy = `GREATEST(similarity(name, $1),similarity(COALESCE(description, ''), $1))`
		} else {
			orderBy = "rating"
		}
	} else {
		orderBy = req.Sort
	}
	QORDERBY += orderBy + ` ` + string(req.Order) + ` `

	QLIMIT := `LIMIT $3 `
	args = append(args, req.Limit) //$3

	QOFFSET := `OFFSET $4`
	args = append(args, req.Offset) //$4

	//Concat queries

	QUERYCOUNT := `SELECT COUNT(*)` + QFROM + QWHERE
	QUERY := QSELECT + QFROM + QWHERE + QORDERBY + QLIMIT + QOFFSET + `;`

	// ====== Executing QUERYCOUNT ======

	err := r.db.QueryRowContext(ctx, QUERYCOUNT, args[0:2]...).Scan(&resp.Total)
	if err != nil {
		return nil, err // Internal
	}

	// ====== Executing QUERY ======

	rows, err := r.db.QueryContext(ctx, QUERY, args...)
	if err != nil {
		return nil, err // Internal
	}
	defer rows.Close()

	// ====== Scanning ======

	for rows.Next() {
		//Scan
		err = rows.Scan(dest...)
		if err != nil {
			return nil, fmt.Errorf("scan user: %w", err) // Internal
		}

		//Nullable and parsable

		//Appending
		resp.Exercises = append(resp.Exercises, ebuff)
	}
	if err := rows.Err(); err != nil {
		return nil, err // Internal
	}

	return &resp, nil
}

// Based on [models.ExerciseGetRequest]
//
// Errors:
//
// - Repository
//   - repository.[ErrNotFound] (Considered as http.StatusNotFound)
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (r *ExercisePostgresRepository) Get(ctx context.Context, req models.ExerciseGetRequest, include models.IncludeQuery) (*models.ExerciseGetResponse, error) {

	// ====== VARS ======

	// Responce

	resp := models.ExerciseGetResponse{}

	//Nullable and parsable fields

	//Query & scan args

	args := []any{}
	dest := []any{}

	// ====== Query building ======

	//Data part

	QSELECT := `SELECT id, user_id, version`
	dest = []any{&resp.Exercise.ID, &resp.Exercise.UserID, &resp.Version}
	{
		//"created_at,name,description,url,unit_preferences,rating,entries?<entries_params>"
		if include.ContainsParam("created_at") {
			QSELECT += `, created_at`
			dest = append(dest, &resp.Exercise.CreatedAt)
		}
		if include.ContainsParam("name") {
			QSELECT += `, name`
			dest = append(dest, &resp.Exercise.Name)
		}
		if include.ContainsParam("description") {
			QSELECT += `, description`
			dest = append(dest, &resp.Exercise.Description)
		}
		if include.ContainsParam("url") {
			QSELECT += `, url`
			dest = append(dest, &resp.Exercise.URL)
		}
		if include.ContainsParam("unit_preferences") {
			QSELECT += `, unit_preferences`
			dest = append(dest, &resp.Exercise.UnitPreferences)
		}
		if include.ContainsParam("rating") {
			QSELECT += `, rating`
			dest = append(dest, &resp.Exercise.Rating)
		}
		QSELECT += ` `
	}
	QFROM := `FROM exercises `

	//Logic part

	QWHERE := `WHERE id = $1`
	args = []any{req.ID} //$1

	//Concat queries

	QUERY := QSELECT + QFROM + QWHERE + `LIMIT 1;`

	// ====== Executing QUERY ======

	row := r.db.QueryRowContext(ctx, QUERY, args...)

	// ====== Scanning ======

	//Scan
	err := row.Scan(dest...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user %w: id = %d", ErrNotFound, req.ID) // NotFound
		}
		return nil, err // Internal
	}

	//Nullable and parsable

	//Appending

	return &resp, nil
}

// Based on [models.ExercisePatchRequest]
//
// Errors:
//
// - Repository
//   - repository.[ErrNotFound] (Considered as http.StatusNotFound)
//   - reposiroty.[ErrRequestConflict] (Considered as http.StatusConflict)
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (r *ExercisePostgresRepository) Patch(ctx context.Context, req models.ExercisePatchRequest, current *models.Exercise, version int64, include models.IncludeQuery) (*models.ExercisePatchResponse, error) {

	// ====== VARS ======

	//Responce

	resp := models.ExercisePatchResponse{}

	//Nullable & parsable fields

	//Query & scan args

	args := []any{}
	dest := []any{}

	// ====== Applying request's values changes to the current ======

	if req.Name != nil {
		current.Name = *req.Name
	}
	if req.Description != nil {
		current.Description = *req.Description
	}
	if req.URL != nil {
		current.URL = *req.URL
	}
	if req.UnitPreferences != nil {
		current.UnitPreferences = *req.UnitPreferences
	}

	// ====== Query building ======

	//Logic part

	QUPDATE := `UPDATE exercises `

	QSET := `SET name = $1, description = $2, url = $3, unit_preferences = $4 `
	args = []any{current.Name, current.Description, current.URL, current.UnitPreferences} //$1 $2 $3 $4

	QWHERE := `WHERE id = $5 AND version = $6`
	args = append(args, req.ID)  //$5
	args = append(args, version) //$6

	//Data part

	QRETURNING := `RETURNING id, user_id, version`
	dest = []any{&resp.Exercise.ID, &resp.Exercise.UserID, &resp.Version}
	{
		//"created_at,name,description,url,unit_preferences,rating,entries?<entries_params>"
		if include.ContainsParam("created_at") {
			QRETURNING += `, created_at`
			dest = append(dest, &resp.Exercise.CreatedAt)
		}
		if include.ContainsParam("name") {
			QRETURNING += `, name`
			dest = append(dest, &resp.Exercise.Name)
		}
		if include.ContainsParam("description") {
			QRETURNING += `, description`
			dest = append(dest, &resp.Exercise.Description)
		}
		if include.ContainsParam("url") {
			QRETURNING += `, url`
			dest = append(dest, &resp.Exercise.URL)
		}
		if include.ContainsParam("unit_preferences") {
			QRETURNING += `, unit_preferences`
			dest = append(dest, &resp.Exercise.UnitPreferences)
		}
		if include.ContainsParam("rating") {
			QRETURNING += `, rating`
			dest = append(dest, &resp.Exercise.Rating)
		}
		QRETURNING += ` `
	}

	//Concat queries

	QUERY := QUPDATE + QSET + QWHERE + QRETURNING + `;`

	// ====== Executing QUERY ======

	row := r.db.QueryRowContext(ctx, QUERY, args...)

	// ====== Scanning ======

	//Scan
	err := row.Scan(dest...)
	if err != nil {
		_, errget := r.Get(ctx, models.ExerciseGetRequest{ID: req.ID}, models.IncludeQuery{})
		if errget != nil {
			if errors.Is(err, ErrNotFound) {
				return nil, fmt.Errorf("Can't Patch exercise with id = %d: %w", req.ID, ErrNotFound) // NotFound
			} else {
				return nil, err // Internal (See Get() errors)
			}
		}
		return nil, fmt.Errorf("Can't Patch exercise with id = %d: %w", req.ID, ErrRequestConflict)
	}

	//Nullable and parsable

	//Appending

	return &resp, nil
}

// Based on [models.ExerciseDeleteRequest]
//
// Errors:
//
// - Repository
//   - repository.[ErrNotFound] (Considered as http.StatusNotFound)
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (r *ExercisePostgresRepository) Delete(ctx context.Context, req models.ExerciseDeleteRequest) error {

	// ====== VARS ======

	//Responce

	//Nullable & parsable fields

	//Query & scan args

	args := []any{}

	// ====== Query building ======

	QUERY := `DELETE FROM exercises WHERE id = $1;`
	args = append(args, req.ID)

	// ====== Executing query ======

	result, err := r.db.ExecContext(ctx, QUERY, args...)
	if err != nil {
		return err // Internal
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err // Internal FATAL (not supported)
	}

	if rows == 0 {
		return fmt.Errorf("%w: id = %d", ErrNotFound, req.ID) // NotFound
	}

	return nil
}
