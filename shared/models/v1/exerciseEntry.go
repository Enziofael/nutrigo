package models

import "time"

type ExerciseEntry struct {
	ID         int64     `json:"id"                binding:"required,min=1"`
	TgID       int64     `json:"tg_id"             binding:"required,min=1"`
	ExerciseID int64     `json:"exercise_id"       binding:"required,min=1"`
	Timestamp  time.Time `json:"timestamp"         binding:"required"`
	Comment    *string   `json:"comment,omitempty" binding:"omitempty,max=1000"`

	Sets *[]ExerciseSet `json:"sets" binding:"required"`
}

type ExerciseEntryCreateRequest struct {
	TgID       int64
	ExerciseID int64
	Comment    *string
}

type ExerciseEntryGetRequest struct {
	ID int64
}

type ExerciseEntryListRequest struct {
	ExerciseID int64
}
