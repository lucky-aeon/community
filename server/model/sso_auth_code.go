package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type SsoAuthCode struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Code        string    `json:"code" gorm:"type:varchar(100);uniqueIndex;not null;comment:授权码"`
	AppKey      string    `json:"app_key" gorm:"type:varchar(50);not null;comment:应用标识"`
	UserId      int       `json:"user_id" gorm:"not null;comment:用户ID"`
	RedirectUrl string    `json:"redirect_url" gorm:"type:varchar(500);not null;comment:回调地址"`
	Used        bool      `json:"used" gorm:"default:false;comment:是否已使用"`
	ExpiresAt   time.Time `json:"expires_at" gorm:"comment:过期时间"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (SsoAuthCode) TableName() string {
	return "sso_auth_codes"
}

func SsoAuthCodeModel() *gorm.DB {
	return mysql.GetInstance().Model(&SsoAuthCode{})
}