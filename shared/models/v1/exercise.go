package models

import (
	"errors"
	"fmt"
	"time"
)

// Model
type Exercise struct {
	ID          int64      `json:"id"                    binding:"required,min=1"     db:"id"`
	TgID        int64      `json:"tg_id"                 binding:"required,min=1"     db:"tg_id"`
	Name        string     `json:"name"                  binding:"required,max=140"   db:"name"`
	Description *string    `json:"description,omitempty" binding:"omitempty,max=1000" db:"description"`
	Technique   *string    `json:"technique,omitempty"   binding:"omitempty,url"      db:"technique"`
	WeightUnit  WeightUnit `json:"weight_unit"           binding:"required"           db:"weight_unit"`
	Rating      int        `json:"rating"                binding:"min=-100,max=100"   db:"rating"`
	CreatedAt   time.Time  `json:"created_at"            binding:"required"`

	Entries *[]ExerciseEntry `json:"entries,omitempty"`
}

func SanitizeWeightUnitString(s string) (string, error) {
	switch s {
	case string(WeightUnit_Kg):
		return string(WeightUnit_Kg), nil
	case string(WeightUnit_Lbs):
		return string(WeightUnit_Lbs), nil
	case string(WeightUnit_Mixed):
		return string(WeightUnit_Mixed), nil
	default:
		return string(WeightUnit_Kg), InvalidWeightUnit
	}
}

func SanitizeRating(r int) (int, error) {
	if r < -100 || r > 100 {
		return 0, errors.New("Rating should be between -100 & 100")
	}
	return r, nil
}

func SanitizeExerciseSortBy(s string) (string, error) {
	switch s {
	case "id", "name", "rating":
		return s, nil
	case "tg_id", "description", "technique", "weight_unit", "entries":
		return "", fmt.Errorf("Forbidden sort field: %s", s)
	default:
		return "", fmt.Errorf("Unknown sort field: %s", s)
	}
}

// =========================================================
// CRUD REQUESTS
// =========================================================

// Create by user
type ExerciseCreateRequest struct {
	//body
	TgID int64  `json:"tg_id" binding:"required"`
	Name string `json:"name"  binding:"required"`
}

func (req ExerciseCreateRequest) Sanitize() (ExerciseCreateRequest, error) {
	var errTgID, errName error

	req.TgID, errTgID = SanitizeTgID(req.TgID)
	req.Name, errName = SanitizeName(req.Name)

	err := errors.Join(errTgID, errName)
	return req, err
}

// Delete by id
type ExerciseDeleteRequest struct {
	//url
	ID int64 `json:"id" binding:"required,min=1"`
}

func (req ExerciseDeleteRequest) Sanitize() (ExerciseDeleteRequest, error) {
	var errID error

	req.ID, errID = SanitizeID(req.ID)

	err := errors.Join(errID)
	return req, err
}

// Get by id
type ExerciseGetRequest struct {
	//url
	ID int64 `json:"id" binding:"required,min=1"`
}

func (req ExerciseGetRequest) Sanitize() (ExerciseGetRequest, error) {
	var errID error

	req.ID, errID = SanitizeID(req.ID)

	err := errors.Join(errID)
	return req, err
}

// List by user
type ExerciseListRequest struct {
	//url
	TgID int64 `json:"tg_id"  binding:"required,min=1"`

	//body
	Offset int    `json:"offset" binding:"required,min=0"`
	Limit  int    `json:"limit"  binding:"required,max=100"`
	SortBy string `json:"sort"   binding:"required"`
	Order  Order  `json:"order"  binding:"required"`
}

func (req ExerciseListRequest) Sanitize() (ExerciseListRequest, error) {
	var errTgID, errOffset, errLimit, errSortBy, errOrder error

	req.TgID, errTgID = SanitizeTgID(req.TgID)
	req.Offset, errOffset = SanitizeOffset(req.Offset)
	req.Limit, errLimit = SanitizeLimit(req.Limit)
	req.SortBy, errSortBy = SanitizeExerciseSortBy(req.SortBy)
	req.Order, errOrder = SanitizeOrder(req.Order)

	err := errors.Join(errTgID, errOffset, errLimit, errSortBy, errOrder)
	return req, err
}

// Count by user
type ExerciseCountRequest struct {
	//url
	TgID int64 `json:"tg_id" binding:"required,min=1"`
}

func (req ExerciseCountRequest) Sanitize() (ExerciseCountRequest, error) {
	var errTgID error

	req.TgID, errTgID = SanitizeTgID(req.TgID)

	err := errors.Join(errTgID)
	return req, err
}

// Patch by id
type ExercisePatchRequest struct {
	//url
	ID int64 `json:"id" binding:"required,min=1"`

	//body optional (at least 1)
	Name        *string `json:"name,omitempty"         binding:"max=140"`
	Description *string `json:"description,omitempty"  binding:"max=1000"`
	Technique   *string `json:"technique,omitempty"    binding:"url"`
	WeightUnit  *string `json:"weight_unit,omitempty"`
	Rating      *int    `json:"rating,omitempty"       binding:"min=-100,max=100"`
}

func (req ExercisePatchRequest) Sanitize() (ExercisePatchRequest, error) {
	var errID, errName, errDescription, errTechnique, errWeightUnit, errRating error

	req.ID, errID = SanitizeID(req.ID)
	if req.Name != nil {
		var n string
		n, errName = SanitizeName(*req.Name)
		req.Name = &n
	}
	if req.Description != nil {
		var d string
		d, errDescription = SanitizeDescription(*req.Description)
		req.Description = &d
	}
	if req.Technique != nil {
		var t string
		t, errTechnique = SanitizeUrl(*req.Technique)
		req.Technique = &t
	}
	if req.WeightUnit != nil {
		var w string
		w, errWeightUnit = SanitizeWeightUnitString(*req.WeightUnit)
		req.WeightUnit = &w
	}
	if req.Rating != nil {
		var r int
		r, errRating = SanitizeRating(*req.Rating)
		req.Rating = &r
	}

	err := errors.Join(errID, errName, errDescription, errTechnique, errWeightUnit, errRating)
	return req, err
}
