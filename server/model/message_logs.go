package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MessageLogs struct {
	ID        int       `gorm:"primarykey" json:"id"`
	Content   string    `json:"content"`
	From      int       `json:"from"`
	To        int       `json:"to"`
	Type      int       `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	DeletedAt time.Time `json:"deletedAt"`
}

func MessageLog() *gorm.DB {
	return mysql.GetInstance().Model(&MessageLogs{})
}
