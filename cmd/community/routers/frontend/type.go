package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

func InitTypeRouters(r *gin.Engine) {
	r.Group("")
}

func list(ctx *gin.Context) {
	parentId, err := strconv.Atoi(ctx.DefaultQuery("parentId", "0"))
	if err != nil {
		parentId = 0
	}
	var typeService services.TypeService
	result.Ok(typeService.List(uint(parentId)), "").Json(ctx)
}
