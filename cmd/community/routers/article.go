package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	services "xhyovo.cn/community/server/service"
)

var (
	article services.ArticleService
)

func InitArticleRouter(r *gin.Engine) {
	group := r.Group("/community/articles")
	group.GET("/articles/:id", articleList)
	group.GET("/articles/list", articleList)
	group.GET("/articles/search", articleSearch)
	group.POST("/articles/publish", articleAdd)
	group.POST("/articles/update", articleUpdate)
	group.POST("/articles/delete/:id", articleDeleted)
	group.Use(middleware.Auth)
}

func articleList(ctx *gin.Context) {

}

func articleGet(c *gin.Context) {

}

func articleSearch(c *gin.Context) {

}

func articleAdd(c *gin.Context) {

}

func articleDeleted(c *gin.Context) {

}

func articleUpdate(c *gin.Context) {

}
