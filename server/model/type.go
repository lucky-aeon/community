package model

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Types struct {
	ID            int       `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	DeletedAt     time.Time `gorm:"index"`
	ParentId      int       `json:"parentId"`
	Title         string    `json:"title"`
	Desc          string    `json:"desc"`
	State         int       `json:"state"`
	Sort          int       `json:"sort"`
	ArticleState  string    `json:"articleState"` // 分类下文章的状态
	ArticleStates []string  `gorm:"-" json:"articleStates"`
	FlagName      string
}

type TypeSimple struct {
	TypeFlag  string `json:"flag" gorm:"column:flag_name"` // flag name
	TypeTitle string `json:"title" gorm:"column:title"`    // title
}

func Type() *gorm.DB {
	return mysql.GetInstance().Model(&Types{})
}
