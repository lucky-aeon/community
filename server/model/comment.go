package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Comments struct {
	ID                 int `gorm:"primarykey" json:"id"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time `gorm:"index"`
	ParentId           int        `json:"parentId"`
	RootId             int        `json:"rootId"`
	Content            string     `json:"content"`
	UserId             int
	BusinessId         int `json:"articleId"`
	TenantId           int
	ChildComments      []*Comments `gorm:"-"`
	ChildCommentNumber int         `gorm:"-"`
}

type ChildCommentNumber struct {
	RootId int `json:"rootId"`
	Number int `json:"number"`
}

func Comment() *gorm.DB {
	return mysql.GetInstance().Model(&Comments{})
}
