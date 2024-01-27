package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Articles struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	Title     string
	Desc      string
	UserId    uint
	State     uint // 状态:草稿/发布/待解决/已解决/已关闭
	Like      uint
	Type      uint
	Users     Users `gorm:"-"`
}

func Article() *gorm.DB {
	return mysql.GetInstance().Model(&Articles{})
}
