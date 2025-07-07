package email

import (
	"strings"
	"testing"
)

func TestGenerateArticlePublishHTML(t *testing.T) {
	// 测试数据
	testData := ArticleData{
		UserName:       "张三",
		UserAvatar:     "https://example.com/avatar.jpg",
		ArticleTitle:   "聊聊最近 AI 相关的感悟：AI IDE | Agent | 解决方案 | 未来方向",
		ArticleContent: "随着人工智能技术的快速发展，我们正在见证一个前所未有的变革时代。从智能IDE的出现到AI Agent的广泛应用，从创新解决方案的涌现到对未来发展方向的思考，每一个环节都在深刻地改变着我们的工作方式和思维模式。",
		ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=225",
		PublishTime:    "2025年01月04日 23:55",
	}

	// 生成HTML
	html := GenerateArticlePublishHTML(testData)

	// 验证HTML包含必要的内容
	tests := []struct {
		name     string
		contains string
	}{
		{"包含用户名", "张三"},
		{"包含文章标题", "聊聊最近 AI 相关的感悟"},
		{"包含文章内容", "随着人工智能技术的快速发展"},
		{"包含发布时间", "2025年01月04日 23:55"},
		{"包含文章链接", "https://code.xhyovo.cn/article/view?articleId=225"},
		{"包含头像链接", "https://example.com/avatar.jpg"},
		{"包含HTML结构", "<!DOCTYPE html>"},
		{"包含CSS样式", "<style>"},
		{"包含按钮", "查看完整文章"},
		{"包含品牌标识", "敲鸭社区"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(html, tt.contains) {
				t.Errorf("HTML不包含期望的内容: %s", tt.contains)
			}
		})
	}

	// 验证HTML结构完整性
	if !strings.Contains(html, "<html") || !strings.Contains(html, "</html>") {
		t.Error("HTML结构不完整")
	}

	// 验证响应式设计
	if !strings.Contains(html, "@media") {
		t.Error("缺少响应式CSS")
	}

	// 打印生成的HTML用于手动检查
	t.Logf("生成的HTML长度: %d 字符", len(html))
}

func TestGenerateArticlePublishHTML_WithLongContent(t *testing.T) {
	// 测试长内容截取
	longContent := strings.Repeat("这是一段很长的文章内容。", 50) // 生成超过200字的内容

	testData := ArticleData{
		UserName:       "李四",
		UserAvatar:     "",
		ArticleTitle:   "测试长文章",
		ArticleContent: longContent,
		ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=1",
		PublishTime:    "2025年01月04日 12:00",
	}

	html := GenerateArticlePublishHTML(testData)

	// 验证内容被截取并添加了省略号
	if !strings.Contains(html, "...") {
		t.Error("长内容应该被截取并添加省略号")
	}

	// 验证默认头像
	if !strings.Contains(html, "via.placeholder.com") {
		t.Error("应该使用默认头像")
	}
}

func TestGenerateArticlePublishHTML_HTMLEscape(t *testing.T) {
	// 测试HTML转义
	testData := ArticleData{
		UserName:       "<script>alert('xss')</script>",
		UserAvatar:     "https://example.com/avatar.jpg",
		ArticleTitle:   "测试<script>标签",
		ArticleContent: "包含<b>HTML</b>标签的内容",
		ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=1",
		PublishTime:    "2025年01月04日 12:00",
	}

	html := GenerateArticlePublishHTML(testData)

	// 验证HTML标签被转义
	if strings.Contains(html, "<script>") {
		t.Error("HTML内容应该被转义以防止XSS攻击")
	}

	// 验证转义后的内容存在
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Error("HTML标签应该被正确转义")
	}
}

func TestIsUserUpdateEvent(t *testing.T) {
	tests := []struct {
		name     string
		template string
		expected bool
	}{
		{
			name:     "用户发布文章模板",
			template: "你关注的用户 ${user.name} 发布了最新文章: ${article.title}",
			expected: true,
		},
		{
			name:     "评论模板",
			template: "用户 ${user.name} 在 ${article.title} 这篇文章中的评论 @ 了你",
			expected: false,
		},
		{
			name:     "其他模板",
			template: "在 ${article.title} 这篇文章中 ${comment.content} 该评论 \"被采纳\"",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsUserUpdateEvent(tt.template)
			if result != tt.expected {
				t.Errorf("IsUserUpdateEvent() = %v, expected %v", result, tt.expected)
			}
		})
	}
}
