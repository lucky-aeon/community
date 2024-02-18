package frontend

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	services "xhyovo.cn/community/server/service"
)

func InitMessageRouters(r *gin.Engine) {
	group := r.Group("/community/message")
	group.GET("", listMsg)
	group.POST("/read", readMsg)

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
		result.Err(err.Error()).Json(ctx)
		return
	}
	var msgService services.MessageService
	number := msgService.ReadMessage(ids, middleware.GetUserId(ctx))
	result.OkWithMsg(nil, fmt.Sprintf("已读%d消息", number)).Json(ctx)
}
