package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Types struct {
	gorm.Model
	Title        string
	Desc         string
	State        uint
	Sort         int
	ArticleState string // 分类下文章的状态
}

func Type() *gorm.DB {
	return mysql.GetInstance().Model(&Types{})
}
