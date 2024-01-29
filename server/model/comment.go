package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Comments struct {
	ID                 uint `gorm:"primarykey" json:"id"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time `gorm:"index"`
	ParentId           uint       `json:"parentId"`
	RootId             uint       `json:"rootId"`
	Content            string     `json:"content"`
	UserId             uint
	BusinessId         uint `json:"articleId"`
	TenantId           uint
	ChildComments      []*Comments `gorm:"-"`
	ChildCommentNumber uint        `gorm:"-"`
}

type ChildCommentNumber struct {
	RootId uint `json:"rootId"`
	Number uint `json:"number"`
}

func Comment() *gorm.DB {
	return mysql.GetInstance().Model(&Comments{})
}
