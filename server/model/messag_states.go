package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MessageStates struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Content   string    `json:"content"`
	From      uint      `json:"from"`
	To        uint      `json:"to"`
	CreatedAt time.Time `json:"createdAt"`
}

func MessageState() *gorm.DB {
	return mysql.GetInstance().Model(&MessageStates{})
}
