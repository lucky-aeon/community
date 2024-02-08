package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type InviteCodes struct {
	ID         int       `gorm:"primaryKey"`
	MemberId   int       `json:"memberId"`
	Code       int       `json:"code"`
	State      bool      `json:"state"` // 状态: false 未使用 true 已使用
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	MemberName string    `json:"memberName" gorm:"-"`
}

type GenerateCode struct {
	Number   int `json:"number"`   // 生成的数量
	MemberId int `json:"memberId"` // 绑定的会员等级
}

func InviteCode() *gorm.DB {
	return mysql.GetInstance().Model(&InviteCodes{})
}
