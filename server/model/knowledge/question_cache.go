package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/postgre"
)

type QuestionCaches struct {
	ID       int    `gorm:"primarykey" json:"id"`
	Question string `json:"question"`
}

func QuestionCache() *gorm.DB {
	return postgre.GetInstance().Model(&QuestionCaches{})
}
