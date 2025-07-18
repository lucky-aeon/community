package model

import (
	"gorm.io/gorm"
	"time"
	"xhyovo.cn/community/pkg/mysql"
)

// Reaction 通用表情回复模型
type Reaction struct {
	ID           int            `json:"id" gorm:"primaryKey;autoIncrement"`
	BusinessType int            `json:"businessType" gorm:"not null;comment:业务类型: 0=文章, 1=评论, 2=课程, 3=分享会, 4=AI日报"`
	BusinessId   int            `json:"businessId" gorm:"not null;comment:业务ID"`
	UserId       int            `json:"userId" gorm:"not null;comment:用户ID"`
	ReactionType string         `json:"reactionType" gorm:"type:varchar(50);not null;comment:表情类型"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// ReactionSummary 表情统计信息
type ReactionSummary struct {
	BusinessType int            `json:"businessType"`
	BusinessId   int            `json:"businessId"`
	ReactionType string         `json:"reactionType"`
	Count        int            `json:"count"`
	UserReacted  bool           `json:"userReacted"`
	Users        []ReactionUser `json:"users"`
}

// ReactionUser 表情回复用户信息
type ReactionUser struct {
	UserId     int    `json:"userId"`
	UserName   string `json:"userName"`
	UserAvatar string `json:"userAvatar"`
}

// 业务类型常量
const (
	BusinessTypeArticle = 0 // 文章
	BusinessTypeComment = 1 // 评论
	BusinessTypeCourse  = 2 // 课程
	BusinessTypeMeeting = 3 // 分享会
	BusinessTypeAINews  = 4 // AI日报
)

// GetBusinessTypeName 获取业务类型名称
func GetBusinessTypeName(businessType int) string {
	switch businessType {
	case BusinessTypeArticle:
		return "文章"
	case BusinessTypeComment:
		return "评论"
	case BusinessTypeCourse:
		return "课程"
	case BusinessTypeMeeting:
		return "分享会"
	case BusinessTypeAINews:
		return "AI日报"
	default:
		return "未知"
	}
}

// ReactionDB 返回反应数据库连接
func ReactionDB() *gorm.DB {
	return mysql.GetInstance()
}

// TableName 指定表名
func (Reaction) TableName() string {
	return "reactions"
}