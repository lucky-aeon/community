package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MessageTemplates struct {
	ID        int `gorm:"primarykey"`
	Content   string
	CreatedAt time.Time
	UpdateAt  time.Time
}

func MessageTemplate() *gorm.DB {
	return mysql.GetInstance().Model(&MessageTemplates{})
}
