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
	g.GET("/tree", treeTypes)
}

func list(ctx *gin.Context) {
	parentId, err := strconv.Atoi(ctx.DefaultQuery("parentId", "0"))
	if err != nil {
		parentId = 0
	}
	var typeService services.TypeService
	result.Ok(typeService.List(parentId), "").Json(ctx)
}

func treeTypes(ctx *gin.Context) {
	var typeService services.TypeService
	types, _ := typeService.PageTypes(1, 99)
	result.Ok(types, "").Json(ctx)
}
