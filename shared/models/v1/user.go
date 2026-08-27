package models

import (
	"fmt"
	"time"
)

func init() {
	userCreateRequest_IncludeQuery = NewIncludeQuery(UserCreateRequest{})
	userListRequest_IncludeQuery = NewIncludeQuery(UserListRequest{})
	userGetRequest_IncludeQuery = NewIncludeQuery(UserGetRequest{}, ExerciseListRequest{})
	userPatchRequest_IncludeQuery = NewIncludeQuery(UserPatchRequest{}, ExerciseListRequest{})
}

// =========================================================
// MODEL
// =========================================================

type User struct {
	ID     int64      `json:"id"  binding:"required,min=1"`
	TgTag  string     `json:"tg_tag" binding:"required,min=3,max=100"`
	Status UserStatus `json:"status" binding:"required,oneof=unknown requested confirmed restricted banned admin"` //SYNC WITH [UserStatus]

	CreatedAt      time.Time   `json:"created_at,omitempty"       binding:"omitempty,min=0"`
	LastMessagedAt time.Time   `json:"last_messaged_at,omitempty" binding:"omitempty,min=0"`
	Context        string      `json:"context,omitempty"`      //binding any
	ContextData    ContextData `json:"context_data,omitempty"` //binding any
	Exercises      []Exercise  `json:"exercises,omitempty"        binding:"omitempty,dive"`
}

type ContextData struct {
	MessageID int               `json:"message_id"         binding:"required"`
	Focus     int               `json:"focus"`            //binding any
	Values    map[string]string `json:"values,omitempty"` //binding any
}

func (cd ContextData) IsZero() bool {
	return cd.MessageID == 0 && cd.Focus == 0 && len(cd.Values) == 0
}

// =========================================================
// DTO
// =========================================================

// ====== GLOBAL LEVEL ======
// /users

// # POST /users
//
// Create user
// by their id
//
// Body:
//   - id
//   - tg_tag
//
// Query:
//   - [ include ] = "" | "created_at,last_messaged_at" (Any combination)
type UserCreateRequest struct {
	ID    int64  `json:"id"  binding:"required"`
	TgTag string `json:"tg_tag" binding:"required,min=3,max=100"`

	Include string `json:"-" form:"include"`
}
type UserCreateResponse struct {
	User    User  `json:"user"`
	Version int64 `json:"version"`
}

var userCreateRequest_IncludeQuery IncludeQuery

func (UserCreateRequest) AllowedParams() string {
	return "created_at,last_messaged_at"
}
func (UserCreateRequest) GetIncludeQuery() IncludeQuery {
	return userCreateRequest_IncludeQuery
}

// Can return errors:
//   - [ErrInvalidID]
//   - [ErrInvalidTgTag]
//   - [ErrInvalidIncludeQuery]
func (req UserCreateRequest) Validate() error {
	if req.ID <= 0 {
		return fmt.Errorf("%w: %w - must be positive. Actual: %d", ErrValidation, ErrInvalidID, req.ID)
	}
	if len(req.TgTag) < 3 || len(req.TgTag) > 100 {
		return fmt.Errorf("%w: %w - len must be in range from 1 to 100 included. Actual: %d", ErrValidation, ErrInvalidTgTag, len(req.TgTag))
	}
	if !ValidateIncludeQuery(req.Include) {
		return fmt.Errorf("%w: %w - see format at IncludeQuery doc. Actual: \"%s\"", ErrValidation, ErrInvalidIncludeQuery, req.Include)
	}
	return nil
}

// # GET /users
//
// List
// users
//
// Query:
//   - [ limit ] = 10 | 1-100
//   - [ offset ] = 0 | >= 0
//   - [ sort ] =  "relevance" | "last_messaged_at,created_at,relevance" (1) (relevance: by search similarity if search is provided, last_messaged_at otherwise)
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

// Can return errors:
//   - [ErrInvalidLimit]
//   - [ErrInvalidOffset]
//   - [ErrInvalidSort]
//   - [ErrInvalidOrder]
//   - [ErrInvalidIncludeQuery]
func (req *UserListRequest) Validate() error {
	if req.Limit < 1 || req.Limit > 100 {
		return fmt.Errorf("%w: %w - len must be in range from 1 to 100 included. Actual: %d", ErrValidation, ErrInvalidLimit, req.Limit)
	}
	if req.Offset < 0 {
		return fmt.Errorf("%w: %w - offset must be non-negative. Actual: %d", ErrValidation, ErrInvalidOffset, req.Offset)
	}
	if allowed := map[string]bool{"last_messaged_at": true, "created_at": true, "relevance": true}; !allowed[req.Sort] {
		return fmt.Errorf("%w: %w - must be one of \"last_messaged_at\", \"created_at\", \"relevance\". Actual: %s", ErrValidation, ErrInvalidSort, req.Sort)
	}
	if !req.Order.Validate() {
		return fmt.Errorf("%w: %w - must be one of \"ASC\", \"DESC\". Actual: %s", ErrValidation, ErrInvalidOrder, string(req.Order))
	}
	if !ValidateIncludeQuery(req.Include) {
		return fmt.Errorf("%w: %w - see format at IncludeQuery doc. Actual: \"%s\"", ErrValidation, ErrInvalidIncludeQuery, req.Include)
	}
	return nil
}

var userListRequest_IncludeQuery IncludeQuery

func (UserListRequest) AllowedParams() string {
	return "created_at,last_messaged_at"
}
func (req UserListRequest) GetIncludeQuery() IncludeQuery {
	return userListRequest_IncludeQuery
}

// ====== USER LEVEL ======
// /users/:id

// # GET /users/:id
//
// Get user
// by their id
//
// Uri:
//   - id
//
// Query:
//   - [ include ] = "" | "created_at,last_messaged_at,context,exercises?<exercises_params>" (Any combination)
//
// <exercises_params>: see [ExerciseListRequest] query params (encode!)
type UserGetRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`

	Include string `json:"-" form:"include"`
}
type UserGetResponse struct {
	User    User  `json:"user"`
	Version int64 `json:"version"`
}

// Can return errors:
//   - [ErrInvalidID]
//   - [ErrInvalidIncludeQuery]
func (req *UserGetRequest) Validate() error {
	if req.ID <= 0 {
		return fmt.Errorf("%w: %w - must be valid non-negative. Actual: %d", ErrValidation, ErrInvalidID, req.ID)
	}
	if !ValidateIncludeQuery(req.Include) {
		return fmt.Errorf("%w: %w - see format at IncludeQuery doc. Actual: \"%s\"", ErrValidation, ErrInvalidIncludeQuery, req.Include)
	}
	return nil
}

var userGetRequest_IncludeQuery IncludeQuery

func (UserGetRequest) AllowedParams() string {
	return "created_at,last_messaged_at,context,exercises?<exercises_params>"
}
func (req UserGetRequest) GetIncludeQuery() IncludeQuery {
	return userGetRequest_IncludeQuery
}

// # PATCH /users/:id
//
// Patch user
// by their id
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
//   - [ include ] = "" | "created_at,last_messaged_at,context,exercises<exercises_params>" (Any combination)
//
// <exercises_params>: see [ExerciseListRequest] query params (encode!)
type UserPatchRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`

	TgTag                *string      `json:"tg_tag,omitempty"           binding:"omitempty,min=3,max=100"`
	Status               *UserStatus  `json:"status,omitempty"           binding:"omitempty,oneof=unknown requested confirmed restricted banned admin"` //SYNC WITH [UserStatus]
	Context              *string      `json:"context,omitempty"`
	ContextData          *ContextData `json:"context_data,omitempty"`
	UpdateLastMessagedAt bool         `json:"update_last_messaged_at,omitempty"`

	Include string `json:"-" form:"include"`
}
type UserPatchResponse struct {
	User    User  `json:"user"`
	Version int64 `json:"version"`
}

// Can return errors:
//   - [ErrInvalidID]
//   - [ErrInvalidTgTag]
//   - [ErrInvalidUserStatus]
//   - [ErrInvalidContextData]
//   - [ErrInvalidIncludeQuery]
func (req *UserPatchRequest) Validate() error {
	if req.ID <= 0 {
		return fmt.Errorf("%w: %w - must be positive. Actual: %d", ErrValidation, ErrInvalidID, req.ID)
	}
	if req.TgTag != nil && (len(*req.TgTag) < 3 || len(*req.TgTag) > 100) {
		return fmt.Errorf("%w: %w - len must be in range from 1 to 100 included. Actual: %d", ErrValidation, ErrInvalidTgTag, len(*req.TgTag))
	}
	if req.Status != nil && !req.Status.Validate() {
		return fmt.Errorf("%w: %w - must be one of \"unknown\", \"requested\", \"confirmed\", \"restricted\", \"banned\", \"admin\". Actual: %s", ErrValidation, ErrInvalidUserStatus, string(*req.Status))
	}
	if req.ContextData != nil && req.ContextData.MessageID <= 0 {
		return fmt.Errorf("%w: %w - must be valid non-negative. Actual: %d", ErrValidation, ErrInvalidContextData, req.ContextData.MessageID)
	}
	if !ValidateIncludeQuery(req.Include) {
		return fmt.Errorf("%w: %w - see format at IncludeQuery doc. Actual: \"%s\"", ErrValidation, ErrInvalidIncludeQuery, req.Include)
	}

	return nil
}

var userPatchRequest_IncludeQuery IncludeQuery

func (UserPatchRequest) AllowedParams() string {
	return "created_at,last_messaged_at,context,exercises<exercises_params>"
}
func (UserPatchRequest) GetIncludeQuery() IncludeQuery {
	return userPatchRequest_IncludeQuery
}

// # DELETE /users/:id
//
// Delete user
// by their id
//
// Uri:
//   - id
type UserDeleteRequest struct {
	ID int64 `json:"-" uri:"id" binding:"required"`
}

// Can return errors:
//   - [ErrInvalidID]
func (req *UserDeleteRequest) Validate() error {
	if req.ID <= 0 {
		return fmt.Errorf("%w: %w - must be positive. Actual: %d", ErrValidation, ErrInvalidID, req.ID)
	}
	return nil
}
