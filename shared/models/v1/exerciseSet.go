package models

type ExerciseSet struct {
	ID      int64 `json:"id"`
	TgID    int64 `json:"tg_id"`
	EntryID int64 `json:"entry_id"`

	WeightKg   float32 `json:"weight_kg"`
	WeightLbs  float32 `json:"weight_lbs"`
	Reps       int     `json:"reps"`
	WeightUnit string  `json:"weight_unit"`
}
