package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Files struct {
	ID         uint `gorm:"primaryKey"`
	FileKey    string
	Size       int64
	Format     string
	UserId     uint
	BusinessId uint
	TenantId   uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func File() *gorm.DB {
	return mysql.GetInstance().Model(&Files{})
}
