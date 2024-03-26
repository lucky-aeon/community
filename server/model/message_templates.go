package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/pkg/mysql"
)

type MessageTemplates struct {
	ID        int            `gorm:"primarykey" json:"id"`
	Content   string         `json:"content" binding:"required" msg:"内容不能未空"`
	EventId   int            `json:"eventId"` // 事件id
	CreatedAt time.LocalTime `json:"createdAt"`
	UpdatedAt time.LocalTime `json:"updateAt"`
	EventName string         `json:"eventName" gorm:"-"`
}

func MessageTemplate() *gorm.DB {
	return mysql.GetInstance().Model(&MessageTemplates{})
}
