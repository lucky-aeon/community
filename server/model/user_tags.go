package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type UserTags struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	CreatedAt *time.LocalTime `json:"createdAt"`
}

type UserTagRelations struct {
	ID        int `json:"id"`
	UserId    int `json:"userId"`
	UserTagId int `json:"UserTagId"`
}

func UserTag() *gorm.DB {
	return mysql.GetInstance().Model(&UserTags{})
}

func UserTagRelation() *gorm.DB {
	return mysql.GetInstance().Model(&UserTagRelations{})
}
