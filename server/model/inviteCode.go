package model

import "time"

type InviteCode struct {
	Id        uint // id
	Code      int
	State     bool // 状态: false 未使用 true 已使用
	CreatedAt time.Time
	UpdatedAt time.Time
}
