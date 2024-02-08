package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MemberInfos struct {
	ID        int       `gorm:"primarykey" json:"id"`
	Name      string    `json:"name"`
	Desc      string    `json:"desc"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func MemberInfo() *gorm.DB {
	return mysql.GetInstance().Model(&MemberInfos{})
}
