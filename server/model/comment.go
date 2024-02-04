package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Comments struct {
	ID                 int         `gorm:"primarykey" json:"id"`
	CreatedAt          time.Time   `json:"createdAt"`
	UpdatedAt          time.Time   `json:"updatedAt"`
	DeletedAt          *time.Time  `gorm:"index"`
	ParentId           int         `json:"parentId"`
	RootId             int         `json:"rootId"`
	Content            string      `json:"content"`
	UserId             int         `json:"userId"`
	BusinessId         int         `json:"articleId"`
	TenantId           int         `json:"tenantId"`
	ChildComments      []*Comments `gorm:"-" json:"childComments"`
	ChildCommentNumber int         `gorm:"-" json:"childCommentNumber"`
}

type ChildCommentNumber struct {
	RootId int `json:"rootId"`
	Number int `json:"number"`
}

func Comment() *gorm.DB {
	return mysql.GetInstance().Model(&Comments{})
}
