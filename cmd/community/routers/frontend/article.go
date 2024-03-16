package frontend

import (
	"strconv"
	"xhyovo.cn/community/pkg/constant"

	"xhyovo.cn/community/pkg/utils"

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
	UserId  int    `json:"userId"`  // 用户id
	data.ListSortStrategy
}

func InitArticleRouter(r *gin.Engine) {
	group := r.Group("/community/articles")
	group.GET("/:id", articleGet)
	group.POST("", articlePageBySearch)
	group.POST("/update", articleSave)
	group.DELETE("/:id", articleDeleted)
	group.POST("/like", articleLike)
	group.GET("/like/state/:articleId", articleLikeState)

	group.Use(middleware.Auth)
}

func articlePageBySearch(ctx *gin.Context) {
	// 获取所有分类
	searchArticle := new(SearchArticle)
	if err := ctx.ShouldBindBodyWith(searchArticle, binding.JSON); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	result.Page(articleService.PageByClassfily(searchArticle.Tags, &model.Articles{
		Title:   searchArticle.Context,
		Content: searchArticle.Context,
		Type:    searchArticle.Type,
		UserId:  searchArticle.UserId,
	}, ginutils.GetPage(ctx), ginutils.GetOderBy(ctx))).Json(ctx)
}

func articleGet(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil || articleId < 1 {
		result.Err("未找到相关文章").Json(c)
		return
	}
	result.Auto(articleService.GetArticleData(articleId, middleware.GetUserId(c))).ErrMsg("未找到相关文章").Json(c)
}

func articleDeleted(c *gin.Context) {
	id := c.Param("id")
	articleId, _ := strconv.Atoi(id)
	if err := articleService.Delete(articleId, middleware.GetUserId(c)); err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(c)
}

func articleSave(c *gin.Context) {
	var o model.Articles
	if err := c.ShouldBindJSON(&o); err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	o.UserId = middleware.GetUserId(c)
	id, err := articleService.SaveArticle(o)
	if err != nil {
		result.Err(utils.GetValidateErr(o, err)).Json(c)
		return
	}
	articleData, err := articleService.GetArticleData(id, o.UserId)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	result.OkWithMsg(articleData, constant.GetArticleMsg(o.State)).Json(c)
}

func articleLike(c *gin.Context) {
	v := c.Query("articleId")
	articleId, err := strconv.Atoi(v)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	userId := middleware.GetUserId(c)
	var msg string = "取消点赞"
	var likeState bool = false
	if articleService.Like(articleId, userId) {
		msg = "点赞"
		likeState = true
	}
	result.OkWithMsg(likeState, msg).Json(c)
}

func articleLikeState(c *gin.Context) {
	v := c.Param("articleId")
	articleId, err := strconv.Atoi(v)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	userId := middleware.GetUserId(c)
	state := articleService.GetLikeState(articleId, userId)
	result.Ok(state, "").Json(c)
}
