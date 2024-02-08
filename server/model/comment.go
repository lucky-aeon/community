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
	Content            string      `json:"content" binding:"required" msg:"请评论内容"`
	FromUserId         int         `json:"FromUserId"`
	ToUserId           int         `json:"toUserId"`
	BusinessId         int         `json:"articleId" binding:"required" msg:"请选择对应的文章进行评论"`
	TenantId           int         `json:"tenantId"`
	ChildComments      []*Comments `gorm:"-" json:"childComments"`
	ChildCommentNumber int         `gorm:"-" json:"childCommentNumber"`
	FromUserName       string      `json:"fromUserName" gorm:"-"`
	ToUserName         string      `json:"toUserName" gorm:"-"`
	ArticleTitle       string      `json:"articleTitle" gorm:"-"`
}

type ChildCommentNumber struct {
	RootId int `json:"rootId"`
	Number int `json:"number"`
}

func Comment() *gorm.DB {
	return mysql.GetInstance().Model(&Comments{})
}
