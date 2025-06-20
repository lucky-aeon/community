package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type SsoApplication struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null;comment:应用名称"`
	AppKey       string    `json:"app_key" gorm:"type:varchar(50);uniqueIndex;not null;comment:应用标识"`
	AppSecret    string    `json:"app_secret" gorm:"type:varchar(100);not null;comment:应用密钥"`
	RedirectUrls string    `json:"redirect_urls" gorm:"type:text;comment:允许的回调地址，多个用逗号分隔"`
	Status       int8      `json:"status" gorm:"type:tinyint;default:1;comment:状态：1启用，0禁用"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (SsoApplication) TableName() string {
	return "sso_applications"
}

func SsoApp() *gorm.DB {
	return mysql.GetInstance().Model(&SsoApplication{})
}