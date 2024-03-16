package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type OperLogs struct {
	Id            int            `gorm:"primarykey" json:"id"`
	RequestMethod string         `json:"requestMethod"`
	RequestInfo   string         `json:"requestInfo"`
	RequestBody   string         `json:"requestBody"`
	ResponseData  string         `json:"responseData"`
	UserId        int            `json:"userId"`
	Ip            string         `json:"ip"`
	ExecAt        string         `json:"execAt"`
	CreatedAt     time.LocalTime `json:"createdAt"`
}

type LogSearch struct {
	RequestMethod string `form:"requestMethod"`
	RequestInfo   string `form:"requestInfo"`
	UserName      string `form:"userName"`
	Ip            string `form:"ip"`
	StartTime     string `form:"startTime"`
	EndTime       string `form:"endTime"`
}

func OperLog() *gorm.DB {
	return mysql.GetInstance().Model(&OperLogs{})
}
