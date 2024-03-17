package routers

import (
	"strconv"
	"xhyovo.cn/community/pkg/log"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitCommentRouters(g *gin.Engine) {
	group := g.Group("/community/comments")
	group.Use(middleware.OperLogger())
	group.POST("/comment", comment, middleware.OperLogger())
	group.DELETE("/:id", deleteComment, middleware.OperLogger())
	group.GET("/byArticleId/:articleId", listCommentsByArticleId)
	group.GET("/byRootId/:rootId", listCommentsByRootId)
	group.GET("/allCommentsByArticleId/:articleId", listAllCommentsByArticleId)
}

// 发布评论
func comment(ctx *gin.Context) {
	var comment model.Comments
	userId := middleware.GetUserId(ctx)
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		log.Warnf("用户id: %d 发布评论失败,err: %s", userId, err.Error())
		result.Err(utils.GetValidateErr(comment, err)).Json(ctx)
		return
	}
	comment.FromUserId = userId

	commentsService := services.NewCommentService(ctx)
	err := commentsService.Comment(&comment)
	msg := "评论成功"
	if err != nil {
		log.Warnf("用户id: %d 保存评论失败,err: %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, msg).Json(ctx)
}

// 删除评论
func deleteComment(ctx *gin.Context) {
	commentId := ctx.Param("id")

	userId := middleware.GetUserId(ctx)
	if commentId == "" {
		log.Warnf("用户id: %d 删除评论失败,err: %s", userId, "评论id为空")
		result.Err("删除评论id不能为空").Json(ctx)
		return
	}
	commentIdInt, _ := strconv.Atoi(commentId)
	var commentsService services.CommentsService
	if !commentsService.DeleteComment(commentIdInt, userId) {
		log.Warnf("用户id: %d 删除评论失败", userId)
		result.Err("删除失败").Json(ctx)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}

// 返回文章下的评论(文章页面展示)
func listCommentsByArticleId(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("articleId"))
	p, limit := page.GetPage(ctx)

	if err != nil {
		log.Warnf("用户id: %d 获取文章下的评论失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	var commentsService services.CommentsService
	comments, count := commentsService.GetCommentsByArticleID(p, limit, articleId)
	result.Ok(page.New(comments, count), "").Json(ctx)
}

// 查询根评论下的评论
func listCommentsByRootId(ctx *gin.Context) {
	rootId, _ := strconv.Atoi(ctx.Param("rootId"))
	p, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "15"))
	var commentsService services.CommentsService
	comments, count := commentsService.GetCommentsByRootID(p, limit, rootId)

	result.Ok(page.New(comments, count), "").Json(ctx)

}

// 查询用户文章下的所有评论，文章id为空则查询所有(管理端)
func listAllCommentsByArticleId(ctx *gin.Context) {
	articleId, err := strconv.Atoi(ctx.Param("articleId"))
	p, limit := page.GetPage(ctx)

	if err != nil {
		log.Warnf("用户id: %d 查询用户文章下的所有评论失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err("文章不可为空").Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	var commentsService services.CommentsService
	comments, count := commentsService.GetAllCommentsByArticleID(p, userId, limit, articleId)

	result.Ok(page.New(comments, count), "").Json(ctx)
}
