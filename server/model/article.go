package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Articles struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `gorm:"index"`
	Title     string
	Desc      string
	UserId    int
	State     int // 状态:草稿/发布/待解决/已解决/已关闭
	Like      int
	Type      int
	Users     Users `gorm:"-"`
}

func Article() *gorm.DB {
	return mysql.GetInstance().Model(&Articles{})
}
