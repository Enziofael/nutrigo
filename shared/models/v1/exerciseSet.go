package models

type ExerciseSet struct {
	ID         int64   `json:"id"`
	TgID       int64   `json:"tg_id"`
	EntryID    int64   `json:"entry_id"`
	
	Weight     float32 `json:"weight"`
	Reps       int     `json:"reps"`
	WeightUnit string  `json:"weight_unit"`
}
