package request

import (
	"xhyovo.cn/community/pkg/time"
)

type ReqMeeting struct {
	Id            int            `json:"id"`
	Title         string         `gorm:"not null" json:"title" binding:"required" msg:"标题不能为空"`
	Description   string         `json:"description"`
	InitiatorTime time.LocalTime `gorm:"not null" json:"initiatorTime" binding:"required" msg:"申请时间不可为空"` // 申请时间

}

// 会议通过 req
type ReqApproveMeeting struct {
	Id               int             `json:"id" binding:"required" msg:"id不能未空"`
	MeetingStartTime *time.LocalTime `gorm:"not null" json:"meetingStartTime" binding:"required" msg:"会议开始时间不能为空"` // 会议开始时间
	MeetingEndTime   *time.LocalTime `gorm:"not null" json:"meetingEndTime" binding:"required" msg:"会议结束时间不能为空"`   // 会议结束时间
	SignupEndTime    *time.LocalTime `gorm:"not null" json:"signupEndTime" binding:"required" msg:"报名截止时间不能为空"`    // 报名截止时间
	MeetingLink      string          `gorm:"not null" json:"meetingLink" binding:"required" msg:"会议链接不可为空"`
}

type ReqPassMeeting struct {
	Id          int    `json:"id" binding:"required" msg:"id不能为空"`
	PassMessage string `json:"passMessage" binding:"required" msg:"pass理由不能为空"`
}

type ReqRecordMeeting struct {
	Id     int    `json:"id" binding:"required" msg:"id不能为空"`
	Record string `json:"record"`
}
