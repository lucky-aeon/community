package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Users struct {
	ID         int             `json:"id"`
	CreatedAt  time.LocalTime  `json:"createdAt,omitempty"`
	UpdatedAt  time.LocalTime  `json:"updatedAt,omitempty"`
	DeletedAt  *time.LocalTime `json:"deletedAt,omitempty" gorm:"index"`
	Name       string          `json:"name"`
	Account    string          `json:"account,omitempty"`
	Password   string
	InviteCode int    `json:"inviteCode,omitempty"`
	Desc       string `json:"desc"`
	Avatar     string `json:"avatar"`
	Subscribe  int    `json:"subscribe"` // 1: 未订阅站内消息 2:订阅站内消息 (发送邮箱)
}

type UserSimple struct {
	UId       int            `json:"id" gorm:"column:id"`
	UName     string         `json:"name" gorm:"column:name"`
	UDesc     string         `json:"desc" gorm:"column:desc"`
	UAvatar   string         `json:"avatar" gorm:"column:avatar"`
	Role      string         `json:"role" gorm:"column:u_role"`
	Account   string         `json:"account" gorm:"account"`
	CreatedAt time.LocalTime `json:"createdAt"`
	Subscribe int            `json:"subscribe"` // 1: 未订阅站内消息 2:订阅站内消息 (发送邮箱)
}

type LoginForm struct {
	Account  string `binding:"email" json:"account" msg:"邮箱格式错误"`
	Password string `binding:"required" json:"password" msg:"密码不能为空"`
}

func User() *gorm.DB {
	return mysql.GetInstance().Model(&Users{})
}
