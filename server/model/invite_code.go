package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type InviteCodes struct {
	ID        int       `gorm:"primaryKey"`
	Code      uint16    `json:"code"`
	State     bool      `json:"state"` // 状态: false 未使用 true 已使用
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func InviteCode() *gorm.DB {
	return mysql.GetInstance().Model(&InviteCodes{})
}
