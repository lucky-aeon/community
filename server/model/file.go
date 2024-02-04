package model

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Files struct {
	ID         int       `gorm:"primaryKey"`
	FileKey    string    `json:"fileKey"`
	Size       int64     `json:"size"`
	Format     string    `json:"format"`
	UserId     int       `json:"userId"`
	BusinessId int       `json:"businessId"`
	TenantId   int       `json:"tenantId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func File() *gorm.DB {
	return mysql.GetInstance().Model(&Files{})
}
