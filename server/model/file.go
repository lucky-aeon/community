package model

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Files struct {
	ID         int       `gorm:"primaryKey"`
	FileKey    string    `json:"filename" from:"filename"`
	Size       int64     `json:"size" from:"size"`
	Format     string    `json:"mimeType" from:"mimeType"`
	UserId     int       `json:"userId" from:"userId"`
	BusinessId int       `json:"businessId" from:"articleId"`
	TenantId   int       `json:"tenantId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func File() *gorm.DB {
	return mysql.GetInstance().Model(&Files{})
}
