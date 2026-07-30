package models

import (
	"errors"
	"strings"
)

type Exercise struct {
	ID   int64 `json:"id"`
	TgID int64 `json:"tg_id"`

	Name        string  `json:"name"`
	Description *string `json:"description"`
	Technique   *string `json:"technique"`
	WeightUnit  string  `json:"weight_unit"`

	Entries *[]ExerciseEntry `json:"entries"`
}

type ExerciseCreateRequest struct {
	TgID int64  `json:"tg_id" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type ExercisePatchRequest struct {
	ID int64 `json:"id"`

	Name        *string `json:"name"`
	Description *string `json:"description"`
	Technique   *string `json:"technique"`
	WeightUnit  *string `json:"weight_unit"`
}

type ExerciseListRequest struct {
	TgID   int64
	Offset int
	Limit  int
	SortBy string
	Order  Order
}

type Order string

const OrderASC Order = "ASC"
const OrderDESC Order = "DESC"

func ValidateOrder(order string) (Order, error) {
	switch strings.ToUpper(order) {
	case "ASC":
		return OrderASC, nil
	case "DESC":
		return OrderDESC, nil
	default:
		return OrderASC, errors.New("Invalid order")
	}
}
