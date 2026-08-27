package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	models "github.com/Enziofael/nutrigo/shared/models/v1"
	"github.com/Enziofael/nutrigo/shared/tools/postgreSql/query"
)

type UserPostgresRepository struct {
	db *sql.DB
}

func NewUserPostgresRepository(db *sql.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

var CreateQuery *query.QueryTemplate

func init() {
	CreateQuery = query.ClonePrototype(query.TEMPLATE_INSERT)
	CreateQuery.MSingle("table", "users")
	CreateQuery.MMultiple("columns", "id,tg_tag")
	CreateQuery.MMultiple("values", "($1, $2)")
	CreateQuery.MOptiobal("returning")
	CreateQuery.MMultiple("returning_columns", "id,tg_tag,status,version")
}

// Based on [models.UserCreateRequest]
//
// Errors:
//
// - Repository
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this)
func (r *UserPostgresRepository) Create(ctx context.Context, req models.UserCreateRequest, include models.IncludeQuery) (*models.UserCreateResponse, error) {

	// ====== VARS ======

	//Responce

	resp := models.UserCreateResponse{}

	q := CreateQuery.Clone()
	q.MMultipleIf("returning_columns", ",created_at", include.ContainsParam("created_at"))
	q.MMultipleIf("returning_columns", ",last_messaged_at", include.ContainsParam("last_messaged_at"))
	QUERY := q.Bake() + ";"

	// ====== Executing ======

	//Db querying

	row := r.db.QueryRowContext(ctx, QUERY, args...)

	// ====== Scanning ======

	//Scan
	err := row.Scan(dest...)
	if err != nil {
		return nil, err // Internal or constraint violation (should be pre-checked in service to avoid constraing violations)
	}

	//Nullable and parsable

	//Appending

	return &resp, nil
}

// Based on [models.UserListRequest]
//
// Errors:
//
// - Repository
//   - repository.[ErrUnmarshalJSON] (Considered as http.StatusInternalServerError)
//     when database stores invalid JSON that can't be unmarshaled.
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (r *UserPostgresRepository) List(ctx context.Context, req models.UserListRequest, include models.IncludeQuery) (*models.UserListResponse, error) {

	// ====== VARS ======

	//Responce

	resp := models.UserListResponse{
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   req.Sort,
		Order:  req.Order,
		Search: req.Search,
	}
	var ubuff models.User //scan buffer

	//Nullable and parsable fields

	var context sql.NullString
	var contextDataJSON sql.NullString

	//Query & scan args

	args := []any{}
	dest := []any{}

	// ====== Query building ======

	//Data part

	QSELECT := `SELECT id, tg_tag, status`
	dest = []any{&ubuff.ID, &ubuff.TgTag, &ubuff.Status}
	{
		//"created_at,last_messaged_at,context"
		if include.ContainsParam("created_at") {
			QSELECT += `, created_at `
			dest = append(dest, &ubuff.CreatedAt)
		}
		if include.ContainsParam("last_messaged_at") {
			QSELECT += `, last_messaged_at `
			dest = append(dest, &ubuff.LastMessagedAt)
		}
		if include.ContainsParam("context") {
			QSELECT += `, context `
			dest = append(dest, &context)
			QSELECT += `, context_data `
			dest = append(dest, &contextDataJSON)
		}
	}
	QFROM := `FROM users `

	//Logic part

	QWHERE := `WHERE `
	{
		if req.Search != "" {
			QWHERE += `similarity(tg_tag, $1) > $2 `
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
			orderBy = "similarity(tg_tag, $1)"
		} else {
			orderBy = "last_messaged_at"
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
		if include.ContainsParam("context") {
			if context.Valid {
				ubuff.Context = context.String
			} else {
				ubuff.Context = ""
			}
			if contextDataJSON.Valid {
				var contextData models.ContextData
				if err := json.Unmarshal([]byte(contextDataJSON.String), &contextData); err == nil {
					ubuff.ContextData = contextData
				} else {
					return nil, fmt.Errorf("can't parse ContextData json: %w", ErrUnmarshalJSON) // Internal FATAL! JSON MUST PARSE!
				}
			} else {
				ubuff.ContextData = models.ContextData{
					Values: make(map[string]string),
				}
			}
		}

		//Appending
		resp.Users = append(resp.Users, ubuff)
	}
	if err = rows.Err(); err != nil {
		return nil, err // Internal
	}

	return &resp, nil
}

// Based on [models.UserGetRequest]
//
// Errors:
//
// - Repository
//   - repository.[ErrNotFound] (Considered as http.StatusNotFound)
//   - repository.[ErrUnmarshalJSON] (Considered as http.StatusInternalServerError)
//     BUT JSON MUST BE VALIDATED BEFORE WRITING context_data TO DATABASE! DATABASE SHOULD ALWAYS STORE A VALID JSON!
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (r *UserPostgresRepository) Get(ctx context.Context, req models.UserGetRequest, include models.IncludeQuery) (*models.UserGetResponse, error) {

	// ====== VARS ======

	//Responce

	resp := models.UserGetResponse{}

	//Nullable & parsable fields

	var context sql.NullString
	var contextDataJSON sql.NullString

	//Query & scan args

	args := []any{}
	dest := []any{}

	// ====== Query building ======

	//Data part

	QSELECT := `SELECT id, tg_tag, status, version`
	dest = []any{&resp.User.ID, &resp.User.TgTag, &resp.User.Status, &resp.Version}
	{
		//"created_at,last_messaged_at,context,exercises?<exercises_params>"
		if include.ContainsParam("created_at") {
			QSELECT += `, created_at`
			dest = append(dest, &resp.User.CreatedAt)
		}
		if include.ContainsParam("last_messaged_at") {
			QSELECT += `, last_messaged_at`
			dest = append(dest, &resp.User.LastMessagedAt)
		}
		if include.ContainsParam("context") {
			QSELECT += `, context`
			dest = append(dest, &context)
			QSELECT += `, context_data`
			dest = append(dest, &contextDataJSON)
		}
		QSELECT += ` `
	}
	QFROM := `FROM users `

	//Logic part

	QWHERE := `WHERE id = $1 `
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

	if context.Valid { // context != null
		resp.User.Context = context.String
	} else { // = null
		resp.User.Context = ""
	}

	if contextDataJSON.Valid { // context_data != null
		err := json.Unmarshal([]byte(contextDataJSON.String), &resp.User.ContextData)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrUnmarshalJSON, err) // Internal FATAL! JSON MUST PARSE!
		}
	} else { // = null
		resp.User.ContextData = models.ContextData{
			Values: make(map[string]string),
		}
	}

	//Appending

	return &resp, nil
}

// Based on [models.UserPatchRequest]
//
// Errors:
//
// - Repository
//   - repository.[ErrNotFound] (Considered as http.StatusNotFound)
//   - reposiroty.[ErrRequestConflict] (Considered as http.StatusConflict)
//   - repository.[ErrMarshalJSON] (Considered as http.StatusInternalServerError) FATAL!
//   - repository.[ErrUnmarshalJSON] (Considered as http.StatusInternalServerError) FATAL!
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (r *UserPostgresRepository) Patch(ctx context.Context, req models.UserPatchRequest, current *models.User, version int64, include models.IncludeQuery) (*models.UserPatchResponse, error) {

	// ====== VARS ======

	//Responce

	resp := models.UserPatchResponse{}

	//Nullable & parsable fields

	var context sql.NullString
	var contextDataJSON sql.NullString

	//Query & scan args

	args := []any{}
	dest := []any{}

	// ====== Applying request's values changes to the current ======

	if req.UpdateLastMessagedAt {
		current.LastMessagedAt = time.Now()
	}
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

	if current.Context == "" { // if zero-value -> context = null
		context = sql.NullString{Valid: false}
	} else { // Valid not-null
		context = sql.NullString{String: current.Context, Valid: true}
	}

	if current.ContextData.IsZero() { // if zero-value -> contextDataJSON = null
		contextDataJSON = sql.NullString{Valid: false}
	} else if marshalled, err := json.Marshal(current.ContextData); err == nil { // Valid not-null
		contextDataJSON = sql.NullString{String: string(marshalled), Valid: true}
	} else {
		return nil, fmt.Errorf("%w: %w", ErrMarshalJSON, err) // Internal (marshalling error)
	}

	// ====== Query building ======

	//Logic part

	QUPDATE := `UPDATE users `

	QSET := `SET tg_tag = $1, status = $2, last_messaged_at = $3, context = $4, context_data = $5 `
	args = []any{current.TgTag, current.Status, current.LastMessagedAt, context, contextDataJSON} //$1 $2 $3 $4 $5

	QWHERE := `WHERE id = $6 AND version = $7`
	args = append(args, req.ID)  //$6
	args = append(args, version) //$7

	//Data part

	QRETURNING := `RETURNING id, tg_tag, status, version`
	dest = []any{&resp.User.ID, &resp.User.TgTag, &resp.User.Status, &resp.Version}
	{
		//"created_at,last_messaged_at,context,exercises<exercises_params>"
		if include.ContainsParam("created_at") {
			QRETURNING += `, created_at `
			dest = append(dest, &resp.User.CreatedAt)
		}
		if include.ContainsParam("last_messaged_at") {
			QRETURNING += `, last_messaged_at `
			dest = append(dest, &resp.User.LastMessagedAt)
		}
		if include.ContainsParam("context") {
			QRETURNING += `, context`
			dest = append(dest, &context)
			QRETURNING += `, context_data`
			dest = append(dest, &contextDataJSON)
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
		_, errget := r.Get(ctx, models.UserGetRequest{ID: req.ID}, models.IncludeQuery{})
		if errget != nil {
			if errors.Is(err, ErrNotFound) {
				return nil, fmt.Errorf("Can't Patch user with tg_id = %d: %w", req.ID, ErrNotFound) // NotFound
			} else {
				return nil, err // Internal (See Get() errors)
			}
		}
		return nil, fmt.Errorf("Can't Patch user with tg_id = %d: %w", req.ID, ErrRequestConflict)
	}

	//Nullable and parsable

	if context.Valid { // context != null
		resp.User.Context = context.String
	} else { // = null
		resp.User.Context = ""
	}

	if contextDataJSON.Valid { // context_data != null
		err := json.Unmarshal([]byte(contextDataJSON.String), &resp.User.ContextData)
		if err != nil {
			return nil, fmt.Errorf("%w - %w: %w", models.ErrInvalidContextData, ErrUnmarshalJSON, err) // Internal FATAL! JSON MUST PARSE!
		}
	} else { // = null
		resp.User.ContextData = models.ContextData{
			Values: make(map[string]string),
		}
	}

	//Appending

	return &resp, nil
}

// Based on [models.UserDeleteRequest]
//
// Errors:
//
// - Repository
//   - repository.[ErrNotFound] (Considered as http.StatusNotFound)
//
// - Database (Considered as http.StatusInternalServerError)
//   - internal errors
//   - constraint violation (Should be pre-checked in service to avoid this if there are any constraints)
func (r *UserPostgresRepository) Delete(ctx context.Context, req models.UserDeleteRequest) error {

	// ====== VARS ======

	//Responce

	//Nullable & parsable fields

	//Query & scan args

	args := []any{}

	// ====== Query building ======

	QUERY := `DELETE FROM users WHERE id = $1;`
	args = append(args, req.ID)

	// ====== Executing QUERY ======

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
