package frontend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

func InitMessageRouters(r *gin.Engine) {
	group := r.Group("/community/message")
	group.GET("/listNoRead", listNoReadMsg)
	group.DELETE("/deleteNoRead", deleteMessage)
}

// 查看用户未读消息
func listNoReadMsg(ctx *gin.Context) {
	var msgService services.MessageService
	message := msgService.ListNoReadMessage(1, 5, 2)
	result.Ok(&message, "").Json(ctx)
}

// 删除用户收到的消息(确认消息),
func deleteMessage(ctx *gin.Context) {
	var ids []int
	if err := ctx.ShouldBindJSON(&ids); err != nil && len(ids) > 0 {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var msgService services.MessageService
	msgService.ReadMessage(ids, middleware.GetUserId(ctx))
	result.Ok(nil, "已读").Json(ctx)
}
