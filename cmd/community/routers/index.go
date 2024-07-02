package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitIndexRouters(ctx *gin.Engine) {
	group := ctx.Group("/community")
	// 社区首页
	group.GET("/user/count", getUserCount)
	group.GET("/rate/page", pageRate)
}

func pageRate(ctx *gin.Context) {

	p, limit := page.GetPage(ctx)
	var noteService services.RateService
	state, notes := noteService.Page(p, limit)

	result.Ok(map[string]interface{}{
		"state": state,
		"data":  notes,
	}, "").Json(ctx)
	return
}

func getUserCount(ctx *gin.Context) {
	var count int64
	model.User().Count(&count)
	result.Ok(count, "").Json(ctx)

}
