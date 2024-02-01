package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MessageLogs struct {
	ID        int `gorm:"primarykey"`
	Content   string
	From      int
	To        int
	Type      int
	CreatedAt time.Time
	DeletedAt time.Time
}

func MessageLog() *gorm.DB {
	return mysql.GetInstance().Model(&MessageLogs{})
}
