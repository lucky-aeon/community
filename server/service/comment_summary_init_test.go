package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

// ==================== 初始化配置区域 ====================
// 请根据需要调整以下参数：

const (
	// 批量处理配置
	BATCH_SIZE = 10 // 每批处理的数量，避免一次性处理太多数据

	// 过滤条件
	MIN_COMMENT_COUNT = 3  // 最少评论数量，少于此数量的不生成总结
	MAX_COMMENT_COUNT = 50 // 最多评论数量，超过此数量的可能需要更长时间

	// 并发控制
	MAX_CONCURRENT = 3 // 最大并发数，避免对LLM服务压力过大

	// 延迟控制
	DELAY_BETWEEN_REQUESTS = 1 * time.Second // 请求间延迟，避免频率过高
)

// 如果只想初始化特定范围的数据，可以设置以下过滤条件：
// - 设置为0表示不限制
var (
	FILTER_ARTICLE_IDS = []int{} // 指定文章ID列表，空表示处理所有
	FILTER_SECTION_IDS = []int{} // 指定章节ID列表，空表示处理所有
)

// ====================================================

// InitializationStats 初始化统计信息
type InitializationStats struct {
	TotalProcessed int
	SuccessCount   int
	SkippedCount   int
	ErrorCount     int
	StartTime      time.Time
	TotalDuration  time.Duration
	Errors         []string
}

// 初始化测试环境
func setupInitTestEnvironment(t *testing.T) {
	log.Init()
	
	chinaLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("加载时区失败: %v", err)
	}
	time.Local = chinaLoc
	
	config.Init()
	appConfig := config.GetInstance()
	
	if appConfig.LLMConfig.ApiKey == "" || appConfig.LLMConfig.Url == "" {
		t.Fatalf("LLM配置不完整，请检查环境变量: LLM_API_KEY, LLM_URL, LLM_MODEL")
	}
	
	db := appConfig.DbConfig
	mysql.Init(db.Username, db.Password, db.Address, db.Database)
	
	t.Logf("✅ 初始化环境准备完成")
	t.Logf("   - LLM服务: %s", appConfig.LLMConfig.Url)
	t.Logf("   - LLM模型: %s", appConfig.LLMConfig.Model)
}

// 获取需要初始化的文章列表
func getArticlesNeedInitialization(t *testing.T) []model.Articles {
	var articles []model.Articles
	
	// 构建查询条件
	query := model.Article().Select("id, title, user_id")
	
	// 如果指定了文章ID列表
	if len(FILTER_ARTICLE_IDS) > 0 {
		query = query.Where("id IN ?", FILTER_ARTICLE_IDS)
	}
	
	// 查询所有文章
	query.Find(&articles)
	
	// 过滤：只处理有足够评论数量且还没有总结的文章
	var needInitArticles []model.Articles
	for _, article := range articles {
		// 检查评论数量
		var commentCount int64
		model.Comment().Where("business_id = ? AND tenant_id = ?", article.ID, 0).Count(&commentCount)
		
		if int(commentCount) < MIN_COMMENT_COUNT {
			continue // 评论太少，跳过
		}
		
		if int(commentCount) > MAX_COMMENT_COUNT {
			t.Logf("⚠️  文章 '%s' (ID: %d) 有 %d 条评论，可能需要较长处理时间", article.Title, article.ID, commentCount)
		}
		
		// 检查是否已有总结
		var summaryCount int64
		model.CommentSummaryModel().Where("business_id = ? AND tenant_id = ?", article.ID, 0).Count(&summaryCount)
		
		if summaryCount > 0 {
			continue // 已有总结，跳过
		}
		
		needInitArticles = append(needInitArticles, article)
	}
	
	return needInitArticles
}

// 获取需要初始化的章节列表
func getSectionsNeedInitialization(t *testing.T) []model.CoursesSections {
	var sections []model.CoursesSections
	
	// 构建查询条件
	query := model.CoursesSection().Select("id, title, user_id")
	
	// 如果指定了章节ID列表
	if len(FILTER_SECTION_IDS) > 0 {
		query = query.Where("id IN ?", FILTER_SECTION_IDS)
	}
	
	// 查询所有章节
	query.Find(&sections)
	
	// 过滤：只处理有足够评论数量且还没有总结的章节
	var needInitSections []model.CoursesSections
	for _, section := range sections {
		// 检查评论数量
		var commentCount int64
		model.Comment().Where("business_id = ? AND tenant_id = ?", section.ID, 1).Count(&commentCount)
		
		if int(commentCount) < MIN_COMMENT_COUNT {
			continue // 评论太少，跳过
		}
		
		if int(commentCount) > MAX_COMMENT_COUNT {
			t.Logf("⚠️  章节 '%s' (ID: %d) 有 %d 条评论，可能需要较长处理时间", section.Title, section.ID, commentCount)
		}
		
		// 检查是否已有总结
		var summaryCount int64
		model.CommentSummaryModel().Where("business_id = ? AND tenant_id = ?", section.ID, 1).Count(&summaryCount)
		
		if summaryCount > 0 {
			continue // 已有总结，跳过
		}
		
		needInitSections = append(needInitSections, section)
	}
	
	return needInitSections
}

// 处理单个文章的总结初始化
func processArticleSummary(t *testing.T, article model.Articles, summaryService *CommentSummaryService) error {
	t.Logf("🔄 处理文章: %s (ID: %d)", article.Title, article.ID)
	
	startTime := time.Now()
	summary, err := summaryService.GetSummary(article.ID, 0)
	duration := time.Since(startTime)
	
	if err != nil {
		return fmt.Errorf("生成总结失败: %v", err)
	}
	
	if summary == nil || summary.Summary == "" {
		return fmt.Errorf("生成的总结为空")
	}
	
	t.Logf("✅ 文章总结生成成功 (耗时: %v, 评论数: %d, 总结长度: %d字符)", 
		duration, summary.CommentCount, len(summary.Summary))
	
	return nil
}

// 处理单个章节的总结初始化
func processSectionSummary(t *testing.T, section model.CoursesSections, summaryService *CommentSummaryService) error {
	t.Logf("🔄 处理章节: %s (ID: %d)", section.Title, section.ID)
	
	startTime := time.Now()
	summary, err := summaryService.GetSummary(section.ID, 1)
	duration := time.Since(startTime)
	
	if err != nil {
		return fmt.Errorf("生成总结失败: %v", err)
	}
	
	if summary == nil || summary.Summary == "" {
		return fmt.Errorf("生成的总结为空")
	}
	
	t.Logf("✅ 章节总结生成成功 (耗时: %v, 评论数: %d, 总结长度: %d字符)", 
		duration, summary.CommentCount, len(summary.Summary))
	
	return nil
}

// TestInitializeCommentSummariesForArticles 初始化文章评论总结
func TestInitializeCommentSummariesForArticles(t *testing.T) {
	// 1. 初始化环境
	setupInitTestEnvironment(t)
	
	// 2. 获取需要初始化的文章
	articles := getArticlesNeedInitialization(t)
	
	if len(articles) == 0 {
		t.Logf("ℹ️  没有找到需要初始化总结的文章")
		t.Logf("   - 可能原因: 文章已有总结，或评论数量少于 %d 条", MIN_COMMENT_COUNT)
		return
	}
	
	t.Logf("📊 文章初始化统计:")
	t.Logf("   - 需要处理的文章数量: %d", len(articles))
	t.Logf("   - 最小评论数量要求: %d", MIN_COMMENT_COUNT)
	t.Logf("   - 批量处理大小: %d", BATCH_SIZE)
	t.Logf("   - 请求间延迟: %v", DELAY_BETWEEN_REQUESTS)
	t.Logf("")
	
	// 3. 初始化服务
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	summaryService := NewCommentSummaryService(ctx)
	
	// 4. 统计信息
	stats := &InitializationStats{
		StartTime: time.Now(),
	}
	
	// 5. 分批处理
	for i := 0; i < len(articles); i += BATCH_SIZE {
		end := i + BATCH_SIZE
		if end > len(articles) {
			end = len(articles)
		}
		
		batch := articles[i:end]
		t.Logf("🚀 处理批次 %d-%d (共 %d 篇文章)", i+1, end, len(articles))
		
		for j, article := range batch {
			stats.TotalProcessed++
			
			err := processArticleSummary(t, article, summaryService)
			if err != nil {
				stats.ErrorCount++
				errorMsg := fmt.Sprintf("文章 '%s' (ID: %d): %v", article.Title, article.ID, err)
				stats.Errors = append(stats.Errors, errorMsg)
				t.Logf("❌ %s", errorMsg)
			} else {
				stats.SuccessCount++
			}
			
			// 请求间延迟（最后一个请求不需要延迟）
			if j < len(batch)-1 || end < len(articles) {
				time.Sleep(DELAY_BETWEEN_REQUESTS)
			}
		}
		
		// 批次间的进度报告
		t.Logf("📈 批次完成 - 成功: %d, 失败: %d, 剩余: %d", 
			stats.SuccessCount, stats.ErrorCount, len(articles)-end)
		t.Logf("")
	}
	
	// 6. 输出最终统计
	stats.TotalDuration = time.Since(stats.StartTime)
	
	t.Logf("🏁 文章评论总结初始化完成!")
	t.Logf("═══════════════════════════════════════")
	t.Logf("📊 最终统计:")
	t.Logf("   - 总处理数量: %d", stats.TotalProcessed)
	t.Logf("   - 成功数量: %d", stats.SuccessCount)
	t.Logf("   - 失败数量: %d", stats.ErrorCount)
	t.Logf("   - 总耗时: %v", stats.TotalDuration)
	t.Logf("   - 平均耗时: %v", stats.TotalDuration/time.Duration(stats.TotalProcessed))
	
	if stats.ErrorCount > 0 {
		t.Logf("❌ 失败详情:")
		for _, err := range stats.Errors {
			t.Logf("   - %s", err)
		}
	}
	
	if stats.SuccessCount == 0 {
		t.Errorf("⚠️  警告: 没有成功初始化任何文章总结")
	}
}

// TestInitializeCommentSummariesForSections 初始化章节评论总结
func TestInitializeCommentSummariesForSections(t *testing.T) {
	// 1. 初始化环境
	setupInitTestEnvironment(t)
	
	// 2. 获取需要初始化的章节
	sections := getSectionsNeedInitialization(t)
	
	if len(sections) == 0 {
		t.Logf("ℹ️  没有找到需要初始化总结的章节")
		t.Logf("   - 可能原因: 章节已有总结，或评论数量少于 %d 条", MIN_COMMENT_COUNT)
		return
	}
	
	t.Logf("📊 章节初始化统计:")
	t.Logf("   - 需要处理的章节数量: %d", len(sections))
	t.Logf("   - 最小评论数量要求: %d", MIN_COMMENT_COUNT)
	t.Logf("   - 批量处理大小: %d", BATCH_SIZE)
	t.Logf("   - 请求间延迟: %v", DELAY_BETWEEN_REQUESTS)
	t.Logf("")
	
	// 3. 初始化服务
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	summaryService := NewCommentSummaryService(ctx)
	
	// 4. 统计信息
	stats := &InitializationStats{
		StartTime: time.Now(),
	}
	
	// 5. 分批处理
	for i := 0; i < len(sections); i += BATCH_SIZE {
		end := i + BATCH_SIZE
		if end > len(sections) {
			end = len(sections)
		}
		
		batch := sections[i:end]
		t.Logf("🚀 处理批次 %d-%d (共 %d 个章节)", i+1, end, len(sections))
		
		for j, section := range batch {
			stats.TotalProcessed++
			
			err := processSectionSummary(t, section, summaryService)
			if err != nil {
				stats.ErrorCount++
				errorMsg := fmt.Sprintf("章节 '%s' (ID: %d): %v", section.Title, section.ID, err)
				stats.Errors = append(stats.Errors, errorMsg)
				t.Logf("❌ %s", errorMsg)
			} else {
				stats.SuccessCount++
			}
			
			// 请求间延迟（最后一个请求不需要延迟）
			if j < len(batch)-1 || end < len(sections) {
				time.Sleep(DELAY_BETWEEN_REQUESTS)
			}
		}
		
		// 批次间的进度报告
		t.Logf("📈 批次完成 - 成功: %d, 失败: %d, 剩余: %d", 
			stats.SuccessCount, stats.ErrorCount, len(sections)-end)
		t.Logf("")
	}
	
	// 6. 输出最终统计
	stats.TotalDuration = time.Since(stats.StartTime)
	
	t.Logf("🏁 章节评论总结初始化完成!")
	t.Logf("═══════════════════════════════════════")
	t.Logf("📊 最终统计:")
	t.Logf("   - 总处理数量: %d", stats.TotalProcessed)
	t.Logf("   - 成功数量: %d", stats.SuccessCount)
	t.Logf("   - 失败数量: %d", stats.ErrorCount)
	t.Logf("   - 总耗时: %v", stats.TotalDuration)
	t.Logf("   - 平均耗时: %v", stats.TotalDuration/time.Duration(stats.TotalProcessed))
	
	if stats.ErrorCount > 0 {
		t.Logf("❌ 失败详情:")
		for _, err := range stats.Errors {
			t.Logf("   - %s", err)
		}
	}
	
	if stats.SuccessCount == 0 {
		t.Errorf("⚠️  警告: 没有成功初始化任何章节总结")
	}
}

// TestInitializeAllCommentSummaries 一次性初始化所有类型的评论总结
func TestInitializeAllCommentSummaries(t *testing.T) {
	t.Logf("🚀 开始初始化所有评论总结数据")
	t.Logf("═══════════════════════════════════════")
	
	// 1. 初始化文章总结
	t.Logf("📚 第一阶段: 初始化文章评论总结")
	t.Run("Articles", func(t *testing.T) {
		TestInitializeCommentSummariesForArticles(t)
	})
	
	t.Logf("")
	t.Logf("═══════════════════════════════════════")
	
	// 2. 初始化章节总结
	t.Logf("📖 第二阶段: 初始化章节评论总结")
	t.Run("Sections", func(t *testing.T) {
		TestInitializeCommentSummariesForSections(t)
	})
	
	t.Logf("")
	t.Logf("🎉 所有评论总结初始化完成!")
}

// TestDryRunInitialization 干跑模式 - 只统计不实际生成总结
func TestDryRunInitialization(t *testing.T) {
	// 初始化环境
	setupInitTestEnvironment(t)
	
	t.Logf("🔍 干跑模式 - 统计需要初始化的数据")
	t.Logf("═══════════════════════════════════════")
	
	// 统计文章
	articles := getArticlesNeedInitialization(t)
	t.Logf("📚 文章统计:")
	t.Logf("   - 需要初始化总结的文章: %d 篇", len(articles))
	
	if len(articles) > 0 {
		t.Logf("   - 前5篇文章示例:")
		for i, article := range articles {
			if i >= 5 {
				break
			}
			var commentCount int64
			model.Comment().Where("business_id = ? AND tenant_id = ?", article.ID, 0).Count(&commentCount)
			t.Logf("     %d. %s (ID: %d, 评论数: %d)", i+1, article.Title, article.ID, commentCount)
		}
		if len(articles) > 5 {
			t.Logf("     ... 还有 %d 篇文章", len(articles)-5)
		}
	}
	
	// 统计章节
	sections := getSectionsNeedInitialization(t)
	t.Logf("")
	t.Logf("📖 章节统计:")
	t.Logf("   - 需要初始化总结的章节: %d 个", len(sections))
	
	if len(sections) > 0 {
		t.Logf("   - 前5个章节示例:")
		for i, section := range sections {
			if i >= 5 {
				break
			}
			var commentCount int64
			model.Comment().Where("business_id = ? AND tenant_id = ?", section.ID, 1).Count(&commentCount)
			t.Logf("     %d. %s (ID: %d, 评论数: %d)", i+1, section.Title, section.ID, commentCount)
		}
		if len(sections) > 5 {
			t.Logf("     ... 还有 %d 个章节", len(sections)-5)
		}
	}
	
	// 估算处理时间
	totalItems := len(articles) + len(sections)
	estimatedTime := time.Duration(totalItems) * (5*time.Second + DELAY_BETWEEN_REQUESTS) // 假设每个总结平均5秒
	
	t.Logf("")
	t.Logf("⏱️  预估处理时间:")
	t.Logf("   - 总条目数: %d", totalItems)
	t.Logf("   - 预估总时长: %v", estimatedTime)
	t.Logf("   - 请求间延迟: %v", DELAY_BETWEEN_REQUESTS)
	t.Logf("   - 批次大小: %d", BATCH_SIZE)
	
	if totalItems == 0 {
		t.Logf("✅ 所有数据已完成初始化，无需处理")
	} else {
		t.Logf("")
		t.Logf("💡 运行建议:")
		t.Logf("   - 如需实际初始化，请运行对应的测试用例")
		t.Logf("   - 建议在业务低峰期进行初始化")
		t.Logf("   - 可以先从少量数据开始测试")
	}
}