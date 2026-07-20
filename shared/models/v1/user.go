package v1

import "time"

type User struct {
	ID             int64     `json:"id"`
	TgID           int64     `json:"tgID"`
	TgTag          string    `json:"tgTag"`
	CreatedAt      time.Time `json:"createdAt"`
	LastMessagedAt time.Time `json:"lastMessagedAt"`
	Status         string    `json:"status"`
}

type UserCreateRequest struct {
	TgID  int64  `json:"tgId" binding:"required"`
	TgTag string `json:"tgTag" binding:"required"`
}

const (
	StatusRequested  = "requested"
	StatusConfirmed  = "confirmed"
	StatusRestricted = "restricted"
	StatusBanned     = "banned"
	StatusAdmin      = "admin"
)
