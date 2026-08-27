package models

import (
	"fmt"
	"net/url"
	"time"
)

func init() {
	exerciseCreateRequest_IncludeQuery = NewIncludeQuery(ExerciseCreateRequest{})
	exerciseListRequest_IncludeQuery = NewIncludeQuery(ExerciseListRequest{}, ExerciseEntryListRequest{})
	exerciseGetRequest_IncludeQuery = NewIncludeQuery(ExerciseGetRequest{}, ExerciseEntryListRequest{})
	exercisePatchRequest_IncludeQuery = NewIncludeQuery(ExercisePatchRequest{}, ExerciseEntryListRequest{})
}

// =========================================================
// MODEL
// =========================================================

type Exercise struct {
	ID     int64 `json:"id"      binding:"required,min=1"`
	UserID int64 `json:"user_id" binding:"required,min=1"`

	CreatedAt       time.Time       `json:"created_at,omitempty"       binding:"omitempty,min=0"`
	Name            string          `json:"name,omitempty"             binding:"omitempty,min=3,max=100"`
	Description     string          `json:"description,omitempty"      binding:"omitempty,max=1000"`
	URL             string          `json:"url,omitempty"              binding:"omitempty,url"`
	UnitPreferences UnitPreferences `json:"unit_preferences,omitempty" binding:"omitempty,min=1,max=31"` // SYNC WITH [UnitPreferences]
	Rating          int             `json:"rating,omitempty"           binding:"omitempty,min=0"`
	Entries         []ExerciseEntry `json:"entries,omitempty"          binding:"omitempty,dive"`
}

// =========================================================
// DTO
// =========================================================

// ====== USER LEVEL ======
// /users/:user_id/exercises

// # POST /users/:user_id/exercises
//
// Create exercise
// by user id
//
// Uri:
//   - user_id
//
// Body: (nil will user default value)
//   - name | 3-100
//   - [ description ] = "" | <1000
//   - [ url ] = "" | url
//   - [ unit_preferences ] = 0 | 0-31 (use [UnitPreferences] constants, combine them with "|" operator)
//
// Query:
//   - [ include ] = "" | "created_at,name,description,url,unit_preferences,rating" (Any combination)
type ExerciseCreateRequest struct {
	UserID int64 `json:"-" uri:"user_id" binding:"required,min=1"`

	Name            string           `json:"name"                       binding:"required,min=3,max=100"`
	Description     *string          `json:"description,omitempty"      binding:"omitempty,max=1000"`
	URL             *string          `json:"url,omitempty"              binding:"omitempty,url"`
	UnitPreferences *UnitPreferences `json:"unit_preferences,omitempty" binding:"omitempty,min=1,max=31"` // SYNC WITH [UnitPreferences]

	Include string `json:"-" form:"include"`
}
type ExerciseCreateResponse struct {
	Exercise Exercise `json:"exercise"`
	Version  int64    `json:"version"`
}

// Can return errors:
//   - [ErrInvalidUserID]
//   - [ErrInvalidName]
//   - [ErrInvalidDescription]
//   - [ErrInvalidUnitPreferences]
//   - [ErrInvalidIncludeQuery]
func (req ExerciseCreateRequest) Validate() error {
	if req.UserID <= 0 {
		return fmt.Errorf("%w: %w - must be positive. Actual: %d", ErrValidation, ErrInvalidUserID, req.UserID)
	}
	if len(req.Name) < 3 || len(req.Name) > 100 {
		return fmt.Errorf("%w: %w - len must be in range from 1 to 100 included. Actual: %d", ErrValidation, ErrInvalidName, len(req.Name))
	}
	if req.Description != nil && len(*req.Description) > 1000 {
		return fmt.Errorf("%w: %w - len can be 1000 max. Actual: %d", ErrValidation, ErrInvalidDescription, len(*req.Description))
	}
	if req.URL != nil {
		if u, err := url.Parse(*req.URL); err != nil || u.Host == "" || u.Scheme == "" {
			return fmt.Errorf("%w: %w - must be valid url. Actual: %s", ErrValidation, ErrInvalidURL, *req.URL)
		}
	}
	if req.UnitPreferences != nil && !req.UnitPreferences.Validate() {
		return fmt.Errorf("%w: %w - must be valid (see models.UnitPreferences). Actual: %b (%d)", ErrValidation, ErrInvalidUnitPreferences, *req.UnitPreferences, *req.UnitPreferences)
	}
	if !ValidateIncludeQuery(req.Include) {
		return fmt.Errorf("%w: %w - see format at IncludeQuery doc. Actual: \"%s\"", ErrValidation, ErrInvalidIncludeQuery, req.Include)
	}
	return nil
}

var exerciseCreateRequest_IncludeQuery IncludeQuery

func (ExerciseCreateRequest) AllowedParams() string {
	return "created_at,name,description,url,unit_preferences,rating"
}
func (req ExerciseCreateRequest) GetIncludeQuery() IncludeQuery {
	return exerciseCreateRequest_IncludeQuery
}

// # GET /users/:user_id/exercises
//
// List exercises
// by user id
//
// Uri:
//   - user_id
//
// Query:
//   - [ limit ] = 10 | 1-100
//   - [ offset ] = 0 | >=0
//   - [ sort ] = "relevance" | "created_at,relevance,rating" (1) (relevance: by search similarity if search is provided, rating otherwise)
//   - [ order ] = "DESC" | "ASC,DESC" (1)
//   - [ search ] = "" | Any search string (encode!)
//   - [ include ] = "" | "created_at,name,description,url,unit_preferences,rating" (Any combination)
type ExerciseListRequest struct {
	UserID int64 `json:"-" uri:"user_id" binding:"required,min=1"`

	Limit   int    `json:"-" form:"limit"   binding:"omitempty,min=1,max=100"`
	Offset  int    `json:"-" form:"offset"  binding:"omitempty,min=0"`
	Sort    string `json:"-" form:"sort"    binding:"omitempty,oneof=created_at relevance rating"`
	Order   Order  `json:"-" form:"order"   binding:"omitempty,oneof=ASC DESC"`
	Search  string `json:"-" form:"search"`
	Include string `json:"-" form:"include"`
}
type ExerciseListResponse struct {
	Exercises []Exercise `json:"exercises"`
	Total     int        `json:"total"`
	Limit     int        `json:"limit"`
	Offset    int        `json:"offset"`
	Sort      string     `json:"sort"`
	Order     Order      `json:"order"`
	Search    string     `json:"search,omitempty"`
}

// Can return errors:
//   - [ErrInvalidUserID]
//   - [ErrInvalidLimit]
//   - [ErrInvalidOffset]
//   - [ErrInvalidSort]
//   - [ErrInvalidOrder]
//   - [ErrInvalidIncludeQuery]
func (req *ExerciseListRequest) Validate() error {
	if req.UserID <= 0 {
		return fmt.Errorf("%w: %w - must be positive. Actual: %d", ErrValidation, ErrInvalidUserID, req.UserID)
	}
	if req.Limit < 1 || req.Limit > 100 {
		return fmt.Errorf("%w: %w - len must be in range from 1 to 100 included. Actual: %d", ErrValidation, ErrInvalidLimit, req.Limit)
	}
	if req.Offset < 0 {
		return fmt.Errorf("%w: %w - offset must be non-negative. Actual: %d", ErrValidation, ErrInvalidOffset, req.Offset)
	}
	if allowed := map[string]bool{"created_at": true, "relevance": true, "rating": true}; !allowed[req.Sort] {
		return fmt.Errorf("%w: %w - must be one of \"created_at\", \"relevance\", \"rating\". Actual: %s", ErrValidation, ErrInvalidSort, req.Sort)
	}
	if !req.Order.Validate() {
		return fmt.Errorf("%w: %w - must be one of \"ASC\", \"DESC\". Actual: %s", ErrValidation, ErrInvalidOrder, string(req.Order))
	}
	if !ValidateIncludeQuery(req.Include) {
		return fmt.Errorf("%w: %w - see format at IncludeQuery doc. Actual: \"%s\"", ErrValidation, ErrInvalidIncludeQuery, req.Include)
	}
	return nil
}

var exerciseListRequest_IncludeQuery IncludeQuery

func (ExerciseListRequest) AllowedParams() string {
	return "created_at,name,description,url,unit_preferences,rating"
}
func (req ExerciseListRequest) GetIncludeQuery() IncludeQuery {
	return exerciseListRequest_IncludeQuery
}

// ====== EXERCISE LEVEL ======
// /exercises/:id

// # GET /exercises/:id
//
// Get exercise
// by its id
//
// Uri:
//   - id
//
// Query:
//   - [ include ] = "" | "created_at,name,description,url,unit_preferences,rating,entries?<entries_params>" (Any combination)
//
// <entries_params>: see [ExerciseEntryListRequest] query params (encode!)
type ExerciseGetRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required,min=1"`

	Include string `json:"-" form:"include"`
}
type ExerciseGetResponse struct {
	Exercise Exercise `json:"exercise"`
	Version  int64    `json:"version"`
}

// Can return errors:
//   - [ErrInvalidID]
//   - [ErrInvalidIncludeQuery]
func (req *ExerciseGetRequest) Validate() error {
	if req.ID <= 0 {
		return fmt.Errorf("%w: %w - must be positive. Actual: %d", ErrValidation, ErrInvalidID, req.ID)
	}
	if !ValidateIncludeQuery(req.Include) {
		return fmt.Errorf("%w: %w - see format at IncludeQuery doc. Actual: \"%s\"", ErrValidation, ErrInvalidIncludeQuery, req.Include)
	}
	return nil
}

var exerciseGetRequest_IncludeQuery IncludeQuery

func (ExerciseGetRequest) AllowedParams() string {
	return "created_at,name,description,url,unit_preferences,rating,entries?<entries_params>"
}
func (req ExerciseGetRequest) GetIncludeQuery() IncludeQuery {
	return exerciseGetRequest_IncludeQuery
}

// # PATCH /exercises/:id
//
// Patch exercise
// by its id
//
// Uri:
//   - id
//
// Body: (at least 1 should be provided. nil won't change the current value)
//   - [ name ] = current | 3-100
//   - [ description ] = current | < 1000
//   - [ url ] = current | url
//   - [ unit_preferences ] = current |  0-31 (use [UnitPreferences] constants, combine them with "|" operator)
//
// Query:
//   - [ include ] = "" | "created_at,name,description,url,unit_preferences,rating,entries?<entries_params>" (Any combination)
//
// <entries_params>: see [ExerciseEntryListRequest] query params (encode!)
type ExercisePatchRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required,min=1"`

	Name            *string          `json:"name,omitempty"             binding:"omitempty,min=3,max=100"`
	Description     *string          `json:"description,omitempty"      binding:"omitempty,max=1000"`
	URL             *string          `json:"url,omitempty"              binding:"omitempty,url"`
	UnitPreferences *UnitPreferences `json:"unit_preferences,omitempty" binding:"omitempty,min=1,max=31"` // SYNC WITH [UnitPreferences]

	Include string `json:"-" form:"include"`
}
type ExercisePatchResponse struct {
	Exercise Exercise `json:"exercise"`
	Version  int64    `json:"version"`
}

// Can return errors:
//   - [ErrInvalidID]
//   - [ErrInvalidName]
//   - [ErrInvalidDescription]
//   - [ErrInvalidURL]
//   - [ErrInvalidUnitPreferences]
//   - [ErrInvalidIncludeQuery]
func (req *ExercisePatchRequest) Validate() error {
	if req.ID <= 0 {
		return fmt.Errorf("%w: %w - must be positive. Actual: %d", ErrValidation, ErrInvalidID, req.ID)
	}
	if req.Name != nil && (len(*req.Name) < 3 || len(*req.Name) > 100) {
		return fmt.Errorf("%w: %w - len must be in range from 1 to 100 included. Actual: %d", ErrValidation, ErrInvalidName, len(*req.Name))
	}
	if req.Description != nil && len(*req.Description) > 1000 {
		return fmt.Errorf("%w: %w - len can be 1000 max. Actual: %d", ErrValidation, ErrInvalidDescription, len(*req.Description))
	}
	if req.URL != nil {
		if u, err := url.Parse(*req.URL); err != nil || u.Host == "" || u.Scheme == "" {
			return fmt.Errorf("%w: %w - must be valid url. Actual: %s", ErrValidation, ErrInvalidURL, *req.URL)
		}
	}
	if req.UnitPreferences != nil && !req.UnitPreferences.Validate() {
		return fmt.Errorf("%w: %w - must be valid (see models.UnitPreferences). Actual: %b (%d)", ErrValidation, ErrInvalidUnitPreferences, *req.UnitPreferences, *req.UnitPreferences)
	}
	if !ValidateIncludeQuery(req.Include) {
		return fmt.Errorf("%w: %w - see format at IncludeQuery doc. Actual: \"%s\"", ErrValidation, ErrInvalidIncludeQuery, req.Include)
	}

	return nil
}

var exercisePatchRequest_IncludeQuery IncludeQuery

func (ExercisePatchRequest) AllowedParams() string {
	return "created_at,name,description,url,unit_preferences,rating,entries?<entries_params>"
}
func (req ExercisePatchRequest) GetIncludeQuery() IncludeQuery {
	return exercisePatchRequest_IncludeQuery
}

// # DELETE /exercises/:id
//
// Delete exercise
// by its id
//
// Uri:
//   - id
type ExerciseDeleteRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required,min=1"`
}

// Can return errors:
//   - [ErrInvalidID]
func (req *ExerciseDeleteRequest) Validate() error {
	if req.ID <= 0 {
		return fmt.Errorf("%w: %w - must be positive. Actual: %d", ErrValidation, ErrInvalidID, req.ID)
	}
	return nil
}
