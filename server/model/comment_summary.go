package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

type CommentSummary struct {
	ID            int            `gorm:"primarykey" json:"id"`
	BusinessId    int            `json:"businessId" gorm:"not null;comment:业务对象ID"`
	TenantId      int            `json:"tenantId" gorm:"not null;comment:租户ID"`
	Summary       string         `json:"summary" gorm:"type:text;not null;comment:AI总结内容"`
	CommentCount  int            `json:"commentCount" gorm:"default:0;comment:参与总结的评论数量"`
	LastCommentId int            `json:"lastCommentId" gorm:"default:0;comment:最后处理的评论ID"`
	CreatedAt     time.LocalTime `json:"createdAt"`
	UpdatedAt     time.LocalTime `json:"updatedAt"`
}

func (CommentSummary) TableName() string {
	return "comment_summaries"
}

func CommentSummaryModel() *gorm.DB {
	return mysql.GetInstance().Model(&CommentSummary{})
}
