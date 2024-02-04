package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MessageTemplates struct {
	ID        int       `gorm:"primarykey" json:"id"`
	Content   string    `json:"content"`
	Event     int       `json:"event"` // 事件id
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`
}

func MessageTemplate() *gorm.DB {
	return mysql.GetInstance().Model(&MessageTemplates{})
}
