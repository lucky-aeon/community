package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Comments struct {
	ID            uint `gorm:"primarykey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time `gorm:"index"`
	ParentId      uint
	RootId        uint
	Content       string
	UserId        uint
	BusinessId    uint
	TenantId      uint
	ChildComments []*Comments `gorm:"-"`
}

func Comment() *gorm.DB {
	return mysql.GetInstance().Model(&Comments{})
}
