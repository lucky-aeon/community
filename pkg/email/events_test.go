package email

import (
	"strings"
	"testing"
	"xhyovo.cn/community/pkg/log"
)

// TestCommentReplyEvent 测试评论回复事件
func TestCommentReplyEvent(t *testing.T) {
	// 初始化日志系统
	log.Init()

	tests := []struct {
		name     string
		template string
		expected bool
	}{
		{
			name:     "评论回复模板",
			template: "在 ${article.title}${course.title}${courses_section.title}，用户 ${user.name} 回复了你的评论 ${comment.content}",
			expected: true,
		},
		{
			name:     "其他模板",
			template: "你关注的用户 ${user.name} 发布了最新文章: ${article.title}",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCommentReplyEvent(tt.template)
			if result != tt.expected {
				t.Errorf("IsCommentReplyEvent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestAdoptionEvent 测试采纳事件
func TestAdoptionEvent(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expected bool
	}{
		{
			name:     "采纳模板",
			template: "在 ${article.title} 这篇文章中 ${comment.content} 该评论 \"被采纳\"",
			expected: true,
		},
		{
			name:     "其他模板",
			template: "用户 ${user.name} 在 ${article.title} 这篇文章中的评论 @ 了你",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsAdoptionEvent(tt.template)
			if result != tt.expected {
				t.Errorf("IsAdoptionEvent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestCourseUpdateEvent 测试课程更新事件
func TestCourseUpdateEvent(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expected bool
	}{
		{
			name:     "课程更新模板",
			template: "你订阅的课程  ${course.title} 更新了章节 ${courses_section.title}",
			expected: true,
		},
		{
			name:     "其他模板",
			template: "你关注的用户 ${user.name} 发布了最新文章: ${article.title}",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsCourseUpdateEvent(tt.template)
			if result != tt.expected {
				t.Errorf("IsCourseUpdateEvent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGenerateCommentReplyHTML 测试评论回复HTML生成
func TestGenerateCommentReplyHTML(t *testing.T) {
	testData := CommentReplyData{
		UserName:        "张三",
		UserAvatar:      "https://example.com/avatar.jpg",
		ReplyContent:    "感谢分享！这个解决方案很有帮助。",
		OriginalComment: "请问这个问题有好的解决方案吗？",
		ArticleTitle:    "Go语言并发编程最佳实践",
		ArticleURL:      "https://code.xhyovo.cn/article/view?articleId=123",
		ReplyTime:       "2025年01月04日 23:55",
	}

	html := GenerateCommentReplyHTML(testData)

	// 验证HTML包含必要的内容
	tests := []struct {
		name     string
		contains string
	}{
		{"包含用户名", "张三"},
		{"包含回复内容", "感谢分享！这个解决方案很有帮助"},
		{"包含原始评论", "请问这个问题有好的解决方案吗"},
		{"包含文章标题", "Go语言并发编程最佳实践"},
		{"包含回复时间", "2025年01月04日 23:55"},
		{"包含文章链接", "https://code.xhyovo.cn/article/view?articleId=123"},
		{"包含HTML结构", "<!DOCTYPE html>"},
		{"包含评论回复标识", "回复了你的评论"},
		{"包含查看按钮", "查看完整对话"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.contains) {
				t.Errorf("HTML不包含期望的内容: %s", tt.contains)
			}
		})
	}

	// 验证HTML长度合理
	if len(html) < 1000 {
		t.Error("生成的HTML长度过短")
	}
}

// TestGenerateAdoptionHTML 测试采纳HTML生成
func TestGenerateAdoptionHTML(t *testing.T) {
	testData := AdoptionData{
		ArticleTitle:   "如何优化MySQL查询性能",
		CommentContent: "可以通过添加索引、优化查询语句、调整数据库配置等方式来提升性能。",
		ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=456",
		AdoptionTime:   "2025年01月04日 23:55",
	}

	html := GenerateAdoptionHTML(testData)

	// 验证HTML包含必要的内容
	tests := []struct {
		name     string
		contains string
	}{
		{"包含文章标题", "如何优化MySQL查询性能"},
		{"包含评论内容", "可以通过添加索引、优化查询语句"},
		{"包含采纳时间", "2025年01月04日 23:55"},
		{"包含文章链接", "https://code.xhyovo.cn/article/view?articleId=456"},
		{"包含HTML结构", "<!DOCTYPE html>"},
		{"包含采纳标识", "恭喜！你的评论被采纳了"},
		{"包含采纳表情", "🎉"},
		{"包含查看按钮", "查看详情"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.contains) {
				t.Errorf("HTML不包含期望的内容: %s", tt.contains)
			}
		})
	}

	// 验证使用绿色主题
	if !strings.Contains(html, "#28a745") {
		t.Error("采纳邮件应该使用绿色主题")
	}
}

// TestGenerateCourseUpdateHTML 测试课程更新HTML生成
func TestGenerateCourseUpdateHTML(t *testing.T) {
	testData := CourseUpdateData{
		CourseTitle:  "Vue.js 3.0 全栈开发实战",
		SectionTitle: "第五章：Composition API 详解",
		CourseURL:    "https://code.xhyovo.cn/course/view?courseId=789",
		UpdateTime:   "2025年01月04日 23:55",
	}

	html := GenerateCourseUpdateHTML(testData)

	// 验证HTML包含必要的内容
	tests := []struct {
		name     string
		contains string
	}{
		{"包含课程标题", "Vue.js 3.0 全栈开发实战"},
		{"包含章节标题", "第五章：Composition API 详解"},
		{"包含更新时间", "2025年01月04日 23:55"},
		{"包含课程链接", "https://code.xhyovo.cn/course/view?courseId=789"},
		{"包含HTML结构", "<!DOCTYPE html>"},
		{"包含更新标识", "课程有新内容啦"},
		{"包含学习表情", "📚"},
		{"包含继续学习按钮", "继续学习"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.contains) {
				t.Errorf("HTML不包含期望的内容: %s", tt.contains)
			}
		})
	}

	// 验证使用橙色主题
	if !strings.Contains(html, "#fd7e14") {
		t.Error("课程更新邮件应该使用橙色主题")
	}
}

// TestHTMLEscaping 测试HTML转义
func TestHTMLEscaping(t *testing.T) {
	// 测试评论回复的HTML转义
	testData := CommentReplyData{
		UserName:        "<script>alert('xss')</script>",
		ReplyContent:    "包含<b>HTML</b>标签的回复",
		OriginalComment: "包含<script>的评论",
		ArticleTitle:    "测试<script>标题",
		ArticleURL:      "https://code.xhyovo.cn/article/view?articleId=1",
		ReplyTime:       "2025年01月04日 23:55",
	}

	html := GenerateCommentReplyHTML(testData)

	// 验证HTML标签被转义
	if strings.Contains(html, "<script>") {
		t.Error("HTML内容应该被转义以防止XSS攻击")
	}

	// 验证转义后的内容存在
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Error("HTML标签应该被正确转义")
	}
}

// TestEmptyFieldHandling 测试空字段处理
func TestEmptyFieldHandling(t *testing.T) {
	// 测试空头像处理
	testData := CommentReplyData{
		UserName:        "测试用户",
		UserAvatar:      "", // 空头像
		ReplyContent:    "测试回复",
		OriginalComment: "原始评论",
		ArticleTitle:    "测试文章",
		ArticleURL:      "https://code.xhyovo.cn/article/view?articleId=1",
		ReplyTime:       "2025年01月04日 23:55",
	}

	html := GenerateCommentReplyHTML(testData)

	// 应该生成默认头像
	if !strings.Contains(html, "via.placeholder.com") {
		t.Error("应该为空头像生成默认头像")
	}

	// 验证HTML仍然有效
	if len(html) < 500 {
		t.Error("即使有空字段，HTML也应该正常生成")
	}
}

// BenchmarkEventTemplateGeneration 性能测试
func BenchmarkEventTemplateGeneration(b *testing.B) {
	testData := CommentReplyData{
		UserName:        "性能测试用户",
		UserAvatar:      "https://example.com/avatar.jpg",
		ReplyContent:    "性能测试回复内容",
		OriginalComment: "性能测试原始评论",
		ArticleTitle:    "性能测试文章",
		ArticleURL:      "https://code.xhyovo.cn/article/view?articleId=1",
		ReplyTime:       "2025年01月04日 23:55",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateCommentReplyHTML(testData)
	}
}

// TestSendEventEmails 测试发送不同类型的事件邮件
func TestSendEventEmails(t *testing.T) {
	// 这个测试需要真实的SMTP配置，默认跳过
	t.Skip("跳过实际邮件发送测试")

	// 如果需要测试，可以设置环境变量
	testEmail := "test@example.com"
	if testEmail == "" {
		t.Skip("跳过邮件发送测试：缺少测试邮箱配置")
	}

	// 初始化邮件服务（需要真实配置）
	// Init("username", "password", "smtp.example.com:587")

	// 测试评论回复邮件
	replyData := CommentReplyData{
		UserName:        "测试用户",
		UserAvatar:      "",
		ReplyContent:    "这是一个测试回复",
		OriginalComment: "这是原始评论",
		ArticleTitle:    "测试文章",
		ArticleURL:      "https://code.xhyovo.cn/article/view?articleId=1",
		ReplyTime:       "2025年01月04日 23:55",
	}

	html := GenerateCommentReplyHTML(replyData)
	// err := Send([]string{testEmail}, html, "敲鸭社区 - 评论回复通知 (测试)")
	// if err != nil {
	//     t.Fatalf("发送评论回复邮件失败: %v", err)
	// }

	t.Logf("评论回复邮件HTML长度: %d", len(html))
}
