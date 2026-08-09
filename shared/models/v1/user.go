package models

import (
	"time"
)

// =========================================================
// MODEL
// =========================================================

type User struct {
	TgID   int64      `json:"tg_id"  binding:"required"`
	TgTag  string     `json:"tg_tag" binding:"required"`
	Status UserStatus `json:"status" binding:"required,oneof=unknown requested confirmed restricted banned admin"` //SYNC WITH [UserStatus]

	CreatedAt      time.Time   `json:"created_at,omitempty"       binding:"omitempty,min=0"`
	LastMessagedAt time.Time   `json:"last_messaged_at,omitempty" binding:"omitempty,min=0"`
	Context        string      `json:"context,omitempty"`       //binding any
	ContextData    ContextData `json:"context_data,omitempty"`  //binding any
	Exercises      []Exercise  `json:"exercises,omitempty"        binding:"omitempty,dive"`
}

type ContextData struct {
	MessageID int               `json:"message_id"         binding="required"`
	Focus     int               `json:"focus"`           //binding any
	Values    map[string]string `json:"values,omitempty"`//binding any
}

// =========================================================
// DTO
// =========================================================

// ====== GLOBAL LEVEL ======
// /users

// # POST /users
//
// Create user
// by tg id
//
// Body:
//   - tg_id
//   - tg_tag
//
// Query:
//   - [ include ] = "" | "created_at,last_messaged_at" (Any combination)
type UserCreateRequest struct {
	TgID  int64  `json:"tg_id"  binding="required"`
	TgTag string `json:"tg_tag" binding="required,min=3,max=100"`

	Include string `json:"-" form:"include"`
}
type UserCreateResponse struct {
	User User `json:"user"`
}

// # GET /users
//
// List
// users
//
// Query:
//   - [ limit ] = 10 | 0 - 100
//   - [ offset ] = 0 | >= 0
//   - [ sort ] = "last_messaged_at" | "last_messaged_at,created_at" (1)
//   - [ order ] = "DESC" | "ASC,DESC" (1)
//   - [ search ] = "" | Any search string (encode!)
//   - [ include ] = "" | "created_at,last_messaged_at,context" (Any combination)
type UserListRequest struct {
	Limit   int    `json:"-" form:"limit"   binding:"omitempty,min=1,max=100"`
	Offset  int    `json:"-" form:"offset"  binding:"omitempty,min=0"`
	Sort    string `json:"-" form:"sort"    binding:"omitempty,oneof=last_messaged_at created_at"`
	Order   Order  `json:"-" form:"order"   binding:"omitempty,oneof=ASC DESC"`
	Search  string `json:"-" form:"search"`
	Include string `json:"-" form:"include"`
}
type UserListResponse struct {
	Users  []User `json:"users"`
	Total  int    `json:"total"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
	Order  Order  `json:"order"`
	Search string `json:"search,omitempty"`
}

// ====== USER LEVEL ======
// /users/:id

// # GET /users/:id
//
// Get user
// by tg id
//
// Uri:
//   - id
//
// Query:
//   - [ include ] = "" | "created_at,last_messaged_at,context,exercises?<exercises_params>" (Any combination)
//
// <exercises_params>: see [ExerciseListRequest] query params (encode!)
type UserGetRequest struct {
	TgID int64 `json:"-" uri="id" binding="required"`

	Include string `json:"-" form:"include"`
}
type UserGetResponse struct {
	User User `json:"user"`
}

// # PATCH /users/:id
//
// Patch user
// by tg id
//
// Uri:
//   - id
//
// Body: (at least 1 should be provided. nil won't change the current value)
//   - [ tg_tag ] = current | 4-32 (https://skybots.ru/telegram-limity)
//   - [ status ] = current | "unknown,requested,confirmed,restricted,banned,admin" (1, use [UserStatus] constants)
//   - [ context ] = current | Any context string
//   - [ context_data ] = current
//
// Query:
//   - [ include ] = "" | "created_at,last_messaged_at,context,exercises?<exercises_params>" (Any combination)
//
// <exercises_params>: see [ExerciseListRequest] query params (encode!)
type UserPatchRequest struct {
	TgID int64 `json:"-" uri:"id" binding="required"`

	TgTag       *string      `json:"tg_tag,omitempty"       binding="omitempty,min=3,max=100"`
	Status      *UserStatus  `json:"status,omitempty"       binding="omitempty,oneof=unknown requested confirmed restricted banned admin"` //SYNC WITH [UserStatus]
	Context     *string      `json:"context,omitempty"`
	ContextData *ContextData `json:"context_data,omitempty"`

	Include string `json:"-" form="include"`
}
type UserPatchResponse struct {
	User User `json:"user"`
}

// # DELETE /users/:id
//
// Delete user
// by tg id
//
// Uri:
//   - tg_id
type UserDeleteRequest struct {
	TgID int64 `json:"-" uri="id" binding="required"`
}
