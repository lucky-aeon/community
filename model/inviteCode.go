package model

import "time"

type InviteCode struct {
	ID        uint
	Code      int
	State     bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
