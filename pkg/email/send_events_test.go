package email

import (
	"strings"
	"testing"
	"xhyovo.cn/community/pkg/log"
)

// 邮件配置 - 实际使用时应从环境变量获取
const (
	testEmailAddr = ""
	testSMTPHost  = ""
	testSMTPUser  = ""
	testSMTPPass  = ""
)

// TestSendCommentReplyEmail 测试发送评论回复邮件
func TestSendCommentReplyEmail(t *testing.T) {
	// 初始化日志系统
	log.Init()

	// 初始化邮件服务
	Init(testSMTPUser, testSMTPPass, testSMTPHost)

	// 创建评论回复测试数据
	testData := CommentReplyData{
		UserName:        "张三",
		UserAvatar:      "https://via.placeholder.com/50x50/28a745/ffffff?text=张",
		ReplyContent:    "感谢你的问题！这个问题确实很重要。我建议你可以先查看官方文档，然后结合实际项目练习。如果还有疑问，随时可以继续讨论。",
		OriginalComment: "请问Go语言的并发编程有什么好的学习资料推荐吗？特别是关于goroutine和channel的使用。",
		ArticleTitle:    "Go语言并发编程最佳实践：从入门到精通",
		ArticleURL:      "https://code.xhyovo.cn/article/view?articleId=123",
		ReplyTime:       "2025年01月05日 00:15",
	}

	// 生成HTML邮件内容
	htmlContent := GenerateCommentReplyHTML(testData)

	// 发送邮件
	err := Send([]string{testEmailAddr}, htmlContent, "敲鸭社区 - 评论回复通知 (测试)")
	if err != nil {
		t.Fatalf("发送评论回复邮件失败: %v", err)
	}

	t.Logf("评论回复邮件已发送到: %s", testEmailAddr)
	t.Log("请检查邮箱查看HTML邮件效果")
}

// TestSendAdoptionEmail 测试发送采纳邮件
func TestSendAdoptionEmail(t *testing.T) {
	// 初始化日志系统
	log.Init()

	// 初始化邮件服务
	Init(testSMTPUser, testSMTPPass, testSMTPHost)

	// 创建采纳测试数据
	testData := AdoptionData{
		ArticleTitle:   "如何优化MySQL数据库查询性能？",
		CommentContent: "可以通过以下几个方面来优化MySQL查询性能：1. 合理设计索引，避免全表扫描；2. 优化SQL语句，避免使用SELECT *；3. 分析执行计划，找出性能瓶颈；4. 考虑分库分表策略；5. 定期维护和优化数据库配置。希望这些建议对你有帮助！",
		ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=456",
		AdoptionTime:   "2025年01月05日 00:20",
	}

	// 生成HTML邮件内容
	htmlContent := GenerateAdoptionHTML(testData)

	// 发送邮件
	err := Send([]string{testEmailAddr}, htmlContent, "敲鸭社区 - 评论采纳通知 (测试)")
	if err != nil {
		t.Fatalf("发送采纳邮件失败: %v", err)
	}

	t.Logf("采纳邮件已发送到: %s", testEmailAddr)
	t.Log("请检查邮箱查看HTML邮件效果")
}

// TestSendCourseUpdateEmail 测试发送课程更新邮件
func TestSendCourseUpdateEmail(t *testing.T) {
	// 初始化日志系统
	log.Init()

	// 初始化邮件服务
	Init(testSMTPUser, testSMTPPass, testSMTPHost)

	// 创建课程更新测试数据
	testData := CourseUpdateData{
		CourseTitle:  "Vue.js 3.0 全栈开发实战：从零基础到项目上线",
		SectionTitle: "第八章：Composition API 深度解析与实战应用",
		CourseURL:    "https://code.xhyovo.cn/course/view?courseId=789",
		UpdateTime:   "2025年01月05日 00:25",
	}

	// 生成HTML邮件内容
	htmlContent := GenerateCourseUpdateHTML(testData)

	// 发送邮件
	err := Send([]string{testEmailAddr}, htmlContent, "敲鸭社区 - 课程更新通知 (测试)")
	if err != nil {
		t.Fatalf("发送课程更新邮件失败: %v", err)
	}

	t.Logf("课程更新邮件已发送到: %s", testEmailAddr)
	t.Log("请检查邮箱查看HTML邮件效果")
}

// TestSendArticleCommentEmail 测试发送文章评论邮件
func TestSendArticleCommentEmail(t *testing.T) {
	// 初始化日志系统
	log.Init()

	// 初始化邮件服务
	Init(testSMTPUser, testSMTPPass, testSMTPHost)

	// 创建文章评论测试数据（使用评论回复模板，但内容不同）
	testData := CommentReplyData{
		UserName:        "李四",
		UserAvatar:      "https://via.placeholder.com/50x50/6c757d/ffffff?text=李",
		ReplyContent:    "这篇文章写得很好！特别是关于性能优化的部分，让我学到了很多新的技巧。期待作者能够分享更多类似的实战经验。",
		OriginalComment: "你订阅的文章有新评论",
		ArticleTitle:    "Python Web开发性能优化指南",
		ArticleURL:      "https://code.xhyovo.cn/article/view?articleId=321",
		ReplyTime:       "2025年01月05日 00:30",
	}

	// 生成HTML邮件内容（复用评论回复模板）
	htmlContent := GenerateCommentReplyHTML(testData)

	// 发送邮件
	err := Send([]string{testEmailAddr}, htmlContent, "敲鸭社区 - 文章评论通知 (测试)")
	if err != nil {
		t.Fatalf("发送文章评论邮件失败: %v", err)
	}

	t.Logf("文章评论邮件已发送到: %s", testEmailAddr)
	t.Log("请检查邮箱查看HTML邮件效果")
}

// TestSendAllEventTypes 测试发送所有事件类型的邮件（批量测试）
func TestSendAllEventTypes(t *testing.T) {
	// 初始化日志系统
	log.Init()

	// 初始化邮件服务
	Init(testSMTPUser, testSMTPPass, testSMTPHost)

	t.Log("开始批量测试所有事件类型的邮件发送...")

	// 1. 用户发布文章
	t.Run("UserUpdate", func(t *testing.T) {
		testData := ArticleData{
			UserName:       "王五",
			UserAvatar:     "https://via.placeholder.com/50x50/dc3545/ffffff?text=王",
			ArticleTitle:   "深入理解JavaScript闭包：原理、应用与性能优化",
			ArticleContent: "闭包是JavaScript中最重要的概念之一，它不仅影响着代码的结构和功能，还直接关系到应用的性能。本文将从原理出发，结合实际案例，深入探讨闭包的各种应用场景，并提供性能优化的最佳实践。无论你是初学者还是有经验的开发者，都能从中获得有价值的见解。",
			ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=999",
			PublishTime:    "2025年01月05日 00:35",
		}

		htmlContent := GenerateArticlePublishHTML(testData)
		err := Send([]string{testEmailAddr}, htmlContent, "敲鸭社区 - 用户发布文章通知 (批量测试)")
		if err != nil {
			t.Errorf("发送用户发布文章邮件失败: %v", err)
		} else {
			t.Log("✓ 用户发布文章邮件发送成功")
		}
	})

	// 2. 评论回复
	t.Run("CommentReply", func(t *testing.T) {
		testData := CommentReplyData{
			UserName:        "赵六",
			UserAvatar:      "https://via.placeholder.com/50x50/007bff/ffffff?text=赵",
			ReplyContent:    "你说得很对！我也遇到过类似的问题，这个解决方案确实很实用。",
			OriginalComment: "这个功能的实现原理是什么？",
			ArticleTitle:    "React Hooks 最佳实践",
			ArticleURL:      "https://code.xhyovo.cn/article/view?articleId=888",
			ReplyTime:       "2025年01月05日 00:40",
		}

		htmlContent := GenerateCommentReplyHTML(testData)
		err := Send([]string{testEmailAddr}, htmlContent, "敲鸭社区 - 评论回复通知 (批量测试)")
		if err != nil {
			t.Errorf("发送评论回复邮件失败: %v", err)
		} else {
			t.Log("✓ 评论回复邮件发送成功")
		}
	})

	// 3. 评论采纳
	t.Run("Adoption", func(t *testing.T) {
		testData := AdoptionData{
			ArticleTitle:   "Docker容器化部署最佳实践",
			CommentContent: "建议使用多阶段构建来减小镜像大小，同时配置合适的健康检查和资源限制。具体可以参考官方文档的production checklist。",
			ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=777",
			AdoptionTime:   "2025年01月05日 00:45",
		}

		htmlContent := GenerateAdoptionHTML(testData)
		err := Send([]string{testEmailAddr}, htmlContent, "敲鸭社区 - 评论采纳通知 (批量测试)")
		if err != nil {
			t.Errorf("发送评论采纳邮件失败: %v", err)
		} else {
			t.Log("✓ 评论采纳邮件发送成功")
		}
	})

	// 4. 课程更新
	t.Run("CourseUpdate", func(t *testing.T) {
		testData := CourseUpdateData{
			CourseTitle:  "Spring Boot 微服务架构实战",
			SectionTitle: "第十二章：分布式事务处理与Saga模式",
			CourseURL:    "https://code.xhyovo.cn/course/view?courseId=666",
			UpdateTime:   "2025年01月05日 00:50",
		}

		htmlContent := GenerateCourseUpdateHTML(testData)
		err := Send([]string{testEmailAddr}, htmlContent, "敲鸭社区 - 课程更新通知 (批量测试)")
		if err != nil {
			t.Errorf("发送课程更新邮件失败: %v", err)
		} else {
			t.Log("✓ 课程更新邮件发送成功")
		}
	})

	t.Log("批量邮件测试完成！请检查邮箱查看所有邮件效果")
}

// TestSendEventEmailsToMultipleRecipients 测试发送事件邮件给多个收件人
func TestSendEventEmailsToMultipleRecipients(t *testing.T) {
	// 跳过多收件人测试，避免发送过多邮件
	t.Skip("跳过多收件人测试以避免邮件发送过量")

	// 初始化日志系统
	log.Init()

	// 初始化邮件服务
	Init(testSMTPUser, testSMTPPass, testSMTPHost)

	// 多个收件人列表（实际使用时应该是真实的邮箱地址）
	recipients := []string{
		testEmailAddr,
		// "test2@example.com",
		// "test3@example.com",
	}

	testData := ArticleData{
		UserName:       "群发测试用户",
		UserAvatar:     "",
		ArticleTitle:   "群发邮件测试",
		ArticleContent: "这是一个用于测试群发邮件功能的测试内容。",
		ArticleURL:     "https://code.xhyovo.cn",
		PublishTime:    "2025年01月05日 01:00",
	}

	htmlContent := GenerateArticlePublishHTML(testData)
	err := Send(recipients, htmlContent, "敲鸭社区 - 群发邮件测试")
	if err != nil {
		t.Fatalf("群发邮件失败: %v", err)
	}

	t.Logf("群发邮件已发送到 %d 个收件人", len(recipients))
}

// TestEmailContentPreview 测试邮件内容预览（不实际发送）
func TestEmailContentPreview(t *testing.T) {
	// 这个测试不发送邮件，只生成内容用于预览

	testCases := []struct {
		name     string
		generate func() string
	}{
		{
			name: "文章发布邮件",
			generate: func() string {
				data := ArticleData{
					UserName:       "预览用户",
					UserAvatar:     "",
					ArticleTitle:   "预览文章标题",
					ArticleContent: "预览文章内容...",
					ArticleURL:     "https://example.com",
					PublishTime:    "2025年01月05日",
				}
				return GenerateArticlePublishHTML(data)
			},
		},
		{
			name: "评论回复邮件",
			generate: func() string {
				data := CommentReplyData{
					UserName:        "回复用户",
					UserAvatar:      "",
					ReplyContent:    "这是回复内容",
					OriginalComment: "原始评论内容",
					ArticleTitle:    "文章标题",
					ArticleURL:      "https://example.com",
					ReplyTime:       "2025年01月05日",
				}
				return GenerateCommentReplyHTML(data)
			},
		},
		{
			name: "采纳邮件",
			generate: func() string {
				data := AdoptionData{
					ArticleTitle:   "文章标题",
					CommentContent: "被采纳的评论内容",
					ArticleURL:     "https://example.com",
					AdoptionTime:   "2025年01月05日",
				}
				return GenerateAdoptionHTML(data)
			},
		},
		{
			name: "课程更新邮件",
			generate: func() string {
				data := CourseUpdateData{
					CourseTitle:  "课程标题",
					SectionTitle: "章节标题",
					CourseURL:    "https://example.com",
					UpdateTime:   "2025年01月05日",
				}
				return GenerateCourseUpdateHTML(data)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			html := tc.generate()

			// 验证HTML长度
			if len(html) < 500 {
				t.Errorf("%s 生成的HTML长度过短: %d", tc.name, len(html))
			}

			// 验证HTML包含基本结构
			if !strings.Contains(html, "<!DOCTYPE html>") {
				t.Errorf("%s 缺少HTML声明", tc.name)
			}

			if !strings.Contains(html, "敲鸭社区") {
				t.Errorf("%s 缺少品牌标识", tc.name)
			}

			t.Logf("%s HTML长度: %d 字符", tc.name, len(html))
		})
	}
}
