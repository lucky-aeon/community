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
	UserAgent     string         `json:"userAgent"`
	Platform      string         `json:"platform"`
	ExecAt        string         `json:"execAt"`
	CreatedAt     time.LocalTime `json:"createdAt"`
	UserName      string         `gorm:"-" json:"userName"`
}

type LogSearch struct {
	RequestMethod string `form:"requestMethod"`
	RequestInfo   string `form:"requestInfo"`
	UserName      string `form:"userName"`
	Ip            string `form:"ip"`
	StartTime     string `form:"startTime"`
	EndTime       string `form:"endTime"`
	Account       string `form:"account"`
}

type LoginLogs struct {
	Id        int            `gorm:"primarykey" json:"id"`
	Account   string         `json:"account"`
	State     string         `json:"state"`
	Browser   string         `json:"browser"`
	Equipment string         `json:"equipment"`
	Ip        string         `json:"ip"`
	CreatedAt time.LocalTime `json:"createdAt"`
}

func OperLog() *gorm.DB {
	return mysql.GetInstance().Model(&OperLogs{})
}

func LoginLog() *gorm.DB {
	return mysql.GetInstance().Model(&LoginLogs{})
}
