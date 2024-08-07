package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/pkg/mysql"
)

type Subscriptions struct {
	ID           int `gorm:"primaryKey" json:"id"`
	SubscriberId int `json:"userId"`                                              // 订阅人
	SendId       int `json:"sendId"`                                              // 发送人   发送人给订阅人发消息 系统默认的是 13
	EventId      int `json:"eventId"`                                             // 事件类型： 人 / 文章
	BusinessId   int `json:"businessId"  binding:"required" msg:"采纳对应业务 id 不能未空"` // 业务id  人id / 文章id
	IndexKey     string
	CreatedAt    time.LocalTime `json:"createdAt"`
	EventName    string         `json:"eventName" gorm:"-"`
	BusinessName string         `json:"businessName" gorm:"-"` // 根据事件类型选择对应的 文章标题 / 用户昵称
}

type SubscriptionState struct {
	SubscriberId int
	EventId      int `json:"eventId"`
	BusinessId   int `json:"businessId"`
}

func Subscription() *gorm.DB {
	return mysql.GetInstance().Model(&Subscriptions{})
}
