package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type Files struct {
	ID         int `gorm:"primaryKey"`
	FileKey    string
	Size       int64
	Format     string
	UserId     int
	BusinessId int
	TenantId   int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func File() *gorm.DB {
	return mysql.GetInstance().Model(&Files{})
}
