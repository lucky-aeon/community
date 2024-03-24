package model

import (
	"xhyovo.cn/community/pkg/time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Articles struct {
	ID        int            `gorm:"primarykey" json:"id"`
	CreatedAt time.LocalTime `json:"createdAt"`
	UpdatedAt time.LocalTime `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"index,omitempty"`
	Title     string         `json:"title" binding:"required" msg:"标题不能未空"`
	Content   string         `json:"content,omitempty" binding:"required" msg:"描述不能未空"`
	UserId    int            `json:"userId,omitempty"`
	State     int            `json:"state"` // 状态:草稿/发布/待解决/已解决/私密提问
	Like      int            `json:"like"`
	Type      int            `json:"type"`
	Tags      []int          `json:"tags" gorm:"-"`
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
	Desc       string         `json:"content,omitempty"`
	CreatedAt  time.LocalTime `json:"createdAt"`
	UpdatedAt  time.LocalTime `json:"updatedAt"`
	StateName  string         `json:"stateName"`
}

func Article() *gorm.DB {
	return mysql.GetInstance().Model(&Articles{})
}
