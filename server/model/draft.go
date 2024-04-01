package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Drafts struct {
	ID        int            `gorm:"primaryKey"`
	Content   string         `json:"content"`
	UserId    int            `json:"userId"`
	CreatedAt time.LocalTime `json:"createdAt"`
}

func Draft() *gorm.DB {
	return mysql.GetInstance().Model(&Drafts{})
}
