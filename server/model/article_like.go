package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Article_Likes struct {
	ID        int `gorm:"primarykey" json:"id"`
	ArticleId int `json:"articleId"`
	UserId    int `json:"userId"`
}

func ArticleLike() *gorm.DB {
	return mysql.GetInstance().Model(&Article_Likes{})
}
