package model

import (
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/mysql"
	localTime "xhyovo.cn/community/pkg/time"
)

// AiNews AI新闻表
type AiNews struct {
	ID          int                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string              `json:"title" gorm:"size:500;not null;comment:标题"`
	Content     string              `json:"content" gorm:"type:longtext;comment:内容"`
	Summary     string              `json:"summary" gorm:"type:text;comment:摘要"`
	SourceURL   string              `json:"source_url" gorm:"size:500;comment:来源URL"`
	SourceName  string              `json:"source_name" gorm:"size:100;comment:来源网站名称"`
	PublishDate localTime.LocalTime `json:"publish_date" gorm:"comment:发布日期"`

	// 大模型自动标注
	Category string `json:"category" gorm:"size:50;comment:AI分类"`
	Tags     string `json:"tags" gorm:"size:200;comment:AI标签(逗号分隔)"`

	// 简单状态控制
	Status int `json:"status" gorm:"default:0;comment:状态 0:隐藏 1:显示"`

	// 存储原始网站文章ID
	Hash string `json:"hash" gorm:"size:64;comment:原始网站文章ID"`

	CreatedAt localTime.LocalTime `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt localTime.LocalTime `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt      `json:"deleted_at,omitempty" gorm:"index"`
}

// CrawlerConfig 爬虫配置表
type CrawlerConfig struct {
	ID          int                 `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string              `json:"name" gorm:"size:100;not null;comment:爬虫名称"`
	URL         string              `json:"url" gorm:"size:500;not null;comment:目标网站URL"`
	Status      int                 `json:"status" gorm:"default:1;comment:状态 1:启用 0:禁用"`
	CronExpr    string              `json:"cron_expr" gorm:"size:50;comment:定时表达式"`
	Description string              `json:"description" gorm:"size:500;comment:描述"`
	CreatedAt   localTime.LocalTime `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   localTime.LocalTime `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt      `json:"deleted_at,omitempty" gorm:"index"`
}

// 表名设置
func (AiNews) TableName() string {
	return "ai_news"
}

func (CrawlerConfig) TableName() string {
	return "crawler_configs"
}

// 数据库连接方法
func AiNewsModel() *gorm.DB {
	return mysql.GetInstance().Model(&AiNews{})
}

func CrawlerConfigModel() *gorm.DB {
	return mysql.GetInstance().Model(&CrawlerConfig{})
}
