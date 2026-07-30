package models

import "time"

type User struct {
	ID             int64       `json:"id"`
	TgID           int64       `json:"tg_id"`
	TgTag          string      `json:"tg_tag"`
	CreatedAt      time.Time   `json:"created_at"`
	LastMessagedAt time.Time   `json:"last_messaged_at"`
	Status         string      `json:"status"`
	Context        string      `json:"context"`
	ContextData    ContextData `json:"context_data"`
}

type ContextData struct {
	MessageID int               `json:"message_id"`
	Focus    int               `json:"focus"`
	Values    map[string]string `json:"values"`
}

type UserCreateRequest struct {
	TgID  int64  `json:"tg_id" binding:"required"`
	TgTag string `json:"name" binding:"required"`
}

type UserStatusUpdateRequest struct {
	Status string `json:"status" binding:"required,ne="`
}

type UserPatchRequest struct {
	TgTag          *string      `json:"tg_tag"`
	LastMessagedAt *time.Time   `json:"last_messaged_at"`
	Status         *string      `json:"status"`
	Context        *string      `json:"context_id"`
	ContextData    *ContextData `json:"context_data"`
}

const (
	StatusRequested  = "requested"
	StatusConfirmed  = "confirmed"
	StatusRestricted = "restricted"
	StatusBanned     = "banned"
	StatusAdmin      = "admin"
)

func ValidateStatus(status string) bool {
	switch status {
	case StatusRequested, StatusConfirmed, StatusRestricted, StatusBanned, StatusAdmin:
		return true
	default:
		return false
	}
}
