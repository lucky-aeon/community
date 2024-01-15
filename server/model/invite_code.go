package model

import "time"

type InviteCode struct {
	ID        uint `gorm:"primaryKey"`
	Code      uint16
	State     bool // 状态: false 未使用 true 已使用
	CreatedAt time.Time
	UpdatedAt time.Time
}
