package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	services "xhyovo.cn/community/server/service"
)

func InitCommentRouters(r *gin.Engine) {
	group := r.Group("/community/admin/comment")
	group.GET("", listComment)
	group.DELETE("/:id", deleteComment)
}

func listComment(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	var c services.CommentsService
	comments, count := c.PageComment(p, limit)
	result.Page(comments, count, nil).Json(ctx)
}

func deleteComment(ctx *gin.Context) {

	id, _ := strconv.Atoi(ctx.Param("id"))

	var c services.CommentsService
	if !c.DeleteComment(id, 0) {
		result.Err("删除失败").Json(ctx)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}
