package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Drafts struct {
	ID        int            `gorm:"primaryKey"`
	Content   string         `json:"content"`
	Type      int            `json:"type"`
	LabelIds  string         `json:"labelIds"`
	Labels    []int          `json:"labels" gorm:"-"`
	UserId    int            `json:"userId"`
	ArticleId int            `json:"articleId"`
	State     int            `json:"state"` // 临时保存文章状态： 编辑 ，发布
	CreatedAt time.LocalTime `json:"createdAt"`
}

func Draft() *gorm.DB {
	return mysql.GetInstance().Model(&Drafts{})
}
