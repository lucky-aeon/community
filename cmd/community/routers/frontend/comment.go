package frontend

import (
	"fmt"
	"strconv"

	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/constants"
	"xhyovo.cn/community/server/service/event"

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
	group.GET("/byArticleId/:articleId", listCommentsByArticleId)
	group.GET("/byRootId/:rootId", listCommentsByRootId)
	group.GET("/allCommentsByArticleId/:articleId", listAllCommentsByArticleId)
	group.Use(middleware.OperLogger())
	group.POST("/comment", comment)
	group.DELETE("/:id", deleteComment)
	group.POST("/adoption", adoption)

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
	var adS services.QAAdoption
	adS.SetAdoptionComment(comments)
	result.Ok(page.New(comments, count), "").Json(ctx)
}

// 查询根评论下的评论
func listCommentsByRootId(ctx *gin.Context) {
	rootId, _ := strconv.Atoi(ctx.Param("rootId"))
	p, limit := page.GetPage(ctx)
	var commentsService services.CommentsService
	comments, count := commentsService.GetCommentsByRootID(p, limit, rootId)
	var adS services.QAAdoption
	adS.SetAdoptionComment(comments)
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
	var adS services.QAAdoption
	adS.SetAdoptionComment(comments)
	result.Ok(page.New(comments, count), "").Json(ctx)
}

// 采纳评论
func adoption(ctx *gin.Context) {
	var adoption model.QaAdoptions
	userId := middleware.GetUserId(ctx)
	if err := ctx.ShouldBindJSON(&adoption); err != nil {
		log.Warnf("用户id: %d 采纳评论参数解析失败,err :%s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	var cS services.CommentsService
	commentId := adoption.CommentId
	comment := cS.GetById(commentId)
	articleId := comment.BusinessId
	adoption.ArticleId = articleId
	var msg string
	// 采纳权限，文章得是本人,评论得存在
	var aS services.ArticleService
	article := aS.GetById(articleId)
	state := article.State
	if article.UserId != userId {
		msg = fmt.Sprintf("用户id: %d 采纳评论无权限,文章id: %d", userId, articleId)
		log.Warnln(msg)
		result.Err(msg).Json(ctx)
		return
	}
	if !(state == constants.Pending || state == constants.Resolved) {
		result.Err("该文章不是 QA 分类,无法进行采纳").Json(ctx)
		return
	}

	if !aS.Auth(userId, articleId) {
		msg = fmt.Sprintf("用户id: %d 采纳评论无权限,文章id: %d", userId, articleId)
		log.Warnln(msg)
		result.Err(msg).Json(ctx)
		return
	}

	if comment.ID == 0 {
		msg = fmt.Sprintf("用户id: %d 采纳评论对应的评论不存在,评论id: %d", userId, commentId)
		log.Warnln(msg)
		result.Err(msg).Json(ctx)
		return
	}
	var adptionS services.QAAdoption
	msg = "取消采纳"
	if adptionS.Adopt(articleId, commentId) {
		var suS services.SubscriptionService

		suS.Send(event.Adoption, constant.NOTICE, userId, comment.FromUserId, services.SubscribeData{CommentId: commentId, ArticleId: articleId, UserId: userId, CurrentBusinessId: articleId})
		msg = "已采纳"
	}
	// 采纳了,但是状态为未解决,则改为已解决
	if adptionS.QAAdoptState(articleId) && state == constants.Pending {
		aS.UpdateState(articleId, constants.Resolved)
	} else if !adptionS.QAAdoptState(articleId) && state == constants.Resolved {
		aS.UpdateState(articleId, constants.Pending)
	}

	result.OkWithMsg(nil, msg).Json(ctx)
}
