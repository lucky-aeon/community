package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MessageStates struct {
	ID        uint `gorm:"primarykey"`
	Content   string
	From      uint
	To        uint
	CreatedAt time.Time
}

func MessageState() *gorm.DB {
	return mysql.GetInstance().Model(&MessageStates{})
}
