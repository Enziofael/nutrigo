package models

import (
	"time"
)

// =========================================================
// MODEL
// =========================================================

type User struct {
	TgID   int64      `json:"tg_id"`
	TgTag  string     `json:"tg_tag"`
	Status UserStatus `json:"status"`

	CreatedAt      time.Time   `json:"created_at,omitempty"`
	LastMessagedAt time.Time   `json:"last_messaged_at,omitempty"`
	Context        string      `json:"context,omitempty"`
	ContextData    ContextData `json:"context_data,omitempty"`
	Exercises      []Exercise  `json:"exercises,omitempty"`
}

type ContextData struct {
	MessageID int               `json:"message_id"`
	Focus     int               `json:"focus"`
	Values    map[string]string `json:"values,omitempty"`
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
	TgID  int64  `json:"tg_id"`
	TgTag string `json:"tg_tag"`

	Include string `json:"-"`
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
//   - [ offset ] = 0 | >= 0
//   - [ limit ] = 10 | 0 - 100
//   - [ sort ] = "last_messaged_at" | "last_messaged_at,created_at" (1)
//   - [ order ] = "DESC" | "ASC,DESC" (1)
//   - [ search ] = "" | Any search string (encode!)
//   - [ include ] = "" | "created_at,last_messaged_at,context" (Any combination)
type UserListRequest struct {
	Offset  int    `json:"-"`
	Limit   int    `json:"-"`
	Sort    string `json:"-"`
	Order   Order  `json:"-"`
	Search  string `json:"-"`
	Include string `json:"-"`
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
	TgID int64 `json:"-"`

	Include string `json:"-"`
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
	TgID int64 `json:"-"`

	TgTag       *string      `json:"tg_tag,omitempty"`
	Status      *UserStatus  `json:"status,omitempty"`
	Context     *string      `json:"context,omitempty"`
	ContextData *ContextData `json:"context_data,omitempty"`

	Include string `json:"-"`
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
	TgID int64 `json:"-"`
}
