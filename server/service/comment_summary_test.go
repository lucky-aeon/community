package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/config"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/mysql"
	"xhyovo.cn/community/server/model"
)

// ==================== 测试配置区域 ====================
// 请根据实际环境填写以下参数：

const (
	TEST_BUSINESS_ID = 59 // TODO: 请填入实际存在评论的文章ID
	TEST_TENANT_ID   = 1  // 租户类型：0=文章 1=章节 2=课程 3=分享会 4=AI日报
)

// 可选：如果需要测试特定文章，请修改上述ID
// 建议选择评论数量在5-20条之间的文章进行测试
// ====================================================

// 测试环境初始化
func setupTestEnvironment(t *testing.T) {
	// 初始化日志
	log.Init()

	// 设置中国时区
	chinaLoc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("加载时区失败: %v", err)
	}
	time.Local = chinaLoc

	// 初始化配置
	config.Init()
	appConfig := config.GetInstance()

	// 验证LLM配置
	if appConfig.LLMConfig.ApiKey == "" || appConfig.LLMConfig.Url == "" {
		t.Fatalf("LLM配置不完整，请检查环境变量: LLM_API_KEY, LLM_URL, LLM_MODEL")
	}

	// 初始化数据库连接
	db := appConfig.DbConfig
	mysql.Init(db.Username, db.Password, db.Address, db.Database)

	// 验证数据库连接
	dbInstance := mysql.GetInstance()
	if dbInstance == nil {
		t.Fatalf("数据库连接失败")
	}

	// 测试数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sqlDB, err := dbInstance.DB()
	if err != nil {
		t.Fatalf("获取数据库实例失败: %v", err)
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		t.Fatalf("数据库连接测试失败: %v", err)
	}

	t.Logf("✅ 测试环境初始化完成")
	t.Logf("   - 数据库连接: %s", db.Address)
	t.Logf("   - LLM服务: %s", appConfig.LLMConfig.Url)
	t.Logf("   - LLM模型: %s", appConfig.LLMConfig.Model)
}

// 获取租户类型名称
func getTenantTypeName(tenantId int) string {
	switch tenantId {
	case 0:
		return "文章"
	case 1:
		return "章节"
	case 2:
		return "课程"
	case 3:
		return "分享会"
	case 4:
		return "AI日报"
	default:
		return "未知类型"
	}
}

// 验证业务对象是否存在
func verifyBusinessObject(t *testing.T, businessId, tenantId int) (string, error) {
	switch tenantId {
	case 0: // 文章
		var article model.Articles
		result := model.Article().Where("id = ?", businessId).First(&article)
		if result.Error != nil {
			return "", result.Error
		}
		return article.Title, nil
	case 1: // 章节
		var section model.CoursesSections
		result := model.CoursesSection().Where("id = ?", businessId).First(&section)
		if result.Error != nil {
			return "", result.Error
		}
		return section.Title, nil
	case 2: // 课程
		var course model.Courses
		result := model.Course().Where("id = ?", businessId).First(&course)
		if result.Error != nil {
			return "", result.Error
		}
		return course.Title, nil
	case 3: // 分享会
		var meeting model.Meetings
		result := model.Meeting().Where("id = ?", businessId).First(&meeting)
		if result.Error != nil {
			return "", result.Error
		}
		return meeting.Title, nil
	case 4: // AI日报
		var aiNews model.AiNews
		result := model.AiNewsModel().Where("id = ?", businessId).First(&aiNews)
		if result.Error != nil {
			return "", result.Error
		}
		return aiNews.Title, nil
	default:
		return "", fmt.Errorf("不支持的租户类型: %d", tenantId)
	}
}

// 验证测试数据
func verifyTestData(t *testing.T) ([]model.Comments, string) {
	tenantTypeName := getTenantTypeName(TEST_TENANT_ID)
	
	// 验证业务对象是否存在
	businessTitle, err := verifyBusinessObject(t, TEST_BUSINESS_ID, TEST_TENANT_ID)
	if err != nil {
		t.Fatalf("%s ID %d 不存在: %v", tenantTypeName, TEST_BUSINESS_ID, err)
	}

	// 获取该业务对象的所有评论
	var comments []model.Comments
	result := model.Comment().Where("business_id = ? AND tenant_id = ?", TEST_BUSINESS_ID, TEST_TENANT_ID).Find(&comments)
	if result.Error != nil {
		t.Fatalf("查询评论失败: %v", result.Error)
	}

	if len(comments) == 0 {
		t.Fatalf("%s ID %d 下没有找到评论，请选择有评论的%s进行测试", tenantTypeName, TEST_BUSINESS_ID, tenantTypeName)
	}

	t.Logf("✅ 测试数据验证完成")
	t.Logf("   - %s标题: %s", tenantTypeName, businessTitle)
	t.Logf("   - %s ID: %d", tenantTypeName, TEST_BUSINESS_ID)
	t.Logf("   - 租户类型: %d (%s)", TEST_TENANT_ID, tenantTypeName)
	t.Logf("   - 评论数量: %d", len(comments))

	// 显示评论内容摘要
	for i, comment := range comments {
		if i >= 5 { // 最多显示5条评论摘要
			t.Logf("   - ... (还有%d条评论)", len(comments)-5)
			break
		}

		content := comment.Content
		if len(content) > 50 {
			content = content[:50] + "..."
		}
		t.Logf("   - 评论%d: %s", i+1, content)
	}

	return comments, businessTitle
}

// TestCommentSummaryWithRealData 使用真实数据测试AI评论总结功能
func TestCommentSummaryWithRealData(t *testing.T) {
	// 1. 初始化测试环境
	setupTestEnvironment(t)

	// 2. 验证测试数据
	comments, businessTitle := verifyTestData(t)

	// 3. 创建测试上下文和服务
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	summaryService := NewCommentSummaryService(ctx)

	t.Logf("🚀 开始AI评论总结测试...")

	// 4. 清理可能存在的旧总结（确保测试的独立性）
	model.CommentSummaryModel().Where("business_id = ? AND tenant_id = ?", TEST_BUSINESS_ID, TEST_TENANT_ID).Delete(&model.CommentSummary{})
	t.Logf("🧹 清理旧测试数据完成")

	// 5. 测试总结生成
	startTime := time.Now()
	summary, err := summaryService.GetSummary(TEST_BUSINESS_ID, TEST_TENANT_ID)
	duration := time.Since(startTime)

	if err != nil {
		t.Fatalf("❌ 生成评论总结失败: %v", err)
	}

	// 6. 验证总结结果
	if summary == nil {
		t.Fatalf("❌ 返回的总结为空")
	}

	if summary.Summary == "" {
		t.Fatalf("❌ 总结内容为空")
	}

	if summary.CommentCount != len(comments) {
		t.Errorf("❌ 评论数量统计错误: 期望 %d, 实际 %d", len(comments), summary.CommentCount)
	}

	if summary.BusinessId != TEST_BUSINESS_ID {
		t.Errorf("❌ 业务ID错误: 期望 %d, 实际 %d", TEST_BUSINESS_ID, summary.BusinessId)
	}

	if summary.TenantId != TEST_TENANT_ID {
		t.Errorf("❌ 租户ID错误: 期望 %d, 实际 %d", TEST_TENANT_ID, summary.TenantId)
	}

	// 7. 输出测试结果
	tenantTypeName := getTenantTypeName(TEST_TENANT_ID)
	t.Logf("✅ AI评论总结生成成功!")
	t.Logf("   - 处理时间: %v", duration)
	t.Logf("   - %s标题: %s", tenantTypeName, businessTitle)
	t.Logf("   - 评论数量: %d", summary.CommentCount)
	t.Logf("   - 总结长度: %d 字符", len(summary.Summary))
	t.Logf("   - 创建时间: %s", summary.CreatedAt)
	t.Logf("   - 更新时间: %s", summary.UpdatedAt)
	t.Logf("───────────────────────────────────────────────────────")
	t.Logf("📝 AI生成的总结内容:")
	t.Logf("%s", summary.Summary)
	t.Logf("───────────────────────────────────────────────────────")

	// 8. 验证总结质量（基本检查）
	if len(summary.Summary) < 50 {
		t.Logf("⚠️  警告: 总结内容较短（%d字符），可能质量不佳", len(summary.Summary))
	}

	if len(summary.Summary) > 1000 {
		t.Logf("⚠️  警告: 总结内容较长（%d字符），可能过于详细", len(summary.Summary))
	}

	// 9. 验证数据库保存
	var savedSummary model.CommentSummary
	result := model.CommentSummaryModel().Where("business_id = ? AND tenant_id = ?", TEST_BUSINESS_ID, TEST_TENANT_ID).First(&savedSummary)
	if result.Error != nil {
		t.Errorf("❌ 总结未正确保存到数据库: %v", result.Error)
	} else {
		t.Logf("✅ 总结已正确保存到数据库 (ID: %d)", savedSummary.ID)
	}
}

// TestCommentSummaryUpdate 测试评论总结更新机制
func TestCommentSummaryUpdate(t *testing.T) {
	// 1. 初始化测试环境
	setupTestEnvironment(t)

	// 2. 验证测试数据
	_, _ = verifyTestData(t)

	// 3. 创建测试上下文和服务
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	summaryService := NewCommentSummaryService(ctx)

	t.Logf("🔄 开始测试评论总结更新机制...")

	// 4. 确保已有总结存在（先生成一次）
	summary1, err := summaryService.GetSummary(TEST_BUSINESS_ID, TEST_TENANT_ID)
	if err != nil {
		t.Fatalf("❌ 初始总结生成失败: %v", err)
	}

	t.Logf("✅ 初始总结已存在 (评论数: %d)", summary1.CommentCount)

	// 5. 测试更新判断机制
	summaryService.UpdateSummaryIfNeeded(TEST_BUSINESS_ID, TEST_TENANT_ID)
	t.Logf("✅ 更新检查完成")

	// 6. 再次获取总结，验证是否有变化
	summary2, err := summaryService.GetSummary(TEST_BUSINESS_ID, TEST_TENANT_ID)
	if err != nil {
		t.Fatalf("❌ 二次获取总结失败: %v", err)
	}

	// 7. 比较两次结果
	if summary1.ID == summary2.ID && summary1.UpdatedAt == summary2.UpdatedAt {
		t.Logf("✅ 更新机制正常: 无变化时不重复生成")
	} else {
		t.Logf("ℹ️  总结发生了更新")
		t.Logf("   - 更新前: %s", summary1.UpdatedAt)
		t.Logf("   - 更新后: %s", summary2.UpdatedAt)
	}
}

// TestCommentSummaryPerformance 性能测试
func TestCommentSummaryPerformance(t *testing.T) {
	// 1. 初始化测试环境
	setupTestEnvironment(t)

	// 2. 验证测试数据
	comments, _ := verifyTestData(t)

	// 3. 创建测试上下文和服务
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	summaryService := NewCommentSummaryService(ctx)

	t.Logf("⚡ 开始性能测试...")

	// 4. 多次调用测试
	const testRounds = 3
	var totalDuration time.Duration

	for i := 0; i < testRounds; i++ {
		startTime := time.Now()
		_, err := summaryService.GetSummary(TEST_BUSINESS_ID, TEST_TENANT_ID)
		duration := time.Since(startTime)
		totalDuration += duration

		if err != nil {
			t.Errorf("❌ 第%d次调用失败: %v", i+1, err)
			continue
		}

		t.Logf("   - 第%d次调用: %v", i+1, duration)
	}

	avgDuration := totalDuration / testRounds
	t.Logf("✅ 性能测试完成")
	t.Logf("   - 评论数量: %d", len(comments))
	t.Logf("   - 平均耗时: %v", avgDuration)
	t.Logf("   - 总测试时间: %v", totalDuration)

	// 性能基准检查
	if avgDuration > 30*time.Second {
		t.Logf("⚠️  警告: 平均响应时间较慢（%v），建议优化", avgDuration)
	} else if avgDuration < 2*time.Second {
		t.Logf("✅ 响应时间优秀（%v）", avgDuration)
	} else {
		t.Logf("✅ 响应时间良好（%v）", avgDuration)
	}
}

// TestCommentSummaryErrorHandling 错误处理测试
func TestCommentSummaryErrorHandling(t *testing.T) {
	// 1. 初始化测试环境
	setupTestEnvironment(t)

	// 2. 创建测试上下文和服务
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	summaryService := NewCommentSummaryService(ctx)

	t.Logf("🚨 开始错误处理测试...")

	// 3. 测试不存在的文章ID
	nonExistentID := 999999999
	summary, err := summaryService.GetSummary(nonExistentID, TEST_TENANT_ID)

	if err == nil && summary == nil {
		t.Logf("✅ 正确处理了不存在的文章ID（返回空结果）")
	} else if err != nil {
		t.Logf("✅ 正确返回了错误: %v", err)
	} else {
		t.Errorf("❌ 未正确处理不存在的文章ID，返回了: %+v", summary)
	}

	// 4. 测试无效的租户ID
	invalidTenantID := 999
	summary, err = summaryService.GetSummary(TEST_BUSINESS_ID, invalidTenantID)

	if err != nil {
		t.Logf("✅ 正确处理了无效的租户ID: %v", err)
	} else if summary == nil || summary.Summary == "" {
		t.Logf("✅ 正确处理了无效的租户ID（返回空结果）")
	} else {
		t.Logf("ℹ️  对无效租户ID有容错处理")
	}

	t.Logf("✅ 错误处理测试完成")
}

// TestCommentSummaryEmptyComments 测试空评论情况
func TestCommentSummaryEmptyComments(t *testing.T) {
	// 1. 初始化测试环境
	setupTestEnvironment(t)
	
	// 2. 创建测试上下文和服务
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	summaryService := NewCommentSummaryService(ctx)
	
	t.Logf("🔍 开始空评论测试...")
	
	// 3. 使用一个不太可能有评论的大ID进行测试
	testBusinessId := 999999999
	testTenantId := 0
	
	// 4. 确保没有评论数据
	var commentCount int64
	model.Comment().Where("business_id = ? AND tenant_id = ?", testBusinessId, testTenantId).Count(&commentCount)
	
	if commentCount > 0 {
		t.Logf("⚠️  测试ID %d 下有 %d 条评论，选择其他ID进行测试", testBusinessId, commentCount)
		// 可以选择删除这些评论或选择其他ID
		return
	}
	
	t.Logf("✅ 确认测试业务对象 (ID: %d, TenantID: %d) 无评论", testBusinessId, testTenantId)
	
	// 5. 测试获取总结
	summary, err := summaryService.GetSummary(testBusinessId, testTenantId)
	
	// 6. 验证结果
	if err != nil {
		t.Errorf("❌ 空评论情况不应该返回错误: %v", err)
		return
	}
	
	if summary != nil {
		t.Errorf("❌ 空评论情况应该返回nil，实际返回: %+v", summary)
		return
	}
	
	t.Logf("✅ 空评论情况处理正确 - 返回nil而不是错误")
	
	// 7. 测试UpdateSummaryIfNeeded是否也能正确处理
	summaryService.UpdateSummaryIfNeeded(testBusinessId, testTenantId)
	t.Logf("✅ UpdateSummaryIfNeeded对空评论的处理完成")
	
	// 8. 再次验证没有创建错误的总结记录
	var summaryCount int64
	model.CommentSummaryModel().Where("business_id = ? AND tenant_id = ?", testBusinessId, testTenantId).Count(&summaryCount)
	
	if summaryCount > 0 {
		t.Errorf("❌ 空评论情况不应该创建总结记录")
	} else {
		t.Logf("✅ 没有为空评论创建总结记录")
	}
	
	t.Logf("✅ 空评论测试完成")
}

// TestCommentSummaryMultipleTenantTypes 测试多种租户类型
func TestCommentSummaryMultipleTenantTypes(t *testing.T) {
	// 1. 初始化测试环境
	setupTestEnvironment(t)
	
	t.Logf("🔄 开始多租户类型测试...")
	
	// 定义要测试的租户类型和对应的业务ID
	// 用户需要根据实际数据填写这些值
	testCases := []struct {
		businessId int
		tenantId   int
		skip       bool // 如果没有对应类型的测试数据，可以跳过
		reason     string
	}{
		{TEST_BUSINESS_ID, 0, false, ""}, // 文章
		{1, 1, true, "请设置有评论的章节ID"},     // 章节 - 需要用户填写
		{1, 2, true, "请设置有评论的课程ID"},     // 课程 - 需要用户填写
		{1, 3, true, "请设置有评论的分享会ID"},    // 分享会 - 需要用户填写
		{1, 4, true, "请设置有评论的AI日报ID"},   // AI日报 - 需要用户填写
	}
	
	gin.SetMode(gin.TestMode)
	ctx := &gin.Context{}
	summaryService := NewCommentSummaryService(ctx)
	
	successCount := 0
	skipCount := 0
	
	for _, tc := range testCases {
		tenantTypeName := getTenantTypeName(tc.tenantId)
		
		if tc.skip {
			t.Logf("⏭️  跳过 %s 测试: %s", tenantTypeName, tc.reason)
			skipCount++
			continue
		}
		
		t.Logf("🧪 测试租户类型: %s (ID: %d, TenantID: %d)", tenantTypeName, tc.businessId, tc.tenantId)
		
		// 验证业务对象存在
		businessTitle, err := verifyBusinessObject(t, tc.businessId, tc.tenantId)
		if err != nil {
			t.Logf("❌ %s ID %d 不存在，跳过: %v", tenantTypeName, tc.businessId, err)
			continue
		}
		
		// 检查是否有评论
		var commentCount int64
		model.Comment().Where("business_id = ? AND tenant_id = ?", tc.businessId, tc.tenantId).Count(&commentCount)
		if commentCount == 0 {
			t.Logf("⚠️  %s '%s' (ID: %d) 没有评论，跳过", tenantTypeName, businessTitle, tc.businessId)
			continue
		}
		
		// 清理可能存在的旧总结
		model.CommentSummaryModel().Where("business_id = ? AND tenant_id = ?", tc.businessId, tc.tenantId).Delete(&model.CommentSummary{})
		
		// 生成总结
		startTime := time.Now()
		summary, err := summaryService.GetSummary(tc.businessId, tc.tenantId)
		duration := time.Since(startTime)
		
		if err != nil {
			t.Errorf("❌ %s总结生成失败: %v", tenantTypeName, err)
			continue
		}
		
		if summary == nil || summary.Summary == "" {
			t.Errorf("❌ %s总结内容为空", tenantTypeName)
			continue
		}
		
		// 输出结果
		t.Logf("✅ %s总结生成成功!", tenantTypeName)
		t.Logf("   - %s: %s", tenantTypeName, businessTitle)
		t.Logf("   - 评论数量: %d", summary.CommentCount)
		t.Logf("   - 总结长度: %d 字符", len(summary.Summary))
		t.Logf("   - 处理时间: %v", duration)
		t.Logf("   - 总结预览: %s", func() string {
			if len(summary.Summary) > 100 {
				return summary.Summary[:100] + "..."
			}
			return summary.Summary
		}())
		t.Logf("───────────────────────────────────────")
		
		successCount++
	}
	
	t.Logf("🏁 多租户类型测试完成")
	t.Logf("   - 成功测试: %d 种类型", successCount)
	t.Logf("   - 跳过测试: %d 种类型", skipCount)
	t.Logf("   - 总计类型: %d 种", len(testCases))
	
	if successCount == 0 {
		t.Logf("⚠️  警告: 没有成功测试任何租户类型，请检查测试数据配置")
	}
}
