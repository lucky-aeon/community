package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type ArticleRelations struct {
	ID        int       `gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `gorm:"index"`
	ParentId  int       `json:"parentId"`
	RootId    int       `json:"rootId"`
	ArticleId int       `json:"articleId"`
}

func ArticleRelation() *gorm.DB {
	return mysql.GetInstance().Model(&ArticleRelations{})
}
