package backend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/utils/page"
	services "xhyovo.cn/community/server/service"
)

func InitArticleRouters(r *gin.Engine) {

}

func ListArticles(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	var a services.ArticleService
	a.PageArticles(p, limit)
}

func UpdateArticle(ctx *gin.Context) {

}

func DeleteArticle(ctx *gin.Context) {

}
