package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Meetings struct {
	Id               int                `gorm:"primarykey" json:"id"`
	InitiatorId      int                `gorm:"not null" json:"initiator_id"`
	Title            string             `gorm:"not null" json:"title"`
	Description      string             `gorm:"not null" json:"description"`
	Context          string             `gorm:"not null" json:"context"`
	InitiatorTime    time.LocalTime     `gorm:"not null" json:"initiator_time"`
	MeetingStartTime time.LocalTime     `gorm:"not null" json:"metting_start_time"`
	SignupEndTime    time.LocalTime     `gorm:"not null" json:"signup_end_time"`
	State            string             `gorm:"not null" json:"state"`
	StateMessage     string             `gorm:"not null" json:"state_message"`
	UpdatedAt        time.LocalTime     `gorm:"not null" json:"updated_at"`
	CreatedAt        time.LocalTime     `json:"createdAt"`
	DeletedAt        gorm.DeletedAt     `gorm:"index" json:"deleted_at"`
	InitiatorName    string             `json:"initiatorName" gorm:"-"` // 发起者昵称
	JoinUsers        []MeetingJoinUsers `json:"joinUsers" gorm:"-"`     // 参与人
	JoinUserCount    int                `json:"joinUserCount" gorm:"-"`
}

func Meeting() *gorm.DB {
	return mysql.GetInstance().Model(&Meetings{})
}
