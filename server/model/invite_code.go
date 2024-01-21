package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type InviteCodes struct {
	ID        uint `gorm:"primaryKey"`
	Code      uint16
	State     bool // 状态: false 未使用 true 已使用
	CreatedAt time.Time
	UpdatedAt time.Time
}

func InviteCode() *gorm.DB {
	return mysql.GetInstance().Model(&InviteCodes{})
}
