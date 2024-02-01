package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MessageStates struct {
	ID        int       `gorm:"primarykey" json:"id"`
	Content   string    `json:"content"`
	From      int       `json:"from"`
	To        int       `json:"to"`
	CreatedAt time.Time `json:"createdAt"`
}

func MessageState() *gorm.DB {
	return mysql.GetInstance().Model(&MessageStates{})
}
