package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Articles struct {
	ID        int        `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `gorm:"index"`
	Title     string     `json:"title"`
	Desc      string     `json:"desc"`
	UserId    int        `json:"userId"`
	State     int        `json:"state"` // 状态:草稿/发布/待解决/已解决/已关闭
	Like      int        `json:"like"`
	Type      int        `json:"type"`
	Users     Users      `gorm:"-" json:"users"`
}

func Article() *gorm.DB {
	return mysql.GetInstance().Model(&Articles{})
}
