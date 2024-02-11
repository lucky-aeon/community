package frontend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/data"
	ginutils "xhyovo.cn/community/pkg/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var (
	articleService = new(services.ArticleService)
)

type SearchArticle struct {
	Tags    []int  `json:"tags"`    // 文章标签
	Context string `json:"context"` // 模糊查询内容
	Type    int    `json:"type"`    // 分类id
	data.ListSortStrategy
}

func InitArticleRouter(r *gin.Engine) {
	group := r.Group("/community")
	group.GET("/articles/:id", articleGet)
	group.GET("/articles", articlePageBySearch)
	group.POST("/articles/update", articleUpdate)
	group.DELETE("/articles/:id", articleDeleted)
	group.Use(middleware.Auth)
}

func articlePageBySearch(ctx *gin.Context) {
	// 获取所有分类
	searchArticle := new(SearchArticle)
	ctx.ShouldBindBodyWith(searchArticle, binding.JSON)

	result.Page(articleService.PageByClassfily(searchArticle.Tags, &model.Articles{
		Title: searchArticle.Context,
		Desc:  searchArticle.Context,
		Type:  searchArticle.Type,
	}, ginutils.GetPage(ctx), ginutils.GetOderBy(ctx))).Json(ctx)
}

func articleGet(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil || articleId < 1 {
		result.Err("未找到相关文章").Json(c)
		return
	}
	result.Auto(articleService.GetArticleData(articleId)).ErrMsg("未找到相关文章").Json(c)
}

func articleDeleted(c *gin.Context) {

}

func articleUpdate(c *gin.Context) {

}
