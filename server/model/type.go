package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Types struct {
	ID            int            `gorm:"primarykey" json:"id"`
	CreatedAt     time.LocalTime `json:"createdAt"`
	UpdatedAt     time.LocalTime `json:"updatedAt"`
	DeletedAt     time.LocalTime `gorm:"index"`
	ParentId      int            `json:"parentId"`
	Title         string         `json:"title"`
	Desc          string         `json:"desc"`
	State         int            `json:"state"`
	Sort          int            `json:"sort"`
	ArticleState  string         `json:"articleState"` // 分类下文章的状态
	ArticleStates []string       `gorm:"-" json:"articleStates"`
	FlagName      string
}

type TypeSimple struct {
	TypeId    int    `json:"id" gorm:"column:id"`
	TypeFlag  string `json:"flag" gorm:"column:flag_name"` // flag name
	TypeTitle string `json:"title" gorm:"column:title"`    // title
}

func Type() *gorm.DB {
	return mysql.GetInstance().Model(&Types{})
}
