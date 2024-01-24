package routers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/service_context"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var (
	article services.ArticleService
)

func InitArticleRouter(r *gin.Engine) {
	group := r.Group("/community")
	group.GET("/articles/:id", articleList)
	group.GET("/articles/list", articleList)
	group.GET("/articles/search", articleSearch)
	group.POST("/articles/publish", articleAdd)
	group.POST("/articles/update", articleUpdate)
	group.POST("/articles/delete/:id", articleDeleted)
	// ------ View ------
	group.GET("/articles/edit/:articleId", articleEdit)
	group.GET("/articles/publish", articleAdd)
	group.Use(middleware.Auth)
}

func articleList(ctx *gin.Context) {

}

func articleGet(c *gin.Context) {

}

func articleSearch(c *gin.Context) {

}

func articleAdd(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "article.publish", nil)
		return
	}
}

func articleDeleted(c *gin.Context) {

}

func articleUpdate(c *gin.Context) {

}

func articleEdit(c *gin.Context) {
	bc := service_context.DataContext(c)
	articleId, err := strconv.Atoi(c.Params.ByName("articleId"))
	if err != nil || articleId < 1 {
		c.Redirect(302, "/") // 非法操作，直接首页
		return
	}
	daoArticle := &dao.Article{}
	result, err := daoArticle.QuerySingle(model.Articles{
		Model: gorm.Model{
			ID: uint(articleId),
		},
		UserId: bc.Auth().ID,
	})
	if err != nil {
		bc.WithError("未找到相关文章").Referer()
		return
	}
	log.Println(articleId, result)
	c.HTML(http.StatusOK, "article.publish", nil)
}
