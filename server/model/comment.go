package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Comments struct {
	gorm.Model
	ParentId   uint
	Content    string
	UserId     uint
	BusinessId uint
	TenantId   uint
}

func Comment() *gorm.DB {
	return mysql.GetInstance().Model(&Comments{})
}
