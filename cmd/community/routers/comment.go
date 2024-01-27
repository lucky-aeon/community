package routers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

func InitCommentRouters(g *gin.Engine) {
	group := g.Group("/community/comments")
	// todo  这里的路由取名字不会取
	group.GET("/byArticleId", listByArticleId)
}

// 返回文章下的评论
func listByArticleId(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Query("id"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "15"))

	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var commentsService services.CommentsService
	comments := commentsService.GetCommentsByArticleID(uint(page), uint(limit), uint(articleId))
	result.Ok(comments, "").Json(ctx)
}
