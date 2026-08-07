package service

import "errors"

var ErrTgIdRequired = errors.New("tg_id is required")
var ErrTgTagRequired = errors.New("tg_tag is required")
var ErrInvalidStatus = errors.New("status is invalid")

var (
	ErrExerciseNotFound   = errors.New("exercise not found")
	ErrInvalidRequest     = errors.New("invalid exercise data")
	ErrIdRequired         = errors.New("id is required")
	ErrNameRequired       = errors.New("name is required")
	ErrWeightUnitRequired = errors.New("weight_unit is required")
)
