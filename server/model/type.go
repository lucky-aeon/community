package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Types struct {
	gorm.Model
	ParentId      int
	Title         string
	Desc          string
	State         int
	Sort          int
	ArticleState  string   // 分类下文章的状态
	ArticleStates []string `gorm:"-"`
}

func Type() *gorm.DB {
	return mysql.GetInstance().Model(&Types{})
}
