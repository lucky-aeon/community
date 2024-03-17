package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	services "xhyovo.cn/community/server/service"
)

func InitMessageRouters(r *gin.Engine) {
	group := r.Group("/community/message")
	group.GET("/unReader/count", getUnReadMsgCount)
	group.GET("", listMsg)
	group.POST("/read", readMsg, middleware.OperLogger())
	group.DELETE("/UnReadMsg/:type", clearUnReadMsg, middleware.OperLogger())

}

func getUnReadMsgCount(ctx *gin.Context) {
	userId := middleware.GetUserId(ctx)
	var msgService services.MessageService
	count := msgService.GetUnReadMessageCountByUserId(userId)
	result.Ok(count, "").Json(ctx)
}

// 查看用户未读消息
func listMsg(ctx *gin.Context) {
	types := ctx.Query("type")
	states := ctx.Query("state")
	p, limit := page.GetPage(ctx)
	atoi, _ := strconv.Atoi(types)
	atoi2, _ := strconv.Atoi(states)
	userId := middleware.GetUserId(ctx)
	var msgService services.MessageService
	message, count := msgService.PageMessage(p, limit, userId, atoi, atoi2)
	result.Ok(page.New(message, count), "").Json(ctx)
}

// 阅读消息
func readMsg(ctx *gin.Context) {
	var ids []int
	if err := ctx.ShouldBindJSON(&ids); err != nil && len(ids) > 0 {
		log.Warnf("用户id: %d 阅读消息参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	var msgService services.MessageService
	number := msgService.ReadMessage(ids, middleware.GetUserId(ctx))
	result.OkWithMsg(nil, fmt.Sprintf("已读%d消息", number)).Json(ctx)
}

// 清除未读消息
func clearUnReadMsg(ctx *gin.Context) {
	msgType, err := strconv.Atoi(ctx.Param("type"))
	if err != nil {
		log.Warnln("用户id: %d 清除未读消息参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	var msgService services.MessageService
	msgService.ClearUnReadMessage(msgType, middleware.GetUserId(ctx))
	result.OkWithMsg(nil, "已清空").Json(ctx)
}
