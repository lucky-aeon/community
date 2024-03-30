package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/request"
	services "xhyovo.cn/community/server/service"
)

func InitArticleRouters(r *gin.Engine) {
	group := r.Group("/community/admin/article")
	group.GET("/page", listArticles)
	group.DELETE("/:id", deleteArticle, middleware.OperLogger())
	group.POST("/state", articleState, middleware.OperLogger())
	group.GET("/states", listStates)
}

func listArticles(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	var a services.ArticleService
	articles, count := a.PageArticles(p, limit)
	result.Page(articles, count, nil).Json(ctx)
}

func deleteArticle(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warnf("删除文章时参数解析失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	var a services.ArticleService
	if err := a.Delete(id); err != nil {
		log.Warnf("删除文章失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}

func articleState(ctx *gin.Context) {
	var topArticle request.TopArticle
	if err := ctx.ShouldBindJSON(&topArticle); err != nil {
		log.Warnf("修改文章状态时参数解析失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	var a services.ArticleService
	if err := a.UpdateArticleState(topArticle); err != nil {
		log.Warnf("修改文章状态失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "修改成功").Json(ctx)
}

func listStates(ctx *gin.Context) {
	result.Ok(constant.ListState(), "").Json(ctx)
}
