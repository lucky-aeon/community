package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

var (
	articleDao = &dao.Article{}
)

func InitArticleRouter(r *gin.Engine) {
	r.GET("/articles", articleList)
	r.POST("/articles", articleAdd)
	r.DELETE("/articles", articleDeleted)
	r.PUT("/articles", articleUpdate)
}

func articleList(c *gin.Context) {
	result, err := articleDao.QueryList(&model.Article{}, c.GetInt("page"), c.GetInt("limit"))
	if err != nil {
		R.Error().setMsg("error querying articles").Res(c)
		return
	}
	R.Ok().setData(result).Res(c)
}

func articleAdd(c *gin.Context) {

}

func articleDeleted(c *gin.Context) {

}

func articleUpdate(c *gin.Context) {

}
