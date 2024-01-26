package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type ArticleRelations struct {
	gorm.Model
	ParentId  int
	RootId    int
	ArticleId int
}

func ArticleRelation() *gorm.DB {
	return mysql.GetInstance().Model(&ArticleRelations{})
}
