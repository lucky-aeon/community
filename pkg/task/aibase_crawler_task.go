package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/crawler"
	"xhyovo.cn/community/server/model"
)

// AIBaseCrawlerTask AIBase爬虫定时任务
type AIBaseCrawlerTask struct {
	name     string
	cronExpr string
	enabled  bool
	config   AIBaseCrawlerConfig
}

// AIBaseCrawlerConfig 爬虫配置
type AIBaseCrawlerConfig struct {
	MaxConsecutive404 int           `json:"max_consecutive_404"` // 最大连续404次数
	RequestInterval   time.Duration `json:"request_interval"`    // 请求间隔
	BatchSize         int           `json:"batch_size"`          // 批处理大小
	MaxRetries        int           `json:"max_retries"`         // 最大重试次数
}

// NewAIBaseCrawlerTask 创建AIBase爬虫任务
func NewAIBaseCrawlerTask(cronExpr string) *AIBaseCrawlerTask {
	return &AIBaseCrawlerTask{
		name:     "AIBase增量爬取",
		cronExpr: cronExpr,
		enabled:  true,
		config: AIBaseCrawlerConfig{
			MaxConsecutive404: 50,
			RequestInterval:   1 * time.Second,
			BatchSize:         10,
			MaxRetries:        3,
		},
	}
}

// GetName 获取任务名称
func (t *AIBaseCrawlerTask) GetName() string {
	return t.name
}

// GetCronExpr 获取cron表达式
func (t *AIBaseCrawlerTask) GetCronExpr() string {
	return t.cronExpr
}

// IsEnabled 是否启用
func (t *AIBaseCrawlerTask) IsEnabled() bool {
	return t.enabled
}

// Execute 执行任务
func (t *AIBaseCrawlerTask) Execute(ctx context.Context) error {
	log.Infof("开始执行AIBase增量爬取任务")

	// 获取最新文章ID
	latestID, err := t.getLatestArticleID()
	if err != nil {
		return fmt.Errorf("获取最新文章ID失败: %v", err)
	}

	// 从下一个ID开始爬取
	startID := latestID + 1
	log.Infof("从ID %d 开始增量爬取", startID)

	// 执行增量爬取
	result, err := t.performIncrementalCrawl(ctx, startID)
	if err != nil {
		return fmt.Errorf("增量爬取失败: %v", err)
	}

	log.Infof("AIBase增量爬取完成，成功: %d 篇，失败: %d 篇，最后ID: %d",
		result.SuccessCount, result.FailedCount, result.LastProcessedID)

	return nil
}

// getLatestArticleID 获取数据库中AIBase的最新文章ID
func (t *AIBaseCrawlerTask) getLatestArticleID() (int, error) {
	db := mysql.GetInstance()

	var result struct {
		ArticleID int `json:"article_id"`
	}

	// 查询AIBase最新文章ID
	err := db.Raw(`
		SELECT 
			CAST(SUBSTRING_INDEX(source_url, '/', -1) AS UNSIGNED) as article_id
		FROM ai_news 
		WHERE source_name = 'AIBase' 
			AND source_url LIKE '%/daily/%'
			AND source_url REGEXP '/daily/[0-9]+$'
		ORDER BY article_id DESC 
		LIMIT 1
	`).Scan(&result).Error

	if err != nil {
		// 如果没有找到记录，从默认ID开始
		log.Warnf("未找到AIBase文章记录，使用默认起始ID: %v", err)
		return 18330, nil // 默认起始ID
	}

	log.Infof("找到最新AIBase文章ID: %d", result.ArticleID)
	return result.ArticleID, nil
}

// CrawlResult 爬取结果
type CrawlResult struct {
	SuccessCount    int   `json:"success_count"`
	FailedCount     int   `json:"failed_count"`
	LastProcessedID int   `json:"last_processed_id"`
	ProcessedIDs    []int `json:"processed_ids"`
	Consecutive404  int   `json:"consecutive_404"`
}

// performIncrementalCrawl 执行增量爬取
func (t *AIBaseCrawlerTask) performIncrementalCrawl(ctx context.Context, startID int) (*CrawlResult, error) {
	db := mysql.GetInstance()
	aibaseCrawler := crawler.NewAIBaseCrawler(startID)

	result := &CrawlResult{
		ProcessedIDs: make([]int, 0),
	}

	currentID := startID
	consecutive404Count := 0

	for {
		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		// 检查是否达到最大连续404次数
		if consecutive404Count >= t.config.MaxConsecutive404 {
			log.Infof("连续遇到 %d 个404，停止爬取", t.config.MaxConsecutive404)
			break
		}

		url := fmt.Sprintf("https://www.aibase.com/zh/daily/%d", currentID)
		log.Infof("正在爬取: %s", url)

		// 添加到处理列表
		result.ProcessedIDs = append(result.ProcessedIDs, currentID)
		result.LastProcessedID = currentID

		// 执行爬取
		article, err := t.crawlSingleArticle(aibaseCrawler, url, currentID)
		if err != nil {
			if strings.Contains(err.Error(), "404") {
				consecutive404Count++
				log.Infof("遇到404 (ID: %d)，连续404次数: %d", currentID, consecutive404Count)
			} else {
				log.Errorf("爬取失败 (ID: %d): %v", currentID, err)
			}
			result.FailedCount++
		} else {
			// 重置404计数
			consecutive404Count = 0

			// 保存文章（FirstOrCreate会自动处理重复）
			if err := t.saveArticle(db, article, currentID); err != nil {
				log.Errorf("保存文章失败 (ID: %d): %v", currentID, err)
				result.FailedCount++
			} else {
				result.SuccessCount++
				log.Infof("处理文章成功: %s (ID: %d)", article.Title, currentID)
			}
		}

		currentID++

		// 请求间隔控制
		time.Sleep(t.config.RequestInterval)
	}

	result.Consecutive404 = consecutive404Count
	return result, nil
}

// crawlSingleArticle 爬取单篇文章（带重试）
func (t *AIBaseCrawlerTask) crawlSingleArticle(crawler *crawler.AIBaseCrawler, url string, id int) (*model.AiNews, error) {
	var lastErr error

	for retry := 0; retry < t.config.MaxRetries; retry++ {
		article, err := crawler.FetchSingleArticle(url, id)
		if err == nil {
			return article, nil
		}

		lastErr = err
		if strings.Contains(err.Error(), "404") {
			// 404错误不重试
			break
		}

		if retry < t.config.MaxRetries-1 {
			log.Warnf("爬取失败，第 %d 次重试 (ID: %d): %v", retry+1, id, err)
			time.Sleep(time.Duration(retry+1) * time.Second)
		}
	}

	return nil, lastErr
}

// saveArticle 保存文章到数据库
func (t *AIBaseCrawlerTask) saveArticle(db *gorm.DB, article *model.AiNews, originalID int) error {
	// 清空主键ID，让数据库自动分配
	article.ID = 0
	// 将原始ID存储在hash字段中
	article.Hash = fmt.Sprintf("%d", originalID)

	// 使用FirstOrCreate避免重复插入
	var existingArticle model.AiNews
	result := db.Where("source_url = ?", article.SourceURL).FirstOrCreate(&existingArticle, article)

	if result.Error != nil {
		return result.Error
	}

	// 如果是新创建的记录，记录日志
	if result.RowsAffected > 0 {
		log.Infof("新文章已保存: %s", article.Title)
	} else {
		log.Infof("文章已存在，跳过: %s", article.Title)
	}

	return nil
}

// SetEnabled 设置启用状态
func (t *AIBaseCrawlerTask) SetEnabled(enabled bool) {
	t.enabled = enabled
}

// UpdateConfig 更新配置
func (t *AIBaseCrawlerTask) UpdateConfig(config AIBaseCrawlerConfig) {
	t.config = config
}
