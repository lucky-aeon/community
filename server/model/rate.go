package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Rates struct {
	ID        int            `gorm:"primarykey" json:"id"`
	UserId    int            `json:"userId"`
	Content   string         `json:"content"`
	Avatar    string         `json:"avatar"`
	CreatedAt time.LocalTime `json:"createdAt"`
	UpdatedAt time.LocalTime `json:"updateAt"`
	Nickname  string         `json:"nickName" gorm:"-"`
}

func Rate() *gorm.DB {
	return mysql.GetInstance().Model(&Rates{})
}
