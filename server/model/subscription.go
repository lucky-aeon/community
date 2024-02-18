package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Subscriptions struct {
	ID           int `gorm:"primaryKey" json:"id"`
	SubscriberId int `json:"userId"`                                          // 订阅人
	SendId       int `json:"sendId"`                                          // 发送人
	EventId      int `json:"eventId"`                                         // 事件类型： 人 / 文章
	BusinessId   int `json:"businessId" binding:"required" msg:"请选择对应业务进行订阅"` // 业务id  人id / 文章id
	IndexKey     string
	CreatedAt    time.Time `json:"createdAt"`
	EventName    string    `json:"eventName" gorm:"-"`
	BusinessName string    `json:"businessName" gorm:"-"` // 根据事件类型选择对应的 文章标题 / 用户昵称
}

func Subscription() *gorm.DB {
	return mysql.GetInstance().Model(&Subscriptions{})
}
