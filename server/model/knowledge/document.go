package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/postgre"
	"xhyovo.cn/community/pkg/time"
)

type Documents struct {
	ID         int                  `gorm:"primarykey" json:"id"`
	CreatedAt  time.LocalTime       `json:"createdAt"`
	Content    string               `json:"Content,omitempty" binding:"required" msg:"原文内容不能为空"`
	Type       constant.ContentType `json:"type,omitempty"` // 假设是字符串类型
	Link       string               `json:"link,omitempty"`
	Remark     string               `json:"remark"`
	BusinessId int                  `json:"businessId"`
	Answer     string               `json:"answer" gorm:"-"`
}

func Document() *gorm.DB {
	return postgre.GetInstance().Model(&Documents{})
}
