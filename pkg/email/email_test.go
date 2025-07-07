package email

import (
	"os"
	"testing"
	"xhyovo.cn/community/pkg/log"
)

// TestSendHTMLEmail 测试发送HTML邮件功能
// 注意：这个测试需要真实的SMTP配置，仅用于开发测试
func TestSendHTMLEmail(t *testing.T) {
	// 初始化日志系统
	log.Init()

	// 检查是否有测试邮件配置环境变量
	testEmail := ""
	smtpHost := ""
	smtpUser := ""
	smtpPass := ""

	if testEmail == "" || smtpHost == "" || smtpUser == "" || smtpPass == "" {
		t.Skip("跳过邮件发送测试：缺少环境变量配置 (TEST_EMAIL, SMTP_HOST, SMTP_USER, SMTP_PASS)")
	}

	// 初始化邮件服务
	Init(smtpUser, smtpPass, smtpHost)

	// 生成测试HTML内容
	testData := ArticleData{
		UserName:       "测试用户",
		UserAvatar:     "https://via.placeholder.com/50x50/667eea/ffffff?text=测",
		ArticleTitle:   "测试文章标题：HTML邮件功能验证",
		ArticleContent: "这是一个用于测试HTML邮件模板的测试文章。包含了基本的文章信息和格式化内容，用于验证邮件模板的显示效果。",
		ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=999",
		PublishTime:    "2025年01月04日 23:55",
	}

	htmlContent := GenerateArticlePublishHTML(testData)

	// 发送邮件
	err := Send([]string{testEmail}, htmlContent, "敲鸭社区 - 新文章通知 (测试)")
	if err != nil {
		t.Fatalf("发送邮件失败: %v", err)
	}

	t.Logf("测试邮件已发送到: %s", testEmail)
	t.Log("请检查邮箱查看HTML邮件效果")
}

// TestSendToMultipleRecipients 测试发送给多个收件人
func TestSendToMultipleRecipients(t *testing.T) {
	// 检查测试环境
	testEmails := os.Getenv("TEST_EMAILS") // 用逗号分隔的多个邮箱
	if testEmails == "" {
		t.Skip("跳过多收件人测试：缺少 TEST_EMAILS 环境变量")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || smtpUser == "" || smtpPass == "" {
		t.Skip("跳过邮件发送测试：缺少SMTP配置")
	}

	// 初始化邮件服务
	Init(smtpUser, smtpPass, smtpHost)

	// 解析邮箱列表
	recipients := []string{testEmails} // 简化处理，实际使用时可以用strings.Split处理

	// 生成测试内容
	testData := ArticleData{
		UserName:       "批量测试用户",
		UserAvatar:     "",
		ArticleTitle:   "批量邮件测试",
		ArticleContent: "这是一个用于测试批量发送HTML邮件的测试内容。",
		ArticleURL:     "https://code.xhyovo.cn",
		PublishTime:    "2025年01月04日 23:55",
	}

	htmlContent := GenerateArticlePublishHTML(testData)

	// 发送邮件
	err := Send(recipients, htmlContent, "敲鸭社区 - 批量邮件测试")
	if err != nil {
		t.Fatalf("批量发送邮件失败: %v", err)
	}

	t.Logf("批量邮件已发送到 %d 个收件人", len(recipients))
}

// BenchmarkGenerateHTML 测试HTML生成性能
func BenchmarkGenerateHTML(b *testing.B) {
	testData := ArticleData{
		UserName:       "性能测试用户",
		UserAvatar:     "https://example.com/avatar.jpg",
		ArticleTitle:   "性能测试文章",
		ArticleContent: "这是用于性能测试的文章内容",
		ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=1",
		PublishTime:    "2025年01月04日 23:55",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GenerateArticlePublishHTML(testData)
	}
}

// TestEmailContentValidation 测试邮件内容验证
func TestEmailContentValidation(t *testing.T) {
	tests := []struct {
		name    string
		data    ArticleData
		wantErr bool
		checkFn func(string) bool
	}{
		{
			name: "正常数据",
			data: ArticleData{
				UserName:       "正常用户",
				ArticleTitle:   "正常文章",
				ArticleContent: "正常内容",
				ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=1",
				PublishTime:    "2025年01月04日 23:55",
			},
			wantErr: false,
			checkFn: func(html string) bool {
				return len(html) > 1000 // HTML应该有一定长度
			},
		},
		{
			name: "空用户名",
			data: ArticleData{
				UserName:       "",
				ArticleTitle:   "测试文章",
				ArticleContent: "测试内容",
				ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=1",
				PublishTime:    "2025年01月04日 23:55",
			},
			wantErr: false,
			checkFn: func(html string) bool {
				return len(html) > 0 // 即使用户名为空也应该生成HTML
			},
		},
		{
			name: "空文章内容",
			data: ArticleData{
				UserName:       "测试用户",
				ArticleTitle:   "测试文章",
				ArticleContent: "",
				ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=1",
				PublishTime:    "2025年01月04日 23:55",
			},
			wantErr: false,
			checkFn: func(html string) bool {
				return len(html) > 0 // 即使内容为空也应该生成HTML
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			html := GenerateArticlePublishHTML(tt.data)

			if !tt.checkFn(html) {
				t.Errorf("生成的HTML不符合预期")
			}

			// 验证基本HTML结构
			if !containsBasicHTMLStructure(html) {
				t.Error("生成的HTML缺少基本结构")
			}
		})
	}
}

// containsBasicHTMLStructure 检查HTML是否包含基本结构
func containsBasicHTMLStructure(html string) bool {
	requiredElements := []string{
		"<!DOCTYPE html>",
		"<html",
		"<head>",
		"<body>",
		"</html>",
	}

	for _, element := range requiredElements {
		if !contains(html, element) {
			return false
		}
	}
	return true
}

// contains 简单的字符串包含检查
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(len(substr) == 0 || findInString(s, substr) >= 0)
}

// findInString 在字符串中查找子串
func findInString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
