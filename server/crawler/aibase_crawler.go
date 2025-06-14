package crawler

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	localTime "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
)

// AIBaseCrawler AIBase网站爬虫
type AIBaseCrawler struct {
	name      string
	baseURL   string
	enabled   bool
	startID   int
	maxCount  int // 最大爬取数量，0表示无限制
	userAgent string
}

// NewAIBaseCrawler 创建AIBase爬虫
func NewAIBaseCrawler(startID int) *AIBaseCrawler {
	return &AIBaseCrawler{
		name:      "AIBase每日资讯",
		baseURL:   "https://www.aibase.com/zh/daily/",
		enabled:   true,
		startID:   startID,
		maxCount:  0, // 默认无限制
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

// NewAIBaseCrawlerWithLimit 创建有限制的AIBase爬虫
func NewAIBaseCrawlerWithLimit(startID, maxCount int) *AIBaseCrawler {
	crawler := NewAIBaseCrawler(startID)
	crawler.maxCount = maxCount
	return crawler
}

// GetName 获取爬虫名称
func (c *AIBaseCrawler) GetName() string {
	return c.name
}

// GetURL 获取目标URL
func (c *AIBaseCrawler) GetURL() string {
	return c.baseURL
}

// IsEnabled 检查是否启用
func (c *AIBaseCrawler) IsEnabled() bool {
	return c.enabled
}

// ArticlesByDate 按日期组织的文章结构
type ArticlesByDate struct {
	Date     string          `json:"date"`
	Articles []*model.AiNews `json:"articles"`
}

// CrawlResult 爬取结果
type CrawlResult struct {
	ArticlesByDate []*ArticlesByDate `json:"articles_by_date"`
	TotalCount     int               `json:"total_count"`
	LastID         int               `json:"last_id"`
}

// Crawl 执行爬取（实现CrawlerInterface接口）
func (c *AIBaseCrawler) Crawl() ([]*model.AiNews, error) {
	result, err := c.CrawlByDate()
	if err != nil {
		return nil, err
	}

	// 将按日期组织的文章展平为单一列表
	var allNews []*model.AiNews
	for _, dateGroup := range result.ArticlesByDate {
		allNews = append(allNews, dateGroup.Articles...)
	}

	return allNews, nil
}

// CrawlByDate 按日期爬取文章
func (c *AIBaseCrawler) CrawlByDate() (*CrawlResult, error) {
	if c.maxCount > 0 {
		fmt.Printf("（限制 %d 篇）", c.maxCount)
	}
	fmt.Println("...")

	articlesByDate := make(map[string][]*model.AiNews)
	currentID := c.startID
	totalCount := 0
	consecutiveErrors := 0
	maxConsecutiveErrors := 5 // 连续错误超过5次就停止

	for {
		// 检查是否达到最大数量限制
		if c.maxCount > 0 && totalCount >= c.maxCount {
			break
		}

		url := fmt.Sprintf("%s%d", c.baseURL, currentID)

		article, err := c.fetchArticle(url, currentID)
		if err != nil {
			consecutiveErrors++
			if consecutiveErrors >= maxConsecutiveErrors {
				break
			}

			currentID--
			continue
		}

		// 重置连续错误计数
		consecutiveErrors = 0

		if article != nil {
			// 按日期分组
			dateKey := time.Time(article.PublishDate).Format("2006-01-02")
			articlesByDate[dateKey] = append(articlesByDate[dateKey], article)
			totalCount++
		}

		currentID--

		// 添加延迟避免请求过快
		time.Sleep(1 * time.Second)
	}

	// 转换为有序的结果
	var result []*ArticlesByDate
	for date, articles := range articlesByDate {
		result = append(result, &ArticlesByDate{
			Date:     date,
			Articles: articles,
		})
	}

	// 按日期排序（最新的在前）
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Date < result[j].Date {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return &CrawlResult{
		ArticlesByDate: result,
		TotalCount:     totalCount,
		LastID:         currentID + 1, // 最后成功的ID
	}, nil
}

// fetchArticle 获取单篇文章
func (c *AIBaseCrawler) fetchArticle(url string, id int) (*model.AiNews, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查404
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("文章不存在 (404)")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	content := string(body)
	return c.parseArticle(content, url, id)
}

// parseArticle 解析文章内容
func (c *AIBaseCrawler) parseArticle(content, url string, id int) (*model.AiNews, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return nil, err
	}

	title := c.extractTitle(doc)
	publishDate := c.extractPublishDate(doc)
	articleContent := c.extractContent(doc)
	summary := c.extractSummary(articleContent)

	// 生成内容hash用于去重
	hashData := fmt.Sprintf("%s|%s|%s", title, articleContent, url)
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(hashData)))

	return &model.AiNews{
		ID:          id, // 设置文章ID
		Title:       title,
		Content:     articleContent,
		Summary:     summary,
		SourceURL:   url,
		SourceName:  "AIBase",
		PublishDate: localTime.LocalTime(publishDate),
		Category:    "AI资讯", // 统一分类
		Tags:        "",     // 不生成标签
		Status:      0,      // 默认隐藏，需要人工审核
		Hash:        hash,   // 内容hash用于去重
	}, nil
}

// extractContent 提取文章内容（保留HTML样式）
func (c *AIBaseCrawler) extractContent(doc *goquery.Document) string {
	// 使用更精确的选择器提取文章内容
	content := ""

	// 首先尝试提取主要内容区域
	doc.Find(".post-content").Each(func(i int, s *goquery.Selection) {
		// 移除不需要的元素但保留文本格式
		s.Find("script, style, nav, header, footer, .advertisement, .ads, .nav, .menu, .sidebar").Remove()

		// 提取HTML内容而不是纯文本
		html, err := s.Html()
		if err == nil && strings.TrimSpace(html) != "" {
			content += html + "\n"
		}
	})

	// 如果没有找到 .post-content，尝试其他选择器
	if content == "" {
		// 尝试查找文章主体内容
		doc.Find("article .leading-8").Each(func(i int, s *goquery.Selection) {
			s.Find("script, style, nav, header, footer, .advertisement, .ads, .nav, .menu, .sidebar").Remove()

			html, err := s.Html()
			if err == nil && strings.TrimSpace(html) != "" {
				content += html + "\n"
			}
		})
	}

	// 如果还是没有找到，使用通用选择器
	if content == "" {
		doc.Find("main, article, .content, .article-content").Each(func(i int, s *goquery.Selection) {
			// 移除导航、页眉、页脚等无关内容，但保留文本格式标签
			s.Find("nav, header, footer, .nav, .menu, .sidebar, .advertisement, .ads, script, style").Remove()

			// 提取段落内容，保留HTML格式
			s.Find("p, div.text, .article-text, h1, h2, h3, h4, h5, h6, ul, ol, li, blockquote, pre, code").Each(func(j int, p *goquery.Selection) {
				html, err := p.Html()
				if err == nil {
					text := strings.TrimSpace(p.Text())
					if len(text) > 20 { // 过滤掉太短的文本
						// 重新构建带标签的HTML
						tagName := goquery.NodeName(p)
						content += fmt.Sprintf("<%s>%s</%s>\n\n", tagName, html, tagName)
					}
				}
			})
		})
	}

	// 清理内容，但保留HTML结构
	content = strings.TrimSpace(content)
	content = regexp.MustCompile(`\n{3,}`).ReplaceAllString(content, "\n\n")

	// 清理一些可能有害的标签，但保留格式标签
	content = c.sanitizeHTML(content)

	return content
}

// sanitizeHTML 清理HTML内容，移除有害标签但保留格式标签
func (c *AIBaseCrawler) sanitizeHTML(html string) string {
	// 移除有害的标签和属性
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html
	}

	// 移除有害标签
	doc.Find("script, style, iframe, object, embed, form, input, button").Remove()

	// 移除有害属性
	doc.Find("*").Each(func(i int, s *goquery.Selection) {
		// 移除事件处理器和危险属性
		dangerousAttrs := []string{
			"onclick", "onload", "onerror", "onmouseover", "onmouseout",
			"onfocus", "onblur", "onchange", "onsubmit", "onreset",
			"javascript:", "vbscript:", "data:", "javascript",
		}

		for _, attr := range dangerousAttrs {
			s.RemoveAttr(attr)
		}

		// 清理href和src属性中的危险内容
		if href, exists := s.Attr("href"); exists {
			if strings.Contains(strings.ToLower(href), "javascript:") ||
				strings.Contains(strings.ToLower(href), "vbscript:") ||
				strings.Contains(strings.ToLower(href), "data:") {
				s.RemoveAttr("href")
			}
		}

		if src, exists := s.Attr("src"); exists {
			if strings.Contains(strings.ToLower(src), "javascript:") ||
				strings.Contains(strings.ToLower(src), "vbscript:") {
				s.RemoveAttr("src")
			}
		}
	})

	// 返回清理后的HTML
	result, err := doc.Find("body").Html()
	if err != nil {
		return html
	}

	return result
}

// extractTitle 提取文章标题
func (c *AIBaseCrawler) extractTitle(doc *goquery.Document) string {
	// 按优先级尝试不同的选择器
	selectors := []string{
		"h1.font-extrabold", // AIBase特定的标题样式
		"h1.md\\:text-4xl",  // AIBase的响应式标题
		"h1",                // 通用h1标签
		"title",             // 页面标题
		".article-title",    // 通用文章标题类
		".post-title",       // 博客标题类
	}

	for _, selector := range selectors {
		if title := strings.TrimSpace(doc.Find(selector).First().Text()); title != "" {
			return title
		}
	}

	return "未知标题"
}

// extractPublishDate 提取发布日期
func (c *AIBaseCrawler) extractPublishDate(doc *goquery.Document) time.Time {
	// 尝试多种日期格式和选择器
	dateSelectors := []string{
		"time[datetime]",                    // 标准时间标签
		".text-surface-500 span:last-child", // AIBase的日期显示位置
		".publish-date",                     // 通用发布日期类
		".date",                             // 通用日期类
		"[data-date]",                       // 数据属性
	}

	for _, selector := range dateSelectors {
		var foundDate time.Time
		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			// 尝试从datetime属性获取
			if datetime, exists := s.Attr("datetime"); exists {
				if date, err := time.Parse(time.RFC3339, datetime); err == nil {
					foundDate = date
					return
				}
			}

			// 尝试从文本内容获取
			dateText := strings.TrimSpace(s.Text())
			if dateText != "" {
				// 尝试多种日期格式
				formats := []string{
					"January 2, 2006",
					"Jan 2, 2006",
					"2006-01-02",
					"2006/01/02",
					"02/01/2006",
					"2006年01月02日",
					"2006年1月2日",
				}

				for _, format := range formats {
					if date, err := time.Parse(format, dateText); err == nil {
						foundDate = date
						return
					}
				}
			}
		})

		if !foundDate.IsZero() {
			return foundDate
		}
	}

	// 如果上面的选择器都没找到，尝试在整个文档中搜索日期模式
	docText := doc.Text()

	// 尝试匹配常见的日期模式
	datePatterns := []struct {
		pattern string
		format  string
	}{
		{`(May\s+\d{1,2},\s+\d{4})`, "January 2, 2006"},
		{`(Jan\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Feb\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Mar\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Apr\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Jun\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Jul\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Aug\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Sep\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Oct\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Nov\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(Dec\s+\d{1,2},\s+\d{4})`, "Jan 2, 2006"},
		{`(\d{4}年\d{1,2}月\d{1,2}日)`, "2006年1月2日"},
		{`(\d{4}-\d{1,2}-\d{1,2})`, "2006-01-02"},
		{`(\d{4}/\d{1,2}/\d{1,2})`, "2006/01/02"},
	}

	for _, dp := range datePatterns {
		re := regexp.MustCompile(dp.pattern)
		if match := re.FindStringSubmatch(docText); len(match) > 1 {
			dateStr := strings.TrimSpace(match[1])
			if date, err := time.Parse(dp.format, dateStr); err == nil {
				return date
			}
		}
	}

	// 如果无法解析日期，返回当前时间
	return time.Now()
}

// extractSummary 提取摘要（从HTML内容中提取纯文本）
func (c *AIBaseCrawler) extractSummary(content string) string {
	// 如果内容包含HTML标签，先提取纯文本
	if strings.Contains(content, "<") && strings.Contains(content, ">") {
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
		if err == nil {
			content = strings.TrimSpace(doc.Text())
		}
	}

	// 取前200个字符作为摘要
	runes := []rune(content)
	if len(runes) <= 200 {
		return content
	}
	return string(runes[:200]) + "..."
}

// FetchSingleArticle 公共方法：爬取单篇文章
func (c *AIBaseCrawler) FetchSingleArticle(url string, id int) (*model.AiNews, error) {
	return c.fetchArticle(url, id)
}
