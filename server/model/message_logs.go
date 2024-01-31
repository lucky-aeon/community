package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MessageLogs struct {
	ID        uint `gorm:"primarykey"`
	Content   string
	From      uint
	To        uint
	Type      uint
	CreatedAt time.Time
	DeletedAt time.Time
}

func MessageLog() *gorm.DB {
	return mysql.GetInstance().Model(&MessageLogs{})
}
