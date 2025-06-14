package frontend

import (
	"strings"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var (
	shareService       = new(services.ShareService)
	shareAiNewsService = new(services.AiNewsService)
)

// InitShareRouter 初始化分享路由
func InitShareRouter(r *gin.Engine) {
	// 公开访问的分享路由（不需要登录）
	publicGroup := r.Group("/community/share")
	{
		publicGroup.GET("/:token", accessShare)              // 访问分享内容
		publicGroup.GET("/business-types", getBusinessTypes) // 获取支持分享的业务类型
	}
	// 分享管理路由（需要登录）
	shareGroup := r.Group("/community/share")
	{
		shareGroup.Use(middleware.Auth) // 添加登录验证中间件
		shareGroup.Use(middleware.OperLogger())
		shareGroup.POST("/create", createShare)        // 创建分享
		shareGroup.GET("/:token/stats", getShareStats) // 获取分享统计
	}

}

// CreateShareRequest 创建分享请求
type CreateShareRequest struct {
	BusinessType string `json:"business_type" binding:"required"` // ai_news, article等
	BusinessID   int    `json:"business_id" binding:"required"`
	ExpireDays   int    `json:"expire_days"` // 过期天数，0表示永久
}

// createShare godoc
// @Summary 创建分享
// @Tags Share
// @Accept json
// @Produce json
// @Param request body CreateShareRequest true "创建分享请求"
// @Success 200 {object} services.ShareResponse
// @Router /community/share/create [post]
func createShare(c *gin.Context) {
	var req CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	// 从Auth中间件获取用户ID
	creatorID := middleware.GetUserId(c)

	shareReq := services.CreateShareRequest{
		BusinessType: req.BusinessType,
		BusinessID:   req.BusinessID,
		CreatorID:    creatorID,
		ExpireDays:   req.ExpireDays,
	}

	shareResp, err := shareService.CreateShare(shareReq)
	if err != nil {
		log.Warnf("创建分享失败: %s", err.Error())
		result.Err("创建分享失败").Json(c)
		return
	}

	result.Ok(shareResp, "分享创建成功").Json(c)
}

// accessShare godoc
// @Summary 访问分享内容
// @Tags Share
// @Produce json
// @Param token path string true "分享token"
// @Success 200 {object} interface{}
// @Router /community/share/{token} [get]
func accessShare(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		result.Err("无效的分享链接").Json(c)
		return
	}

	// 获取分享信息
	share, err := shareService.GetShareByToken(token)
	if err != nil {
		log.Warnf("获取分享信息失败: token=%s, err=%s", token, err.Error())
		result.Err("分享链接无效或已过期").Json(c)
		return
	}

	// 记录访问
	visitorInfo := services.VisitorInfo{
		IP:        getClientIP(c),
		UserAgent: c.GetHeader("User-Agent"),
		Referer:   c.GetHeader("Referer"),
	}

	// 如果用户已登录，记录用户ID
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(int); ok {
			visitorInfo.UserID = uid
		}
	}

	// 异步记录访问，不影响响应速度
	go func() {
		if err := shareService.RecordView(share.ID, visitorInfo); err != nil {
			log.Warnf("记录分享访问失败: shareID=%d, err=%s", share.ID, err.Error())
		}
	}()

	// 根据业务类型获取具体内容
	switch share.BusinessType {
	case services.BusinessTypeAiNews:
		content, err := getAiNewsContent(share.BusinessID)
		if err != nil {
			result.Err("获取内容失败").Json(c)
			return
		}
		result.Ok(content, "").Json(c)
	default:
		result.Err("不支持的分享类型").Json(c)
	}
}

// getShareStats godoc
// @Summary 获取分享统计
// @Tags Share
// @Produce json
// @Param token path string true "分享token"
// @Success 200 {object} interface{}
// @Router /community/share/{token}/stats [get]
func getShareStats(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		result.Err("无效的分享token").Json(c)
		return
	}

	// 获取分享信息
	share, err := shareService.GetShareByToken(token)
	if err != nil {
		log.Warnf("获取分享信息失败: token=%s, err=%s", token, err.Error())
		result.Err("分享链接无效或已过期").Json(c)
		return
	}

	// 统计信息
	stats := map[string]interface{}{
		"total_views": share.TotalViews,
		"created_at":  share.CreatedAt,
		"expire_at":   share.ExpireAt,
	}

	result.Ok(stats, "").Json(c)
}

// getAiNewsContent 获取AI新闻内容
func getAiNewsContent(id int) (interface{}, error) {
	article, err := shareAiNewsService.GetNewsById(id)
	if err != nil {
		return nil, err
	}

	// 查询评论数量
	var commentCount int64
	model.Comment().Where("business_id = ? AND tenant_id = ? AND deleted_at IS NULL", article.ID, 4).Count(&commentCount)

	response := map[string]interface{}{
		"id":            article.ID,
		"title":         article.Title,
		"content":       article.Content,
		"summary":       article.Summary,
		"category":      article.Category,
		"tags":          article.Tags,
		"publish_date":  article.PublishDate.String(),
		"created_at":    article.CreatedAt.String(),
		"comment_count": commentCount,
		"share_type":    services.BusinessTypeAiNews,
	}

	return response, nil
}

// getClientIP 获取客户端真实IP
func getClientIP(c *gin.Context) string {
	clientIP := c.ClientIP()

	// 检查代理头
	if forwarded := c.GetHeader("X-Forwarded-For"); forwarded != "" {
		// X-Forwarded-For可能包含多个IP，取第一个
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}

	return clientIP
}

// getBusinessTypes godoc
// @Summary 获取支持分享的业务类型列表
// @Tags Share
// @Produce json
// @Success 200 {object} interface{}
// @Router /community/share/business-types [get]
func getBusinessTypes(c *gin.Context) {
	businessTypes := services.GetAllowedBusinessTypes()

	// 构造响应数据，包含业务类型和描述
	typeMap := map[string]string{
		services.BusinessTypeAiNews: "AI日报",
		// services.BusinessTypeArticle: "文章",  // 未来可能开放
		// services.BusinessTypePost:    "帖子",  // 未来可能开放
	}

	var response []map[string]interface{}
	for _, businessType := range businessTypes {
		response = append(response, map[string]interface{}{
			"type":        businessType,
			"description": typeMap[businessType],
		})
	}

	result.Ok(response, "").Json(c)
}
