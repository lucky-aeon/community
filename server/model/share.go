package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	localTime "xhyovo.cn/community/pkg/time"
)

// Share 分享记录表
type Share struct {
	ID           int                  `json:"id" gorm:"primaryKey;autoIncrement"`
	BusinessType string               `json:"business_type" gorm:"size:50;not null;comment:业务类型：ai_news, article, post等"`
	BusinessID   int                  `json:"business_id" gorm:"not null;comment:业务ID"`
	ShareToken   string               `json:"share_token" gorm:"size:32;uniqueIndex;not null;comment:分享令牌"`
	CreatorID    int                  `json:"creator_id" gorm:"comment:创建分享的用户ID"`
	TotalViews   int                  `json:"total_views" gorm:"default:0;comment:该分享链接的浏览量"`
	ExpireAt     *localTime.LocalTime `json:"expire_at" gorm:"comment:过期时间，NULL表示永久有效"`
	IsActive     int                  `json:"is_active" gorm:"default:1;comment:是否有效 1:有效 0:无效"`
	CreatedAt    localTime.LocalTime  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    localTime.LocalTime  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt       `json:"deleted_at,omitempty" gorm:"index"`
}

// ShareView 分享浏览记录表
type ShareView struct {
	ID        int                 `json:"id" gorm:"primaryKey;autoIncrement"`
	ShareID   int                 `json:"share_id" gorm:"not null;index;comment:分享记录ID"`
	VisitorIP string              `json:"visitor_ip" gorm:"size:45;comment:访问者IP"`
	VisitorID int                 `json:"visitor_id" gorm:"comment:访问者用户ID（如果已登录）"`
	UserAgent string              `json:"user_agent" gorm:"type:text;comment:用户代理"`
	Referer   string              `json:"referer" gorm:"size:500;comment:来源页面"`
	VisitedAt localTime.LocalTime `json:"visited_at" gorm:"autoCreateTime;comment:访问时间"`
}

// 表名设置
func (Share) TableName() string {
	return "shares"
}

func (ShareView) TableName() string {
	return "share_views"
}

// 数据库连接方法
func ShareModel() *gorm.DB {
	return mysql.GetInstance().Model(&Share{})
}

func ShareViewModel() *gorm.DB {
	return mysql.GetInstance().Model(&ShareView{})
}

// 便捷方法
func ShareDB() *gorm.DB {
	return ShareModel()
}

func ShareViewDB() *gorm.DB {
	return ShareViewModel()
}
