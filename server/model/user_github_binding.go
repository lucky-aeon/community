package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/pkg/time"
)

// UserGitHubBinding 用户GitHub账号绑定模型
type UserGitHubBinding struct {
	ID              int            `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID          int            `json:"user_id" gorm:"column:user_id;not null;uniqueIndex;comment:用户ID"`
	GitHubID        int64          `json:"github_id" gorm:"column:github_id;not null;uniqueIndex;comment:GitHub用户ID"`
	GitHubUsername  string         `json:"github_username" gorm:"column:github_username;type:varchar(100);not null;comment:GitHub用户名"`
	GitHubEmail     string         `json:"github_email" gorm:"column:github_email;type:varchar(255);comment:GitHub邮箱"`
	GitHubAvatar    string         `json:"github_avatar" gorm:"column:github_avatar;type:varchar(500);comment:GitHub头像URL"`
	AccessToken     string         `json:"access_token" gorm:"column:access_token;type:varchar(255);comment:GitHub Access Token"`
	BoundAt         time.LocalTime `json:"bound_at" gorm:"column:bound_at;comment:绑定时间"`
	CreatedAt       time.LocalTime `json:"created_at" gorm:"column:created_at"`
	UpdatedAt       time.LocalTime `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (UserGitHubBinding) TableName() string {
	return "user_github_bindings"
}

// UserGitHubBindingModel 获取模型实例
func UserGitHubBindingModel() *gorm.DB {
	return mysql.GetInstance().Model(&UserGitHubBinding{})
}

// GitHubUser GitHub用户信息结构体（用于OAuth回调处理）
type GitHubUser struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
}

// GitHubOAuthResponse GitHub OAuth令牌响应结构体
type GitHubOAuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}