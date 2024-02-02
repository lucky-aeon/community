package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Types struct {
	ID            int `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time `gorm:"index"`
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
