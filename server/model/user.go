package model

import (
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
)

type Users struct {
	ID         int        `json:"id"`
	CreatedAt  time.Time  `json:"createdAt,omitempty"`
	UpdatedAt  time.Time  `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty" gorm:"index"`
	Name       string     `json:"name"`
	Account    string     `json:"account,omitempty"`
	Password   string
	InviteCode int    `json:"inviteCode,omitempty"`
	Desc       string `json:"desc"`
	Avatar     string `json:"avatar"`
}

type UserSimple struct {
	UId     int    `json:"id" gorm:"column:id"`
	UName   string `json:"name" gorm:"column:name"`
	UDesc   string `json:"desc" gorm:"column:desc"`
	UAvatar string `json:"avatar" gorm:"column:avatar"`
}

func User() *gorm.DB {
	return mysql.GetInstance().Model(&Users{})
}
