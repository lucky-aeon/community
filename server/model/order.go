package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Orders struct {
	ID              int            `gorm:"primaryKey" json:"id"`
	InviteCode      string         `json:"inviteCode"`
	Price           int            `json:"price"`
	Purchaser       int            `json:"purchaser"`       // 购买者
	Creator         int            `json:"creator"`         // 订单创建者，和邀请码创建人同步
	AcquisitionType int            `json:"acquisitionType"` // 获取类型:1 赠予，2 购买
	CreatedAt       time.LocalTime `json:"createdAt"`
	PurchaserName   string         `json:"purchaserName" gorm:"-"`
	CreatorName     string         `json:"creatorName" gorm:"-"`
}

func Order() *gorm.DB {
	return mysql.GetInstance().Model(&Orders{})
}
