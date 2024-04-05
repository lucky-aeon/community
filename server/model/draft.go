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
	State     int            `json:"state"` // 是否存在临时文本,1:存在，2：不存在
	CreatedAt time.LocalTime `json:"createdAt"`
}

func Draft() *gorm.DB {
	return mysql.GetInstance().Model(&Drafts{})
}
