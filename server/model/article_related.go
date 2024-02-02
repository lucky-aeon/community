package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type ArticleRelations struct {
	ID        int `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `gorm:"index"`
	ParentId  int
	RootId    int
	ArticleId int
}

func ArticleRelation() *gorm.DB {
	return mysql.GetInstance().Model(&ArticleRelations{})
}
