package models

import "errors"

var ErrValidation = errors.New("validation failed")
var ErrInvalidID = errors.New("invalid id")
var ErrInvalidUserID = errors.New("invalid user_id")
var ErrInvalidExerciseID = errors.New("invalid exercise_id")
var ErrInvalidTgTag = errors.New("invalid tg_tag")
var ErrInvalidName = errors.New("invalid name")
var ErrInvalidDescription = errors.New("invalid description")
var ErrInvalidSort = errors.New("invalid sort")
var ErrInvalidOrder = errors.New("invalid order")
var ErrInvalidUserStatus = errors.New("invalid user_status")
var ErrInvalidIncludeQuery = errors.New("invalid include")
var ErrInvalidLimit = errors.New("invalid limit")
var ErrInvalidOffset = errors.New("invalid offset")
var ErrInvalidContextData = errors.New("invelid context_data")
var ErrInvalidURL = errors.New("invalid url")
