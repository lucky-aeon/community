package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

// it's issue or answer
type Articles struct {
	gorm.Model
	Title  string
	Desc   string
	UserId uint
	State  uint // 状态:草稿/发布/待解决/已解决/已关闭
	Like   uint
	Type   uint
	Users  Users `gorm:"-"`
}

func Article() *gorm.DB {
	return mysql.GetInstance().Model(&Articles{})
}
