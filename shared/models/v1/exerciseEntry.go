package models

import (
	"time"
)

// =========================================================
// MODEL
// =========================================================

type ExerciseEntry struct {
	ID         int64     `json:"id"          binding:"required"`
	ExerciseID int64     `json:"exercise_id" binding:"required"`
	Timestamp  time.Time `json:"timestamp"   binding:"required"`

	Comment string        `json:"comment,omitempty" binding:"omitempty,max=1000"`
	Sets    []ExerciseSet `json:"sets,omitempty"    binding:"omitempty,dive"`
}

type ExerciseSet struct {
	ID           int64                `json:"id,omitempty"       binding:"required"`
	EntryID      int64                `json:"entry_id,omitempty" binding:"required"`
	Measurements ExerciseMeasurements `json:"measurements"       binding:"required"`
	Reps         int                  `json:"reps"               binding:"omitempty,min=0"`
}

// SHOULD SYNC fields with [UnitPreferences] !

// # Kg & lbs:
//
// You CAN use both kg AND lbs
// in CREATE and PATCH requests
//
// for example, if an inventory with different units of measurement
// is used for one exercise.
//
// In responces both kg and lbs will be returned regardless.
// Weight preferences are only for handy input in forms
//
// # Km & Miles:
//
// In responces depending on preferences:
//   - Kilometers will be converted to miles (miles + miles(km))
//   - Miles will be converted to kilometers (km + km(miles))
//   - But NOT both at a time
type ExerciseMeasurements struct {
	Kg    float64       `json:"kg,omitempty"    binding:"omitempty,min=0"`
	Lbs   float64       `json:"lbs,omitempty"   binding:"omitempty,min=0"`
	Time  time.Duration `json:"time,omitempty"  binding:"omitempty,min=0"`
	Km    float64       `json:"km,omitempty"    binding:"omitempty,min=0"`
	Miles float64       `json:"miles,omitempty" binding:"omitempty,min=0"`
}

// =========================================================
// DTO
// =========================================================

// ===== EXERCISE LEVEL =====
// /exercises/:id/entry

// # POST /exercise/:id/entry
//
// Create entry by exercise id 
// with its sets
//
// Sets IDs will be ignored, you shouldn't fill them
// (keep zero values)
//
// Uri:
//   - exercise_id
//
// Body:
//   - sets
//   - [ comment ] = "" | <1000
//
// Query:
//   - [ include ] = "" | "comment,sets"+
type ExerciseEntryCreateRequest struct {
	ExerciseID int64 `json:"-" uri:"id" binding:"required"`

	Sets    []ExerciseSet `json:"sets"              binding:"required,dive"`
	Comment *string       `json:"comment,omitempty" binding:"omitempty,max=1000"`

	Include string `json:"-" form:"include"`
}
type ExerciseEntryCreateResponse struct {
	Entry ExerciseEntry `json:"entry"`
}

// # GET /exercises/:id/entry
//
// List entries
// by exercise id
//
// Uri:
//   - exercise_id
//
// Query:
//   - [ limit ] = 10 | 1-100
//   - [ offset ] = 0 | >=0
//   - [ sort ] = "timestamp" | "timestamp,exercise_id" (1)
//   - [ order ] = "DESC" | "ASC,DESC" (use [Order] constants)
//   - [ include ] = "" | "comment,sets"+
type ExerciseEntryListRequest struct {
	ExerciseID int64 `json:"-" uri:"id" binding:"required"`

	Limit   int    `json:"-" form:"limit"   binding:"omitempty,min=1,max=100"`
	Offset  int    `json:"-" form:"offset"  binding:"omitempty,min=0"`
	Sort    string `json:"-" form:"sort"    binding:"omitempty,oneof=timestamp exercise_id"`
	Order   Order  `json:"-" form:"order"   binding:"omitempty,oneof=ASC DESC"`
	Include string `json:"-" form:"include"`
}
type ExerciseEntryListResponse struct {
	Entries []ExerciseEntry `json:"entries"`
	Total   int             `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	Sort    string          `json:"sort"`
	Order   Order           `json:"order"`
}

// ===== EXERCISE_ENTRY LEVEL =====
// /exercise_entries/:id

// # GET /exercise_entries/:id
//
// Get entry
// by its id
//
// Uri:
//   - id
//
// Query:
//   - [ include ] = "" | "comment,sets" (Any combination)
type ExerciseEntryGetRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`

	Include string `json:"-" form:"include"`
}
type ExerciseEntryGetResponse struct {
	Entry ExerciseEntry `json:"entry"`
}

// # PATCH /exercise_entries/:id
//
// Patch entry by its id
// and overwrite all its sets
//
// Be careful not to delete or lose
// sets accidently by overwriting
// with new ones.
//
// Sets IDs will be ignored. You shouldn't fill them
// (keep zero values)
//
// Uri:
//	- id
//
// Body: (at least 1 should be provided. nil won't change the current value)
// Body:
//   - [ sets ]
//   - [ comment ] = "" | <1000
//
// Query:
//   - [ include ] = "" | "comment,sets"+
type ExerciseEntryPatchRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`

	Sets    *[]ExerciseSet `json:"sets,omitempty"    binding:"omitempty,dive"`
	Comment *string        `json:"comment,omitempty" binding:"omitempty,max=1000"`

	Include string `json:"-" form:"include"`
}
type ExerciseEntryPatchResponse struct {
	Entry ExerciseEntry `json:"entry"`
}


// # DELETE /exercise_entries/:id
//
// Delete entry by its id
// with all its sets
//
// Uri:
//   - id
type ExerciseEntryDeleteRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`
}
