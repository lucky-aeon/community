package services

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/crawler"
	"xhyovo.cn/community/server/model"
)

// CrawlerService 爬虫服务
type CrawlerService struct {
	aibaseCrawler *crawler.AIBaseCrawler
}

// NewCrawlerService 创建爬虫服务
func NewCrawlerService() *CrawlerService {
	return &CrawlerService{
		aibaseCrawler: crawler.NewAIBaseCrawler(0), // 默认ID会被覆盖
	}
}

// InitializeData 初始化数据 - 从指定ID向前爬取历史数据直到404（一次性调用）- 使用并发爬取
func (s *CrawlerService) InitializeData(startID int) error {
	log.Printf("开始初始化数据（并发模式）：从ID %d 向前爬取直到404", startID)

	// 使用并发爬取，默认使用10个协程
	workers := 10

	log.Printf("使用 %d 个协程进行并发爬取", workers)

	// 调用新的并发爬取方法（爬取直到404）
	return s.ConcurrentCrawlUntil404(startID, workers)
}

// ScheduledCrawl 定时爬取 - 从数据库最新AIBase文章ID向后递增爬取新文章
func (s *CrawlerService) ScheduledCrawl() error {
	log.Println("开始定时爬取新文章")

	db := mysql.GetInstance()

	// 获取AIBase来源的最新文章ID（从source_url中解析）
	var latestArticle model.AiNews
	result := db.Where("source_name = 'AIBase' AND source_url LIKE '%/daily/%'").
		Order("CAST(SUBSTRING_INDEX(source_url, '/', -1) AS UNSIGNED) DESC").
		First(&latestArticle)

	var startID int
	if result.Error != nil {
		// 数据库为空，从默认ID开始
		startID = 18330 // 可以根据实际情况调整
		log.Printf("数据库中没有AIBase文章，从默认ID %d 开始爬取", startID)
	} else {
		// 从URL中解析出文章ID，例如：https://www.aibase.com/zh/daily/18391 -> 18391
		urlParts := strings.Split(latestArticle.SourceURL, "/")
		if len(urlParts) > 0 {
			latestID, err := strconv.Atoi(urlParts[len(urlParts)-1])
			if err != nil {
				log.Printf("解析最新文章ID失败: %v，使用默认起始ID", err)
				startID = 18330
			} else {
				startID = latestID + 1 // 从最新ID的下一个开始
				log.Printf("从最新AIBase文章ID %d 的下一个ID %d 开始爬取", latestID, startID)
			}
		} else {
			startID = 18330
		}
	}

	currentID := startID
	successCount := 0
	consecutive404Count := 0
	max404Count := 3 // 连续3个404就停止

	for {
		if consecutive404Count >= max404Count {
			log.Printf("连续遇到 %d 个404，停止爬取", max404Count)
			break
		}

		url := fmt.Sprintf("https://www.aibase.com/zh/daily/%d", currentID)
		log.Printf("正在爬取: %s", url)

		article, err := s.aibaseCrawler.FetchSingleArticle(url, currentID)
		if err != nil {
			if err.Error() == "文章不存在 (404)" {
				consecutive404Count++
				log.Printf("遇到404 (ID: %d)，连续404次数: %d", currentID, consecutive404Count)
			} else {
				log.Printf("爬取失败 (ID: %d): %v", currentID, err)
			}
			currentID++ // 向后递增
			continue
		}

		// 重置404计数
		consecutive404Count = 0

		// 检查文章是否已存在（通过source_url检查）
		var existingArticle model.AiNews
		result := db.Where("source_url = ?", url).First(&existingArticle)
		if result.Error == nil {
			log.Printf("文章已存在，跳过 (URL: %s)", url)
			currentID++
			continue
		}

		// 设置hash为原始文章ID，不设置主键ID（让数据库自动分配）
		article.ID = 0                              // 清空ID，让数据库自动分配
		article.Hash = fmt.Sprintf("%d", currentID) // hash存储原始AIBase文章ID

		// 保存文章
		if err := db.Create(article).Error; err != nil {
			log.Printf("保存文章失败 (原始ID: %d): %v", currentID, err)
		} else {
			successCount++
			log.Printf("成功爬取并保存新文章: %s (原始ID: %d, 数据库ID: %d)", article.Title, currentID, article.ID)
		}

		currentID++ // 向后递增

		// 添加延迟避免请求过快
		time.Sleep(1 * time.Second)
	}

	log.Printf("定时爬取完成，成功爬取 %d 篇新文章", successCount)
	return nil
}

// 保留原有的并发爬取方法（用于测试或特殊场景）
func (s *CrawlerService) ConcurrentCrawlFromID(startID, count, workers int) ([]*model.AiNews, error) {
	log.Printf("开始并发爬取：从ID %d 向前爬取 %d 篇文章，使用 %d 个协程", startID, count, workers)

	jobs := make(chan int, count)
	results := make(chan *model.AiNews, count)
	errors := make(chan error, count)

	var wg sync.WaitGroup

	// 启动工作协程
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go s.crawlWorker(jobs, results, errors, &wg)
	}

	// 发送任务
	go func() {
		defer close(jobs)
		for i := 0; i < count; i++ {
			jobs <- startID - i
		}
	}()

	// 等待所有协程完成
	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	// 收集结果
	var articles []*model.AiNews
	var crawlErrors []error

	for results != nil || errors != nil {
		select {
		case article, ok := <-results:
			if !ok {
				results = nil
			} else if article != nil {
				articles = append(articles, article)
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
			} else if err != nil {
				crawlErrors = append(crawlErrors, err)
			}
		}
	}

	log.Printf("并发爬取完成，成功: %d 篇，失败: %d 篇", len(articles), len(crawlErrors))
	return articles, nil
}

func (s *CrawlerService) crawlWorker(jobs <-chan int, results chan<- *model.AiNews, errors chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range jobs {
		url := fmt.Sprintf("https://www.aibase.com/zh/daily/%d", id)
		article, err := s.aibaseCrawler.FetchSingleArticle(url, id)

		if err != nil {
			errors <- fmt.Errorf("ID %d: %v", id, err)
			continue
		}

		results <- article
	}
}

func (s *CrawlerService) BatchInsertArticles(articles []*model.AiNews) error {
	if len(articles) == 0 {
		return nil
	}

	db := mysql.GetInstance()

	log.Printf("开始批量插入 %d 篇文章", len(articles))

	// 使用事务批量插入
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	for _, article := range articles {
		// 检查文章是否已存在
		var existingArticle model.AiNews
		result := tx.Where("id = ?", article.ID).First(&existingArticle)
		if result.Error == nil {
			log.Printf("文章已存在，跳过: ID %d", article.ID)
			continue
		}

		if err := tx.Create(article).Error; err != nil {
			log.Printf("插入文章失败: ID %d, 错误: %v", article.ID, err)
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Printf("批量插入完成")
	return nil
}

func (s *CrawlerService) ConcurrentCrawlAndSave(startID, count, workers int) error {
	articles, err := s.ConcurrentCrawlFromID(startID, count, workers)
	if err != nil {
		return err
	}

	return s.BatchInsertArticles(articles)
}

// ConcurrentCrawlUntil404 并发爬取直到遇到404
func (s *CrawlerService) ConcurrentCrawlUntil404(startID, workers int) error {
	log.Printf("开始并发爬取：从ID %d 向前爬取直到连续404，使用 %d 个协程", startID, workers)

	db := mysql.GetInstance()

	// 用于跟踪连续404的计数
	consecutive404Count := 0
	max404Count := 10 // 连续10个404就停止

	currentID := startID
	successCount := 0
	totalProcessed := 0

	// 创建工作队列和结果队列
	jobs := make(chan int, workers*2) // 缓冲区大小为协程数的2倍
	results := make(chan crawlResult, workers*2)

	var wg sync.WaitGroup

	// 启动工作协程
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go s.crawlWorkerWithResult(jobs, results, &wg)
	}

	// 结果处理协程
	done := make(chan bool)
	go func() {
		defer close(done)

		for result := range results {
			totalProcessed++

			if result.error != nil {
				if result.error.Error() == "文章不存在 (404)" {
					consecutive404Count++
					log.Printf("遇到404 (ID: %d)，连续404次数: %d", result.id, consecutive404Count)

					if consecutive404Count >= max404Count {
						log.Printf("连续遇到 %d 个404，停止爬取", max404Count)
						close(jobs) // 停止发送新任务
						return
					}
				} else {
					log.Printf("爬取失败 (ID: %d): %v", result.id, result.error)
				}
				continue
			}

			// 重置404计数
			consecutive404Count = 0

			// 通过source_url检查文章是否已存在
			url := fmt.Sprintf("https://www.aibase.com/zh/daily/%d", result.id)
			var existingArticle model.AiNews
			dbResult := db.Where("source_url = ?", url).First(&existingArticle)
			if dbResult.Error == nil {
				log.Printf("文章已存在，跳过: URL %s", url)
				continue
			}

			// 设置hash为原始文章ID，清除主键ID让数据库自动分配
			result.article.ID = 0                              // 清空ID，让数据库自动分配
			result.article.Hash = fmt.Sprintf("%d", result.id) // hash存储原始AIBase文章ID

			// 保存文章
			if err := db.Create(result.article).Error; err != nil {
				log.Printf("保存文章失败 (原始ID: %d): %v", result.id, err)
			} else {
				successCount++
				log.Printf("成功爬取并保存文章: %s (原始ID: %d, 数据库ID: %d)", result.article.Title, result.id, result.article.ID)
			}
		}
	}()

	// 发送任务
	go func() {
		defer close(jobs)

		for {
			select {
			case jobs <- currentID:
				currentID-- // 向前递减
			case <-done:
				return
			}
		}
	}()

	// 等待所有工作协程完成
	wg.Wait()
	close(results)

	// 等待结果处理完成
	<-done

	log.Printf("并发爬取完成，总处理: %d 篇，成功: %d 篇", totalProcessed, successCount)
	return nil
}

// crawlResult 爬取结果结构
type crawlResult struct {
	id      int
	article *model.AiNews
	error   error
}

// crawlWorkerWithResult 工作协程（返回详细结果）
func (s *CrawlerService) crawlWorkerWithResult(jobs <-chan int, results chan<- crawlResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for id := range jobs {
		url := fmt.Sprintf("https://www.aibase.com/zh/daily/%d", id)
		article, err := s.aibaseCrawler.FetchSingleArticle(url, id)

		results <- crawlResult{
			id:      id,
			article: article,
			error:   err,
		}

		// 添加小延迟避免请求过快
		time.Sleep(100 * time.Millisecond)
	}
}
