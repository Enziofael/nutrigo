package models

import (
	"time"
)

type User struct {
	ID             uint
	TgID           uint
	TgTag          string
	CreatedAt      time.Time
	LastMessagedAt time.Time
	Status         Status
}

type Status string

const (
	StatusRequested  Status = "requested"
	StatusConfirmed  Status = "confirmed"
	StatusRestricted Status = "restricted"
	StatusBanned     Status = "banned"
	StatusAdmin      Status = "admin"
)

func (s Status) Valid() bool {
	switch s {
	case StatusRequested, StatusConfirmed, StatusRestricted, StatusBanned, StatusAdmin:
		return true
	default:
		return false
	}
}
