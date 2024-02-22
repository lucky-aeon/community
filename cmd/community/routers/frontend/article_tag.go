package frontend

import (
	"github.com/gin-gonic/gin"
	ginutils "xhyovo.cn/community/pkg/gin"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

var (
	articleTagService = new(services.ArticleTagService)
)

func InitArticleTagRouter(r *gin.Engine) {
	rg := r.Group("/community/tags")
	rg.GET("", getArticleTags)
	rg.GET("/hot", getHotTags)
}

func getArticleTags(r *gin.Context) {
	qp := ginutils.GetPage(r)
	title := r.DefaultQuery("title", "")
	result.Auto(articleTagService.QueryList(qp.Page, qp.Limit, title)).Json(r)
}

func getHotTags(r *gin.Context) {
	qp := ginutils.GetPage(r)
	result.Auto(articleTagService.QueryHotTags(qp.Limit)).Json(r)
}
