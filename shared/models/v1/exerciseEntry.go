package models

import "time"

type ExerciseEntry struct {
	ID         int64 `json:"id"`
	TgID       int64 `json:"tg_id"`
	ExerciseID int64 `json:"exercise_id"`

	Timestamp time.Time `json:"timestamp"`
	Comment   string    `json:"comment"`

	Sets *[]ExerciseSet `json:"sets"`
}
