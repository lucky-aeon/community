package model

import "time"

type User struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
	Name       string
	Account    string
	Password   string
	InviteCode int
}
