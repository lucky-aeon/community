package frontend

import (
	"strconv"
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/service"
)

// ReactionRequest 表情回复请求
type ReactionRequest struct {
	CommentId    int    `json:"commentId" binding:"required"`
	ReactionType string `json:"reactionType" binding:"required"`
}

// InitCommentReactionRouters 初始化评论表情回复路由
func InitCommentReactionRouters(g *gin.Engine) {
	group := g.Group("/community/comments/reactions")
	
	// 获取表情类型列表
	group.GET("/types", getExpressionTypes)
	
	// 获取评论表情统计
	group.GET("/:commentId", getCommentReactionSummary)
	
	// 需要登录的接口
	group.Use(middleware.OperLogger())
	
	// 切换表情回复状态
	group.POST("/toggle", toggleReaction)
	
	// 添加表情回复
	group.POST("/add", addReaction)
	
	// 移除表情回复
	group.DELETE("/remove", removeReaction)
}

// getExpressionTypes 获取所有表情类型
func getExpressionTypes(ctx *gin.Context) {
	reactionService := services.NewCommentReactionService(ctx)
	types, err := reactionService.GetAllExpressionTypes()
	if err != nil {
		log.Errorf("获取表情类型失败: %v", err)
		result.Err("获取表情类型失败").Json(ctx)
		return
	}
	
	result.Ok(types, "获取成功").Json(ctx)
}

// getCommentReactionSummary 获取评论表情统计
func getCommentReactionSummary(ctx *gin.Context) {
	commentIdStr := ctx.Param("commentId")
	commentId, err := strconv.Atoi(commentIdStr)
	if err != nil {
		log.Warnf("评论ID参数错误: %v", err)
		result.Err("评论ID参数错误").Json(ctx)
		return
	}
	
	currentUserId := middleware.GetUserId(ctx)
	if currentUserId == 0 {
		// 未登录用户，不显示用户是否已回复的信息
		currentUserId = -1
	}
	
	reactionService := services.NewCommentReactionService(ctx)
	summaries, err := reactionService.GetReactionSummary(commentId, currentUserId)
	if err != nil {
		log.Errorf("获取评论表情统计失败: %v", err)
		result.Err("获取统计失败").Json(ctx)
		return
	}
	
	result.Ok(summaries, "获取成功").Json(ctx)
}

// toggleReaction 切换表情回复状态
func toggleReaction(ctx *gin.Context) {
	var req ReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warnf("表情回复请求参数错误: %v", err)
		result.Err("请求参数错误").Json(ctx)
		return
	}
	
	userId := middleware.GetUserId(ctx)
	if userId == 0 {
		result.Err("请先登录").Json(ctx)
		return
	}
	
	reactionService := services.NewCommentReactionService(ctx)
	added, err := reactionService.ToggleReaction(req.CommentId, userId, req.ReactionType)
	if err != nil {
		log.Errorf("切换表情回复失败: %v", err)
		result.Err(err.Error()).Json(ctx)
		return
	}
	
	message := "表情回复已移除"
	if added {
		message = "表情回复已添加"
	}
	
	result.Ok(map[string]interface{}{
		"added": added,
		"commentId": req.CommentId,
		"reactionType": req.ReactionType,
	}, message).Json(ctx)
}

// addReaction 添加表情回复
func addReaction(ctx *gin.Context) {
	var req ReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warnf("添加表情回复请求参数错误: %v", err)
		result.Err("请求参数错误").Json(ctx)
		return
	}
	
	userId := middleware.GetUserId(ctx)
	if userId == 0 {
		result.Err("请先登录").Json(ctx)
		return
	}
	
	reactionService := services.NewCommentReactionService(ctx)
	err := reactionService.AddReaction(req.CommentId, userId, req.ReactionType)
	if err != nil {
		log.Errorf("添加表情回复失败: %v", err)
		result.Err(err.Error()).Json(ctx)
		return
	}
	
	result.Ok(nil, "表情回复添加成功").Json(ctx)
}

// removeReaction 移除表情回复
func removeReaction(ctx *gin.Context) {
	var req ReactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warnf("移除表情回复请求参数错误: %v", err)
		result.Err("请求参数错误").Json(ctx)
		return
	}
	
	userId := middleware.GetUserId(ctx)
	if userId == 0 {
		result.Err("请先登录").Json(ctx)
		return
	}
	
	reactionService := services.NewCommentReactionService(ctx)
	err := reactionService.RemoveReaction(req.CommentId, userId, req.ReactionType)
	if err != nil {
		log.Errorf("移除表情回复失败: %v", err)
		result.Err(err.Error()).Json(ctx)
		return
	}
	
	result.Ok(nil, "表情回复移除成功").Json(ctx)
}