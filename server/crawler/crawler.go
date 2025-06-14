package crawler

import (
	"xhyovo.cn/community/server/model"
)

// CrawlerInterface 爬虫接口，不同网站实现不同的爬虫策略
type CrawlerInterface interface {
	// GetName 获取爬虫名称
	GetName() string
	// GetURL 获取目标URL
	GetURL() string
	// Crawl 执行爬取，返回新闻列表
	Crawl() ([]*model.AiNews, error)
	// IsEnabled 检查是否启用
	IsEnabled() bool
}

// CrawlerManager 爬虫管理器
type CrawlerManager struct {
	crawlers []CrawlerInterface
}

// NewCrawlerManager 创建爬虫管理器
func NewCrawlerManager() *CrawlerManager {
	return &CrawlerManager{
		crawlers: make([]CrawlerInterface, 0),
	}
}

// RegisterCrawler 注册爬虫
func (cm *CrawlerManager) RegisterCrawler(crawler CrawlerInterface) {
	cm.crawlers = append(cm.crawlers, crawler)
}

// RunAll 运行所有启用的爬虫
func (cm *CrawlerManager) RunAll() ([]*model.AiNews, error) {
	var allNews []*model.AiNews

	for _, crawler := range cm.crawlers {
		if !crawler.IsEnabled() {
			continue
		}

		news, err := crawler.Crawl()
		if err != nil {
			// 记录错误但继续执行其他爬虫
			continue
		}

		allNews = append(allNews, news...)
	}

	return allNews, nil
}

// GetEnabledCrawlers 获取所有启用的爬虫
func (cm *CrawlerManager) GetEnabledCrawlers() []CrawlerInterface {
	var enabled []CrawlerInterface
	for _, crawler := range cm.crawlers {
		if crawler.IsEnabled() {
			enabled = append(enabled, crawler)
		}
	}
	return enabled
}
