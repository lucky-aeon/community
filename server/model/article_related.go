package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type ArticleRelations struct {
	ID        int            `gorm:"primarykey"`
	CreatedAt time.LocalTime `json:"createdAt"`
	UpdatedAt time.LocalTime `json:"updatedAt"`
	DeletedAt time.LocalTime `gorm:"index"`
	ParentId  int            `json:"parentId"`
	RootId    int            `json:"rootId"`
	ArticleId int            `json:"articleId"`
}

func ArticleRelation() *gorm.DB {
	return mysql.GetInstance().Model(&ArticleRelations{})
}
