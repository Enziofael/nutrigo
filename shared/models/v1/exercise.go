package models

import (
	"time"
)

// =========================================================
// MODEL
// =========================================================

type Exercise struct {
	ID        int64     `json:"id"         binding:"required"`
	TgID      int64     `json:"tg_id"      binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`

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
// /users/:tg_id/exercises

// # POST /users/:tg_id/exercises
//
// Create exercise
// by tg_id
//
// Uri:
//   - tg_id
//
// Body: (nil will user default value)
//   - name | 3-100
//   - [ description ] = "" | <1000
//   - [ url ] = "" | url
//   - [ unit_preferences ] = 0 | 0-31 (use [UnitPreferences] constants, combine them with "|" operator)
//
// Query:
//   - [ include ] = "" | "name,description,url,unit_preferences,rating,created_at" (Any combination)
type ExerciseCreateRequest struct {
	TgID int64 `json:"-" uri:"tg_id" binding:"required"`

	Name           string           `json:"name"                      binding:"required,min=3,max=100"`
	Description    *string          `json:"description,omitempty"     binding:"omitempty,max=1000"`
	URL            *string          `json:"url,omitempty"             binding:"omitempty,url"`
	UnitPreference *UnitPreferences `json:"unit_preferences,omitempty" binding:"omitempty,min=1,max=31"` // SYNC WITH [UnitPreferences]

	Include string `json:"-" form:"include"`
}
type ExerciseCreateResponse struct {
	Exercise Exercise `json:"exercise"`
}

// # GET /users/:tg_id/exercises
//
// List exercises
// by tg_id
//
// Uri:
//   - tg_id
//
// Query:
//   - [ limit ] = 10 | 1-100
//   - [ offset ] = 0 | >=0
//   - [ sort ] = "relevance" | "created_at,relevance" (1) (search relevance if search is provided, rating otherwise)
//   - [ order ] = "DESC" | "ASC,DESC" (1)
//   - [ search ] = "" | Any search string (encode!)
//   - [ include ] = "" | "name, description,url,unit_preferences,rating,entries?<entries_params>" (Any combination)
//
// <entries_params>: see [ExerciseEntryListRequest] query params (encode!)
type ExerciseListRequest struct {
	TgID int64 `json:"-" uri:"tg_id" binding:"required"`

	Limit   int    `json:"-" form:"limit"   binding:"omitempty,min=1,max=100"`
	Offset  int    `json:"-" form:"offset"  binding:"omitempty,min=0"`
	Sort    string `json:"-" form:"sort"    binding:"omitempty,oneof=created_at relevance"`
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
//   - [ include ] = "" | "name, description,url,unit_preferences,rating,entries?<entries_params>" (Any combination)
//
// <entries_params>: see [ExerciseEntryListRequest] query params (encode!)
type ExerciseGetRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`

	Include string `json:"-"`
}
type ExerciseGetResponse struct {
	Exercise Exercise `json:"exercise"`
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
//   - [ include ] = "" | "name, description,url,unit_preferences,rating,entries?<entries_params>" (Any combination)
//
// <entries_params>: see [ExerciseEntryListRequest] query params (encode!)
type ExercisePatchRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`

	Name            *string          `json:"name,omitempty"            binding:"omitempty,min=3,max=100"`
	Description     *string          `json:"description,omitempty"     binding:"omitempty,max=100"`
	URL             *string          `json:"url,omitempty"             binding:"omitempty,url"`
	UnitPreferences *UnitPreferences `json:"unit_preferences,omitempty" binding:"omitempty,min=1,max=31"` // SYNC WITH [UnitPreferences]

	Include string `json:"-"`
}
type ExercisePatchResponse struct {
	Exercise Exercise `json:"exercise"`
}

// # DELETE /exercises/:id
//
// Delete exercise
// by its id
//
// Uri:
//   - id
type ExerciseDeleteRequest struct {
	ID int64 `json:"-"`
}
