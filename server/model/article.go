package model

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Articles struct {
	ID        int        `gorm:"primarykey" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index,omitempty"`
	Title     string     `json:"title"`
	Desc      string     `json:"desc,omitempty"`
	UserId    int        `json:"userId,omitempty"`
	State     int        `json:"state"` // 状态:草稿/发布/待解决/已解决/已关闭
	Like      int        `json:"like"`
	Type      int        `json:"type"`
	Tags      string     `json:"tags"`
	Users     `gorm:"-" json:"user"`
}

type ArticleData struct {
	ID         int    `gorm:"primarykey" json:"id"`
	Title      string `json:"title"`
	State      int    `json:"state"` // 状态:草稿/发布/待解决/已解决/已关闭
	Like       int    `json:"like"`
	Tags       any    `json:"tags"`
	TypeSimple `json:"type"`
	UserSimple `json:"user"`
	Desc       string    `json:"desc,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func Article() *gorm.DB {
	return mysql.GetInstance().Model(&Articles{})
}
