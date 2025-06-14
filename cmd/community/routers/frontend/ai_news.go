package frontend

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var (
	aiNewsService      = new(services.AiNewsService)
	aiNewsShareService = new(services.ShareService)
)

// InitAiNewsRouter 初始化AI新闻路由
func InitAiNewsRouter(r *gin.Engine) {
	group := r.Group("/community/ai-news")
	{
		group.Use(middleware.OperLogger())
		group.GET("/dates", getHistoryDates)    // 获取有AI日报的历史日期
		group.GET("/daily", getDailyNews)       // 获取指定日期的AI日报
		group.GET("/detail/:id", getNewsDetail) // 获取AI日报详情
	}
}

// HistoryDate 历史日期响应结构
type HistoryDate struct {
	Date      string `json:"date"`      // 日期 YYYY-MM-DD
	DateLabel string `json:"dateLabel"` // 显示标签，如"1月15日"
	Count     int    `json:"count"`     // 当天文章数量
}

// ShareStats 分享统计结构
type ShareStats struct {
	ShareCount int `json:"shareCount"` // 分享次数
	ShareViews int `json:"shareViews"` // 分享浏览量
}

// getHistoryDates godoc
// @Summary 获取有AI日报的历史日期列表
// @Tags AI News
// @Produce json
// @Success 200 {array} HistoryDate
// @Router /community/ai-news/dates [get]
func getHistoryDates(c *gin.Context) {
	dates, err := aiNewsService.GetHistoryDates()
	if err != nil {
		log.Warnf("获取AI新闻历史日期失败: %s", err.Error())
		result.Err("获取历史日期失败").Json(c)
		return
	}

	// 转换为前端需要的格式
	var historyDates []HistoryDate
	for _, date := range dates {
		parsedDate, err := time.Parse("2006-01-02", date.Date)
		if err != nil {
			continue
		}

		// 生成显示标签，如"1月15日"
		dateLabel := parsedDate.Format("1月2日")

		historyDates = append(historyDates, HistoryDate{
			Date:      date.Date,
			DateLabel: dateLabel,
			Count:     date.Count,
		})
	}

	result.Ok(historyDates, "").Json(c)
}

// DailyNewsItem 日报文章项 - 移除了SourceURL和Category字段
type DailyNewsItem struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Summary      string     `json:"summary"`
	Tags         string     `json:"tags"`
	PublishDate  string     `json:"publishDate"`
	Content      string     `json:"content,omitempty"` // 可选，用于详情页
	CommentCount int64      `json:"commentCount"`
	ShareStats   ShareStats `json:"shareStats"` // 分享统计
}

// DailyNewsResponse 日报响应结构
type DailyNewsResponse struct {
	Date     string          `json:"date"`
	Articles []DailyNewsItem `json:"articles"`
	Total    int             `json:"total"`
}

// getDailyNews godoc
// @Summary 获取指定日期的AI日报
// @Tags AI News
// @Produce json
// @Param date query string true "日期 YYYY-MM-DD"
// @Param with_content query bool false "是否包含文章内容"
// @Success 200 {object} DailyNewsResponse
// @Router /community/ai-news/daily [get]
func getDailyNews(c *gin.Context) {
	date := c.Query("date")
	if date == "" {
		// 如果没有指定日期，获取最新日期的新闻
		latestDate, err := aiNewsService.GetLatestDate()
		if err != nil {
			log.Warnf("获取最新日期失败: %s", err.Error())
			result.Err("获取最新日期失败").Json(c)
			return
		}
		date = latestDate
	}

	// 验证日期格式
	_, err := time.Parse("2006-01-02", date)
	if err != nil {
		result.Err("日期格式错误，请使用 YYYY-MM-DD 格式").Json(c)
		return
	}

	// 检查是否包含文章内容
	withContent := c.Query("with_content") == "true"

	articles, err := aiNewsService.GetDailyNews(date, withContent)
	if err != nil {
		log.Warnf("获取日期 %s 的AI新闻失败: %s", date, err.Error())
		result.Err("获取AI日报失败").Json(c)
		return
	}

	// 批量查询所有文章的评论数量，避免N+1查询问题
	var commentCounts []struct {
		BusinessId int   `json:"business_id"`
		Count      int64 `json:"count"`
	}
	if len(articles) > 0 {
		articleIds := make([]int, len(articles))
		for i, article := range articles {
			articleIds[i] = article.ID
		}

		model.Comment().
			Select("business_id, COUNT(*) as count").
			Where("business_id IN ? AND tenant_id = ? AND deleted_at IS NULL", articleIds, 4).
			Group("business_id").
			Scan(&commentCounts)
	}

	// 建立评论数量映射
	commentCountMap := make(map[int]int64)
	for _, cc := range commentCounts {
		commentCountMap[cc.BusinessId] = cc.Count
	}

	// 批量查询分享统计
	var shareStatsMap map[int]services.ShareStatistics
	if len(articles) > 0 {
		articleIds := make([]int, len(articles))
		for i, article := range articles {
			articleIds[i] = article.ID
		}

		shareStatsMap, err = aiNewsShareService.GetStatisticsByBusinessIDs(services.BusinessTypeAiNews, articleIds)
		if err != nil {
			log.Warnf("获取分享统计失败: %s", err.Error())
			// 分享统计失败不影响主要功能，继续执行
			shareStatsMap = make(map[int]services.ShareStatistics)
		}
	}

	// 转换为前端需要的格式 - 移除了SourceURL和Category字段
	var dailyItems []DailyNewsItem
	for _, article := range articles {
		// 获取分享统计
		shareStats := shareStatsMap[article.ID]

		item := DailyNewsItem{
			ID:           article.ID,
			Title:        article.Title,
			Summary:      article.Summary,
			Tags:         article.Tags,
			PublishDate:  article.PublishDate.String(),
			CommentCount: commentCountMap[article.ID], // 从映射中获取评论数量
			ShareStats: ShareStats{
				ShareCount: shareStats.ShareCount,
				ShareViews: shareStats.TotalViews,
			},
		}

		// 如果需要内容，则包含
		if withContent {
			item.Content = article.Content
		}

		dailyItems = append(dailyItems, item)
	}

	response := DailyNewsResponse{
		Date:     date,
		Articles: dailyItems,
		Total:    len(dailyItems),
	}

	result.Ok(response, "").Json(c)
}

// NewsDetailResponse AI新闻详情响应结构
type NewsDetailResponse struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	Summary      string     `json:"summary"`
	Category     string     `json:"category"`
	Tags         string     `json:"tags"`
	PublishDate  string     `json:"publishDate"`
	CreatedAt    string     `json:"createdAt"`
	CommentCount int64      `json:"commentCount"` // 评论数量
	ShareStats   ShareStats `json:"shareStats"`   // 分享统计
}

// getNewsDetail godoc
// @Summary 获取AI新闻详情
// @Tags AI News
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} NewsDetailResponse
// @Router /community/ai-news/detail/{id} [get]
func getNewsDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		result.Err("无效的文章ID").Json(c)
		return
	}

	article, err := aiNewsService.GetNewsById(id)
	if err != nil {
		log.Warnf("获取AI新闻详情失败: ID=%d, err=%s", id, err.Error())
		result.Err("文章不存在或已删除").Json(c)
		return
	}

	// 查询评论数量
	var commentCount int64
	model.Comment().Where("business_id = ? AND tenant_id = ? AND deleted_at IS NULL", article.ID, 4).Count(&commentCount)

	// 查询分享统计
	shareStatsMap, err := aiNewsShareService.GetStatisticsByBusinessIDs(services.BusinessTypeAiNews, []int{article.ID})
	if err != nil {
		log.Warnf("获取分享统计失败: %s", err.Error())
		shareStatsMap = make(map[int]services.ShareStatistics)
	}
	shareStats := shareStatsMap[article.ID]

	response := NewsDetailResponse{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		Summary:      article.Summary,
		Category:     article.Category,
		Tags:         article.Tags,
		PublishDate:  article.PublishDate.String(),
		CreatedAt:    article.CreatedAt.String(),
		CommentCount: commentCount,
		ShareStats: ShareStats{
			ShareCount: shareStats.ShareCount,
			ShareViews: shareStats.TotalViews,
		},
	}

	result.Ok(response, "").Json(c)
}
