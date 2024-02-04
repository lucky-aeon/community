package routers

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitCommentRouters(g *gin.Engine) {
	group := g.Group("/community/comments")

	group.POST("/comment", comment)
	group.DELETE("/:id", deleteComment)
	group.GET("/byArticleId/:articleId", listCommentsByArticleId)
	group.GET("/byRootId/:rootId", listCommentsByRootId)
	group.GET("/allCommentsByArticleId/:articleId", listAllCommentsByArticleId)
}

// 发布评论
func comment(ctx *gin.Context) {
	var comment model.Comments
	var userId int
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		result.Err(utils.GetValidateErr(comment, err)).Json(ctx)
		return
	}
	comment.UserId = userId

	commentsService := services.NewCommentService(ctx)
	err := commentsService.Comment(&comment)
	msg := "评论成功"
	if err != nil {
		msg = err.Error()
	}
	result.Ok(nil, msg).Json(ctx)
}

// 删除评论
func deleteComment(ctx *gin.Context) {
	commentId := ctx.Param("id")

	userId := middleware.GetUserId(ctx)
	if commentId == "" {
		result.Err("删除评论id不能为空").Json(ctx)
		return
	}
	commentIdInt, _ := strconv.Atoi(commentId)
	var commentsService services.CommentsService
	commentsService.DeleteComment(commentIdInt, userId)
	result.Ok(nil, "删除成功").Json(ctx)
}

// 返回文章下的评论(文章页面展示)
func listCommentsByArticleId(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("articleId"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "15"))

	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var commentsService services.CommentsService
	comments, count := commentsService.GetCommentsByArticleID(page, limit, articleId)
	result.Ok(map[string]interface{}{
		"data":  &comments,
		"count": &count,
	}, "").Json(ctx)
}

// 查询根评论下的评论
func listCommentsByRootId(ctx *gin.Context) {
	rootId, _ := strconv.Atoi(ctx.Param("rootId"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "15"))
	var commentsService services.CommentsService
	comments, count := commentsService.GetCommentsByRootID(page, limit, rootId)

	result.Ok(map[string]interface{}{
		"data":  comments,
		"count": count,
	}, "").Json(ctx)

}

// 查询文章下的所有评论(管理端)
func listAllCommentsByArticleId(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("articleId"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "15"))

	if err != nil {
		result.Err("文章不可为空").Json(ctx)
		return
	}
	var commentsService services.CommentsService
	comments, count := commentsService.GetAllCommentsByArticleID(page, limit, articleId)
	result.Ok(map[string]interface{}{
		"data":  comments,
		"count": count,
	}, "").Json(ctx)
}
