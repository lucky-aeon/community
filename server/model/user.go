package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Users struct {
	ID         uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `gorm:"index"`
	Name       string
	Account    string
	Password   string
	InviteCode uint16
	Desc       string
	Avatar     string // todo 暂时为url
}

func User() *gorm.DB {
	return mysql.GetInstance().Model(&Users{})
}
