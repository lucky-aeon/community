package frontend

import (
	"strconv"

	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"

	"github.com/gin-gonic/gin"
)

var (
	chatService = new(services.ChatService)
)

func InitChatRouter(r *gin.Engine) {
	group := r.Group("/community/chat")
	group.Use(middleware.OperLogger())
	{
		group.GET("/models", getAIModels)
		group.POST("/groups", createChatGroup)
		group.GET("/groups", getUserChatGroups)
		//group.GET("/groups/:id", getChatGroup)
		group.PUT("/groups/:id", updateChatGroup)
		group.DELETE("/groups/:id", deleteChatGroup)
		group.POST("/groups/:id/messages", sendMessage)
		group.GET("/groups/:id/messages", getMessages)
	}
}

// getAIModels godoc
// @Summary 获取所有可用的AI模型
// @Tags Chat
// @Produce json
// @Success 200 {array} model.AIModels
// @Router /community/chat/models [get]
func getAIModels(c *gin.Context) {
	models, err := chatService.GetAIModels()
	if err != nil {
		log.Warnf("用户id: %d 获取AI模型列表失败: %s", middleware.GetUserId(c), err.Error())
		result.Err(err.Error()).Json(c)
		return
	}
	result.Ok(models, "").Json(c)
}

// createChatGroup godoc
// @Summary 创建新的对话分组
// @Tags Chat
// @Accept json
// @Produce json
// @Param title body string true "对话分组标题"
// @Success 200 {object} model.ChatGroups
// @Router /community/chat/groups [post]
func createChatGroup(c *gin.Context) {
	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("用户id: %d 创建对话分组参数解析失败: %s", middleware.GetUserId(c), err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	userID := middleware.GetUserId(c)
	group, err := chatService.CreateChatGroup(int64(userID), req.Title)
	if err != nil {
		log.Warnf("用户id: %d 创建对话分组失败: %s", userID, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	result.OkWithMsg(group, "创建成功").Json(c)
}

// getUserChatGroups godoc
// @Summary 获取当前用户的所有对话分组
// @Tags Chat
// @Produce json
// @Param page query int false "页码(默认1)"
// @Param limit query int false "每页数量(默认10)"
// @Success 200 {array} model.ChatGroups
// @Router /community/chat/groups [get]
func getUserChatGroups(c *gin.Context) {
	userID := middleware.GetUserId(c)
	p, limit := page.GetPage(c)

	groups, total, err := chatService.GetUserChatGroups(int64(userID), p, limit)
	if err != nil {
		log.Warnf("用户id: %d 获取对话分组列表失败: %s", userID, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	result.Page(groups, total, nil).Json(c)
}

// getChatGroup godoc
// @Summary 获取对话分组详情
// @Tags Chat
// @Produce json
// @Param id path int true "对话分组ID"
// @Success 200 {object} model.ChatGroups
// @Router /community/chat/groups/{id} [get]
func getChatGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warnf("用户id: %d 获取对话分组详情参数解析失败: %s", middleware.GetUserId(c), err.Error())
		result.Err("无效的ID").Json(c)
		return
	}

	group, err := chatService.GetChatGroup(id)
	if err != nil {
		log.Warnf("用户id: %d 获取对话分组详情失败,分组id: %d, err: %s", middleware.GetUserId(c), id, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	result.Ok(group, "").Json(c)
}

// updateChatGroup godoc
// @Summary 更新对话分组
// @Tags Chat
// @Accept json
// @Produce json
// @Param id path int true "对话分组ID"
// @Param title body string true "新的对话分组标题"
// @Success 200 {object} model.ChatGroups
// @Router /community/chat/groups/{id} [put]
func updateChatGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warnf("用户id: %d 更新对话分组参数解析失败: %s", middleware.GetUserId(c), err.Error())
		result.Err("无效的ID").Json(c)
		return
	}

	var req struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("用户id: %d 更新对话分组标题参数解析失败: %s", middleware.GetUserId(c), err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	userID := middleware.GetUserId(c)
	err = chatService.UpdateChatGroup(id, int64(userID), req.Title)
	if err != nil {
		log.Warnf("用户id: %d 更新对话分组失败,分组id: %d, err: %s", userID, id, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	result.OkWithMsg(nil, "更新成功").Json(c)
}

// deleteChatGroup godoc
// @Summary 删除对话分组
// @Tags Chat
// @Produce json
// @Param id path int true "对话分组ID"
// @Success 200 {object} string "success"
// @Router /community/chat/groups/{id} [delete]
func deleteChatGroup(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warnf("用户id: %d 删除对话分组参数解析失败: %s", middleware.GetUserId(c), err.Error())
		result.Err("无效的ID").Json(c)
		return
	}

	userID := middleware.GetUserId(c)
	err = chatService.DeleteChatGroup(id, int64(userID))
	if err != nil {
		log.Warnf("用户id: %d 删除对话分组失败,分组id: %d, err: %s", userID, id, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	result.OkWithMsg(nil, "删除成功").Json(c)
}

// sendMessage godoc
// @Summary 发送消息到对话分组
// @Tags Chat
// @Accept json
// @Produce json
// @Param id path int true "对话分组ID"
// @Param request body model.SendMessageRequest true "发送消息请求"
// @Success 200 {object} model.ChatCompletionChunk
// @Router /community/chat/groups/{id}/messages [post]
func sendMessage(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warnf("用户id: %d 发送消息参数解析失败: %s", middleware.GetUserId(c), err.Error())
		result.Err("无效的分组ID").Json(c)
		return
	}

	var req model.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warnf("用户id: %d 发送消息请求参数解析失败: %s", middleware.GetUserId(c), err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	userID := middleware.GetUserId(c)
	userMessage, aiMessage, err := chatService.SendMessage(groupID, strconv.FormatInt(int64(userID), 10), &req)
	if err != nil {
		log.Warnf("用户id: %d 发送消息失败,分组id: %d, err: %s", userID, groupID, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	result.OkWithMsg(gin.H{
		"userMessage": userMessage,
		"aiMessage":   aiMessage,
	}, "发送成功").Json(c)
}

// getMessages godoc
// @Summary 获取对话分组的消息列表
// @Tags Chat
// @Produce json
// @Param id path int true "对话分组ID"
// @Param page query int false "页码(默认1)"
// @Param page_size query int false "每页数量(默认20)"
// @Success 200 {object} model.ChatMessages
// @Router /community/chat/groups/{id}/messages [get]
func getMessages(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Warnf("用户id: %d 获取消息列表参数解析失败: %s", middleware.GetUserId(c), err.Error())
		result.Err("无效的分组ID").Json(c)
		return
	}

	p, pageSize := page.GetPage(c)

	messages, total, err := chatService.GetMessages(groupID, p, pageSize)
	if err != nil {
		log.Warnf("用户id: %d 获取消息列表失败,分组id: %d, err: %s", middleware.GetUserId(c), groupID, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	result.Page(messages, total, nil).Json(c)
}
