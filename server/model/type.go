package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Types struct {
	gorm.Model
	Title string
	Desc  string
	State uint
	Sort  int
}

func Type() *gorm.DB {
	return mysql.GetInstance().Model(&Types{})
}
