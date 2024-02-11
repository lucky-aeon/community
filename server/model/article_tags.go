package model

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

// 文章标签
type ArticleTags struct {
	Id          int        `json:"id"`
	TagName     string     `json:"tag"`
	Description string     `json:"description"`
	UserId      int        `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index"`
}
type ArticleTagSimple struct {
	TagId          int    `json:"id" gorm:"column:id"`
	TagName        string `json:"name"`
	TagDescription string `json:"description" gorm:"column:description"`
}

// 标签与文章关联
type ArticleTagRelations struct {
	ArticleId int `json:"article_id"`
	TagId     int `json:"tag_id"`
}
type ArticleTagUserRelations struct {
	UserId int `json:"user_id"`
	TagId  int `json:"tag_id"`
}

func ArticleTag() *gorm.DB {
	return mysql.GetInstance().Model(&ArticleTags{})
}

func ArticleTagRelation() *gorm.DB {
	return mysql.GetInstance().Model(&ArticleTagRelations{})
}

func ArticleTagUserRelation() *gorm.DB {
	return mysql.GetInstance().Model(&ArticleTagUserRelations{})
}
