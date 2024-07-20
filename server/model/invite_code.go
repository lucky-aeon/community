package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/pkg/mysql"
)

type InviteCodes struct {
	ID              int            `gorm:"primaryKey" json:"id"`
	MemberId        int            `json:"memberId"`
	Code            string         `json:"code"`
	State           bool           `json:"state"` // 状态: false 未使用 true 已使用
	CreatedAt       time.LocalTime `json:"createdAt"`
	UpdatedAt       time.LocalTime `json:"updatedAt"`
	MemberName      string         `json:"memberName" gorm:"-"`
	AcquisitionType int            `json:"acquisitionType"` // 获取类型:1 购买，2赠予
	Creator         int            `json:"creator"`
}

type GenerateCode struct {
	Number          int `json:"number"`   // 生成的数量
	MemberId        int `json:"memberId"` // 绑定的会员等级
	Creator         int `json:"creator"`
	AcquisitionType int `json:"acquisitionType"` // 获取类型:1 购买，2赠予
}

func InviteCode() *gorm.DB {
	return mysql.GetInstance().Model(&InviteCodes{})
}
