package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type MeetingJoinUsers struct {
	Id        int       `gorm:"primarykey" json:"id"`
	MeetingId int       `json:"meetingId"`
	UserId    int       `json:"userId"`
	UserName  string    `json:"userName"`
	CreatedAt time.Time `json:"createdAt"`
}

func MeetingJoinUser() *gorm.DB {
	return mysql.GetInstance().Model(&MeetingJoinUsers{})
}
