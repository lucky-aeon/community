package model

import (
	"fmt"
	"gorm.io/gorm"
	sysTime "time"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type Meetings struct {
	Id               int                `gorm:"primarykey" json:"id"`
	InitiatorId      int                `gorm:"not null" json:"initiatorId"` // 申请人
	Title            string             `gorm:"not null" json:"title" binding:"required" msg:"标题不能未空"`
	Description      string             `gorm:"not null" json:"description"`                                     // 描述
	Record           string             `gorm:"not null" json:"record"`                                          // 内容，用户不可更改，由管理员进行管理
	InitiatorTime    time.LocalTime     `gorm:"not null" json:"initiatorTime" binding:"required" msg:"申请时间不可为空"` // 申请时间
	MeetingStartTime *time.LocalTime    `gorm:"not null" json:"meetingStartTime"`                                // 会议开始时间
	MeetingEndTime   *time.LocalTime    `gorm:"not null" json:"meetingEndTime"`                                  // 会议结束时间
	SignupEndTime    *time.LocalTime    `gorm:"not null" json:"signupEndTime"`                                   // 报名截止时间
	State            string             `gorm:"not null" json:"state"`                                           // 状态
	StateMessage     string             `gorm:"not null" json:"stateMessage"`                                    // 状态消息
	MeetingLink      string             `gorm:"not null" json:"meetingLink"`
	UpdatedAt        time.LocalTime     `gorm:"not null" json:"updatedAt"`
	CreatedAt        time.LocalTime     `json:"createdAt"`
	DeletedAt        gorm.DeletedAt     `gorm:"index" json:"deletedAt"`
	InitiatorName    string             `json:"initiatorName" gorm:"-"` // 发起者昵称
	JoinUsers        []MeetingJoinUsers `json:"joinUsers" gorm:"-"`     // 参与人
	JoinUserCount    int                `json:"joinUserCount" gorm:"-"`
	InitiatorAvatar  string             `json:"initiatorAvatar" gorm:"-"`
}

func Meeting() *gorm.DB {
	return mysql.GetInstance().Model(&Meetings{})
}

func (m *Meetings) PrintLog() string {

	return fmt.Sprintf("会议id:%d,会议标题:%s,会议状态:%s,会议报名截止时间:%v,会议开始时间:%v,会议结束时间:%v",
		m.Id, m.Title, m.State, sysTime.Time(*m.SignupEndTime).Format("2006-01-02 15:04:05"),
		sysTime.Time(*m.MeetingStartTime).Format("2006-01-02 15:04:05"),
		sysTime.Time(*m.MeetingEndTime).Format("2006-01-02 15:04:05"))
}
