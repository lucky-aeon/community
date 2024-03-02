package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// 文章标签
type ArticleTags struct {
	Id          int             `json:"id"`
	TagName     string          `json:"tag" binding:"required" msg:"标签不能未空"`
	Description string          `json:"desc"`
	UserId      int             `json:"user_id"`
	CreatedAt   time.LocalTime  `json:"created_at"`
	UpdatedAt   time.LocalTime  `json:"updated_at"`
	DeletedAt   *time.LocalTime `gorm:"index"`
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

type TagArticleCount struct {
	ArticleCount int
	TagId        int
	TagName      string
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
