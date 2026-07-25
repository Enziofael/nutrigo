package models

type Exercise struct {
	ID   int64 `json:"id"`
	TgID int64 `json:"tg_id"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Technique   string `json:"technique"`
	WeightUnit  string `json:"weight_unit"`

	Entries *[]ExerciseEntry `json:"entries"`
}
