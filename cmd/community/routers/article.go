package routers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

func InitArticleRouter(r *gin.Engine) {
	group := r.Group("/community")
	group.GET("/articles/:id", articleGet)
	group.GET("/articles", articleList)
	group.GET("/articles/search", articleSearch)
	group.POST("/articles/update", articleUpdate)
	group.DELETE("/articles/:id", articleDeleted)
	group.Use(middleware.Auth)
}

func articleList(ctx *gin.Context) {
	// 获取所有分类

}

func articleGet(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Params.ByName("articleId"))
	if err != nil || articleId < 1 {
		result.Err("未找到相关文章").Json(c)
		return
	}
	daoArticle := &dao.Article{}
	r, err := daoArticle.QuerySingle(model.Articles{
		Model: gorm.Model{
			ID: uint(articleId),
		},
	})
	if err != nil {
		result.Err("未找到相关文章").Json(c)
		return
	}
	result.Ok(r, "未找到相关文章").Json(c)
}

func articleSearch(c *gin.Context) {

}

func articleDeleted(c *gin.Context) {

}

func articleUpdate(c *gin.Context) {

}
