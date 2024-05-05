package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/pkg/mysql"
)

type MessageLogs struct {
	ID        int            `gorm:"primarykey" json:"id"`
	Content   string         `json:"content"`
	From      int            `json:"from"`
	To        int            `json:"to"`
	Type      int            `json:"type"`
	ArticleId int            `json:"articleId"`
	EventId   int            `json:"eventId"`
	CreatedAt time.LocalTime `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt"`
}

func MessageLog() *gorm.DB {
	return mysql.GetInstance().Model(&MessageLogs{})
}
