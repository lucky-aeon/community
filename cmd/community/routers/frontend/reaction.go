package frontend

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

func InitReactionRouters(g *gin.Engine) {
	group := g.Group("/community/reactions")
	group.Use(middleware.OperLogger())

	// 切换表情回复状态
	group.POST("/toggle", toggleUniversalReaction)

	// 获取单个业务的表情统计
	group.GET("/:businessType/:businessId", getReactionSummary)

	// 批量获取多个业务的表情统计
	group.GET("/:businessType/batch", getReactionSummaryBatch)
}

// toggleUniversalReaction 切换通用表情回复状态
func toggleUniversalReaction(ctx *gin.Context) {
	userId := middleware.GetUserId(ctx)

	var req struct {
		BusinessType int    `json:"businessType" binding:"min=0"`
		BusinessId   int    `json:"businessId" binding:"required"`
		ReactionType string `json:"reactionType" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Warnf("用户ID: %d 切换表情回复参数解析失败: %v", userId, err)
		result.Err("参数错误").Json(ctx)
		return
	}

	reactionService := services.NewReactionService(ctx)

	// 验证业务类型
	if !reactionService.ValidateBusinessType(req.BusinessType) {
		log.Warnf("用户ID: %d 提供的业务类型无效: %d", userId, req.BusinessType)
		result.Err("无效的业务类型").Json(ctx)
		return
	}

	// 验证表情类型
	if !reactionService.ValidateReactionType(req.ReactionType) {
		log.Warnf("用户ID: %d 提供的表情类型无效: %s", userId, req.ReactionType)
		result.Err("无效的表情类型").Json(ctx)
		return
	}

	added, err := reactionService.ToggleReaction(req.BusinessType, req.BusinessId, userId, req.ReactionType)
	if err != nil {
		log.Errorf("用户ID: %d 切换表情回复失败: %v", userId, err)
		result.Err("操作失败").Json(ctx)
		return
	}

	message := "表情回复已移除"
	if added {
		message = "表情回复已添加"
	}

	result.Ok(map[string]interface{}{
		"added": added,
	}, message).Json(ctx)
}

// getReactionSummary 获取单个业务的表情统计
func getReactionSummary(ctx *gin.Context) {
	businessType, err := strconv.Atoi(ctx.Param("businessType"))
	if err != nil {
		log.Warnf("解析业务类型失败: %v", err)
		result.Err("业务类型参数错误").Json(ctx)
		return
	}

	businessId, err := strconv.Atoi(ctx.Param("businessId"))
	if err != nil {
		log.Warnf("解析业务ID失败: %v", err)
		result.Err("业务ID参数错误").Json(ctx)
		return
	}

	currentUserId := middleware.GetUserId(ctx)

	reactionService := services.NewReactionService(ctx)

	// 验证业务类型
	if !reactionService.ValidateBusinessType(businessType) {
		log.Warnf("提供的业务类型无效: %d", businessType)
		result.Err("无效的业务类型").Json(ctx)
		return
	}

	summaries, err := reactionService.GetReactionSummary(businessType, businessId, currentUserId)
	if err != nil {
		log.Errorf("获取表情统计失败: %v", err)
		result.Err("获取表情统计失败").Json(ctx)
		return
	}

	result.Ok(summaries, "获取成功").Json(ctx)
}

// getReactionSummaryBatch 批量获取多个业务的表情统计
func getReactionSummaryBatch(ctx *gin.Context) {
	businessType, err := strconv.Atoi(ctx.Param("businessType"))
	if err != nil {
		log.Warnf("解析业务类型失败: %v", err)
		result.Err("业务类型参数错误").Json(ctx)
		return
	}

	businessIdsStr := ctx.Query("businessIds")
	if businessIdsStr == "" {
		result.Err("业务ID列表不能为空").Json(ctx)
		return
	}

	// 解析业务ID列表
	businessIdStrs := strings.Split(businessIdsStr, ",")
	businessIds := make([]int, 0, len(businessIdStrs))

	for _, idStr := range businessIdStrs {
		id, err := strconv.Atoi(strings.TrimSpace(idStr))
		if err != nil {
			log.Warnf("解析业务ID失败: %v", err)
			result.Err("业务ID格式错误").Json(ctx)
			return
		}
		businessIds = append(businessIds, id)
	}

	currentUserId := middleware.GetUserId(ctx)

	reactionService := services.NewReactionService(ctx)

	// 验证业务类型
	if !reactionService.ValidateBusinessType(businessType) {
		log.Warnf("提供的业务类型无效: %d", businessType)
		result.Err("无效的业务类型").Json(ctx)
		return
	}

	summaryMap, err := reactionService.GetReactionSummaryBatch(businessType, businessIds, currentUserId)
	if err != nil {
		log.Errorf("批量获取表情统计失败: %v", err)
		result.Err("批量获取表情统计失败").Json(ctx)
		return
	}

	result.Ok(summaryMap, "获取成功").Json(ctx)
}
