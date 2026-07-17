package v1

import "time"

type User struct {
	ID             uint      `json:"id"`
	TgID           uint      `json:"tgID"`
	TgTag          string    `json:"tgTag"`
	CreatedAt      time.Time `json:"createdAt"`
	LastMessagedAt time.Time `json:"lastMessagedAt"`
	Status         string    `json:"status"`
}

type UserCreateRequest struct {
	TgID  uint   `json:"tgId" binding:"required"`
	TgTag string `json:"tgTag" binding:"required"`
}

const (
	StatusRequested  = "requested"
	StatusConfirmed  = "confirmed"
	StatusRestricted = "restricted"
	StatusBanned     = "banned"
	StatusAdmin      = "admin"
)
