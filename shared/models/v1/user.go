package models

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type User struct {
	TgID           int64       `json:"tg_id"              binding:"required,min=1" db:"tg_id"`
	TgTag          string      `json:"tg_tag"             binding:"required,min=1"  db:"tg_tag"`
	CreatedAt      time.Time   `json:"created_at"         binding:"required"       db:"created_at"`
	LastMessagedAt time.Time   `json:"last_messaged_at"   binding:"required"       db:"last_messaged_at"`
	Status         string      `json:"status"             binding:"required"       db:"status"`
	Context        string      `json:"context"            binding:"omitempty"      db:"context"`
	ContextData    ContextData `json:"context_data"       binding:"omitempty"      db:"context_data"`
}

func SanitizeTgTag(tag string) (string, error) {
	clean := strings.TrimSpace(tag)
	if clean == "" {
		return "", fmt.Errorf("tg_tag can't be empty")
	}
	return clean, nil
}

func SanitizeStatus(status string) (string, error) {
	switch status {
	case StatusRequested, StatusConfirmed, StatusRestricted, StatusBanned, StatusAdmin:
		return status, nil
	default:
		return "", fmt.Errorf("Invalid status: %s", status)
	}
}

type ContextData struct {
	MessageID int               `json:"message_id"`
	Focus     int               `json:"focus"`
	Values    map[string]string `json:"values"`
}

// =========================================================
// CRUD REQUESTS
// =========================================================

// Create by id
type UserCreateRequest struct {
	//body
	TgID  int64  `json:"tg_id"  binding:"required"`
	TgTag string `json:"tg_tag" binding:"required"`
}

func (req UserCreateRequest) Sanitize() (UserCreateRequest, error) {
	var errTgID, errTgTag error

	req.TgID, errTgID = SanitizeTgID(req.TgID)
	req.TgTag, errTgTag = SanitizeTgTag(req.TgTag)

	err := errors.Join(errTgID, errTgTag)
	return req, err
}

// Delete by id
type UserDeleteRequest struct {
	//url
	TgID int64 `json:"tg_id"  binding:"required"`
}

func (req UserDeleteRequest) Sanitize() (UserDeleteRequest, error) {
	var errTgID error

	req.TgID, errTgID = SanitizeTgID(req.TgID)

	err := errors.Join(errTgID)
	return req, err
}

// Get by id
type UserGetRequest struct {
	//url
	TgID int64 `json:"tg_id"  binding:"required"`
}

func (req UserGetRequest) Sanitize() (UserGetRequest, error) {
	var errTgID error

	req.TgID, errTgID = SanitizeTgID(req.TgID)

	err := errors.Join(errTgID)
	return req, err
}

// Patch by id
type UserPatchRequest struct {
	//url
	TgID int64 `json:"tg_id" binding:"required"`

	//body optional (at least 1)
	TgTag       *string      `json:"tg_tag"`
	Status      *string      `json:"status"`
	Context     *string      `json:"context_id"`
	ContextData *ContextData `json:"context_data"`
}

func (req UserPatchRequest) Sanitize() (UserPatchRequest, error) {
	var errTgID, errTgTag, errStatus, errLeast error

	req.TgID, errTgID = SanitizeTgID(req.TgID)
	if req.TgTag != nil {
		var t string
		t, errTgTag = SanitizeTgTag(*req.TgTag)
		req.TgTag = &t
	}
	if req.Status != nil {
		var s string
		s, errStatus = SanitizeStatus(*req.Status)
		req.Status = &s
	}
	if req.TgTag == nil && req.Status == nil && req.Context == nil && req.ContextData == nil {
		errLeast = errors.New("At least 1 optional field should be provided")
	}

	err := errors.Join(errTgID, errTgTag, errStatus, errLeast)
	return req, err
}
