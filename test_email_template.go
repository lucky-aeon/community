package main

import (
	"fmt"
	"xhyovo.cn/community/pkg/email"
	"xhyovo.cn/community/pkg/log"
)

func main() {
	// 初始化日志
	log.Init()

	// 创建测试数据
	testData := email.ArticleData{
		UserName:       "测试用户",
		UserAvatar:     "", // 不设置头像，避免防盗链问题
		ArticleTitle:   "测试文章标题：HTML邮件功能验证",
		ArticleContent: "这是一个用于测试HTML邮件模板的测试文章。包含了基本的文章信息和格式化内容，用于验证邮件模板的显示效果。",
		ArticleURL:     "https://code.xhyovo.cn/article/view?articleId=999",
		PublishTime:    "2025年01月04日 23:55",
	}

	// 生成HTML
	htmlContent := email.GenerateArticlePublishHTML(testData)

	fmt.Println("HTML邮件模板已生成")
	fmt.Printf("HTML长度: %d 字符\n", len(htmlContent))

	// 显示HTML的前500个字符以检查格式
	if len(htmlContent) > 500 {
		fmt.Println("HTML前500字符:")
		fmt.Println(htmlContent[:500])
	} else {
		fmt.Println("完整HTML:")
		fmt.Println(htmlContent)
	}

	// 检查是否包含测试数据
	if len(htmlContent) > 0 {
		fmt.Println("✓ HTML生成成功")
	} else {
		fmt.Println("✗ HTML生成失败")
	}

	// 检查关键内容是否存在
	checks := []string{
		"测试用户",
		"测试文章标题",
		"这是一个用于测试HTML邮件模板",
		"2025年01月04日 23:55",
		"https://code.xhyovo.cn/article/view?articleId=999",
	}

	for _, check := range checks {
		if contains(htmlContent, check) {
			fmt.Printf("✓ 包含: %s\n", check)
		} else {
			fmt.Printf("✗ 缺少: %s\n", check)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findInString(s, substr) >= 0
}

func findInString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
