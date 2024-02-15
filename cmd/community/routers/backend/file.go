package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	services "xhyovo.cn/community/server/service"
)

func InitFileRouters(r *gin.Engine) {
	group := r.Group("/community/admin/file")
	group.GET("", listFiles)
}

func listFiles(ctx *gin.Context) {

	p, limit := page.GetPage(ctx)
	userId := ctx.DefaultQuery("userId", "0")

	uId, _ := strconv.Atoi(userId)

	var fileS services.FileService

	files, count := fileS.PageFiles(p, limit, uId)
	result.Ok(page.New(files, count), "").Json(ctx)
}
