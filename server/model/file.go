package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Files struct {
	ID         int            `gorm:"primaryKey"`
	FileKey    string         `json:"fileKey" from:"filename"`
	Size       int64          `json:"size" from:"size"`
	Format     string         `json:"mimeType" from:"mimeType"`
	UserId     int            `json:"userId" from:"userId"`
	BusinessId int            `json:"businessId" from:"articleId"` // todo 遗弃
	TenantId   int            `json:"tenantId"`
	CreatedAt  time.LocalTime `json:"createdAt"`
	UpdatedAt  time.LocalTime `json:"updatedAt"`
	UserName   string         `json:"userName" gorm:"-"`
	SizeName   string         `json:"sizeName"`
}

func File() *gorm.DB {
	return mysql.GetInstance().Model(&Files{})
}
