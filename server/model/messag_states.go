package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"

	"xhyovo.cn/community/pkg/mysql"
)

type MessageStates struct {
	ID        int            `gorm:"primarykey" json:"id"`
	Content   string         `json:"content"`   // 内容
	From      int            `json:"from"`      // 发送人
	To        int            `json:"to"`        // 接收人
	State     int            `json:"state"`     // 消息状态:未读:1 已读:0
	Type      int            `json:"type"`      // 消息类型:通知消息:1 @:2
	ArticleId int            `json:"articleId"` // 文章id
	EventId   int            `json:"eventId"`   // 事件id
	CreatedAt time.LocalTime `json:"createdAt"` // 发送时间
}

func MessageState() *gorm.DB {
	return mysql.GetInstance().Model(&MessageStates{})
}
