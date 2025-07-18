package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/pkg/mysql"
)

// CommentReaction 评论表情回复模型
type CommentReaction struct {
	ID           int            `gorm:"primarykey" json:"id"`
	CommentId    int            `json:"commentId" gorm:"column:comment_id"`
	UserId       int            `json:"userId" gorm:"column:user_id"`
	ReactionType string         `json:"reactionType" gorm:"column:reaction_type"`
	CreatedAt    time.LocalTime `json:"createdAt"`
	UpdatedAt    time.LocalTime `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	
	// 关联字段
	UserName     string `json:"userName" gorm:"-"`
	UserAvatar   string `json:"userAvatar" gorm:"-"`
}

// ExpressionType 表情类型配置模型
type ExpressionType struct {
	ID        int            `gorm:"primarykey" json:"id"`
	Code      string         `json:"code" gorm:"column:code"`
	Name      string         `json:"name" gorm:"column:name"`
	ImagePath string         `json:"imagePath" gorm:"column:image_path"`
	SortOrder int            `json:"sortOrder" gorm:"column:sort_order"`
	IsActive  bool           `json:"isActive" gorm:"column:is_active"`
	CreatedAt time.LocalTime `json:"createdAt"`
	UpdatedAt time.LocalTime `json:"updatedAt"`
}

// CommentReactionSummary 评论表情统计汇总
type CommentReactionSummary struct {
	CommentId    int    `json:"commentId"`
	ReactionType string `json:"reactionType"`
	Count        int    `json:"count"`
	UserReacted  bool   `json:"userReacted"` // 当前用户是否已回复此表情
	Users        []ReactionUser `json:"users" gorm:"-"` // 回复此表情的用户列表
}


// CommentReactionDetail 评论表情详情（包含用户信息）
type CommentReactionDetail struct {
	CommentId    int    `json:"commentId"`
	ReactionType string `json:"reactionType"`
	Users        []struct {
		UserId   int    `json:"userId"`
		UserName string `json:"userName"`
		UserAvatar string `json:"userAvatar"`
	} `json:"users"`
}

// CommentReactionDB 获取评论表情回复数据库实例
func CommentReactionDB() *gorm.DB {
	return mysql.GetInstance().Model(&CommentReaction{})
}

// ExpressionTypeDB 获取表情类型配置数据库实例
func ExpressionTypeDB() *gorm.DB {
	return mysql.GetInstance().Model(&ExpressionType{})
}

// TableName 指定表名
func (CommentReaction) TableName() string {
	return "comment_reactions"
}

// TableName 指定表名
func (ExpressionType) TableName() string {
	return "expression_types"
}