package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

type OperLogs struct {
	Id           int       `gorm:"primarykey" json:"id"`
	RequestType  string    `json:"requestType"`
	RequestInfo  string    `json:"requestInfo"`
	RequestBody  string    `json:"requestBody"`
	ResponseData string    `json:"responseData"`
	UserId       int       `json:"userId"`
	Ip           string    `json:"ip"`
	ExecAt       string    `json:"execAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

func OperLog() *gorm.DB {
	return mysql.GetInstance().Model(&OperLogs{})
}
