package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/pkg/mysql"
)

type MemberInfos struct {
	ID        int            `gorm:"primarykey" json:"id"`
	Name      string         `json:"name"`
	Desc      string         `json:"desc"`
	Money     int            `json:"money"`
	CreatedAt time.LocalTime `json:"createdAt"`
	UpdatedAt time.LocalTime `json:"updatedAt"`
}

func MemberInfo() *gorm.DB {
	return mysql.GetInstance().Model(&MemberInfos{})
}
