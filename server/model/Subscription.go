package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Subscriptions struct {
	ID         uint `gorm:"primaryKey"`
	userId     uint // 订阅人
	event      uint // 事件类型： 人 / 文章
	businessId uint // 业务id  人id / 文章id
}

func Subscription() *gorm.DB {
	return mysql.GetInstance().Model(&Subscriptions{})
}
