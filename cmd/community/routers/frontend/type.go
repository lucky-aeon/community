package frontend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

func InitTypeRouters(r *gin.Engine) {
	g := r.Group("/community/classfiys")
	g.GET("", list)
}

func list(ctx *gin.Context) {
	parentId, err := strconv.Atoi(ctx.DefaultQuery("parentId", "0"))
	if err != nil {
		parentId = 0
	}
	var typeService services.TypeService
	result.Ok(typeService.List(parentId), "").Json(ctx)
}
