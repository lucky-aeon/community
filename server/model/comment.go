package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/pkg/mysql"
)

type Comments struct {
	ID                 int            `gorm:"primarykey" json:"id"`
	CreatedAt          time.LocalTime `json:"createdAt"`
	UpdatedAt          time.LocalTime `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt `gorm:"index"`
	ParentId           int            `json:"parentId"`
	RootId             int            `json:"rootId"`
	Content            string         `json:"content" binding:"required" msg:"请评论内容"`
	FromUserId         int            `json:"FromUserId"`
	ToUserId           int            `json:"toUserId"`
	BusinessId         int            `json:"businessId" binding:"required" msg:"评论对象不可为空"`
	BusinessUserId     int            `json:"businessUserId"`
	TenantId           int            `json:"tenantId"`
	ChildComments      []*Comments    `gorm:"-" json:"childComments"`
	ChildCommentNumber int            `gorm:"-" json:"childCommentNumber"`
	FromUserName       string         `json:"fromUserName" gorm:"-"`
	ToUserName         string         `json:"toUserName" gorm:"-"`
	ArticleTitle       string         `json:"articleTitle" gorm:"-"`
	FromUserAvatar     string         `json:"fromUserAvatar" gorm:"-"`
	ToUserAvatar       string         `json:"toUserAvatar" gorm:"-"`
	AdoptionState      bool           `json:"adoptionState" gorm:"-"`
	Reactions          []ReactionSummary `json:"reactions" gorm:"-"`
}

type ChildCommentNumber struct {
	RootId int `json:"rootId"`
	Number int `json:"number"`
}

func Comment() *gorm.DB {
	return mysql.GetInstance().Model(&Comments{})
}
