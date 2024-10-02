package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/postgre"
)

type AnswerCaches struct {
	ID              int                  `gorm:"primarykey" json:"id"`
	QuestionCacheId int                  `json:"questionCacheId"`
	Answer          string               `json:"answer"`
	Type            constant.ContentType `json:"type"`
	Link            string               `json:"link"`
	Content         string               `json:"content"`
	DocumentId      int                  `json:"documentId"`
	Remark          string               `json:"remake"`
}

func AnswerCache() *gorm.DB {
	return postgre.GetInstance().Model(&AnswerCaches{})
}
