package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
)

func InitIndexRouters(ctx *gin.Engine) {
	group := ctx.Group("/community")
	// 社区首页
	group.GET("/user/count", getUserCount)
}

func getUserCount(ctx *gin.Context) {
	var count int64
	model.User().Count(&count)
	result.Ok(count, "").Json(ctx)

}
