package frontend

import (
	"strconv"

	"xhyovo.cn/community/pkg/utils/page"

	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/request"

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
	Tags    []string `json:"tags"`    // 文章标签
	Context string   `json:"context"` // 模糊查询内容
	Type    string   `json:"type"`    // 分类id
	UserId  int      `json:"userId"`  // 用户id
	State   int      `json:"state"`
	data.ListSortStrategy
}

func InitArticleRouter(r *gin.Engine) {
	group := r.Group("/community/articles")

	group.GET("/:id", articleGet)
	group.GET("/top", articleTop)
	group.POST("", articlePageBySearch)
	group.GET("/like/state/:articleId", articleLikeState)
	group.Use(middleware.OperLogger())
	group.POST("/update", articleSave)
	group.DELETE("/:id", articleDeleted)
	group.POST("/like", articleLike)

}

func articlePageBySearch(ctx *gin.Context) {
	// 获取所有分类
	searchArticle := new(SearchArticle)
	if err := ctx.ShouldBindBodyWith(searchArticle, binding.JSON); err != nil {
		log.Warnf("用户id: %d 分页获取文章参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	state := searchArticle.State
	if state < 1 || state > 6 {
		log.Warnf("用户id: %d 搜索文章状态参数错误,当前状态: %d", middleware.GetUserId(ctx), state)
		result.Err("文章状态非法").Json(ctx)
		return
	}

	searchUserId := searchArticle.UserId
	currentUserId := middleware.GetUserId(ctx)

	if (state == constant.Draft || state == constant.QADraft || state == constant.PrivateQuestion) && searchUserId != 0 && searchUserId != currentUserId {
		log.Warnf("用户id: %d 搜索文章状态不可选择草稿以及私密提问", middleware.GetUserId(ctx))
		result.Err("搜索文章状态不可选择草稿以及私密提问").Json(ctx) //
		return
	}
	if state == 0 {
		log.Warnf("用户id: %d 查询文章必须带上文章状态", middleware.GetUserId(ctx))
		result.Err("查询文章必须带上文章状态").Json(ctx)
		return
	}
	var userS services.UserService
	flag, err := userS.IsAdmin(currentUserId)
	if err != nil {
		log.Warnf("用户id: %d 校验身份出现错误: %s", middleware.GetUserId(ctx), err)
		result.Err("校验身份出现错误").Json(ctx)
		return
	}

	// TA 用户并且 不是管理员
	if searchUserId != currentUserId && !flag && (state == constant.Draft || state == constant.QADraft || state == constant.PrivateQuestion) {
		log.Warnf("用户id: %d 非法查询文章,查询文章状态: %s", middleware.GetUserId(ctx), state)
		result.Err("你没有权限查询该状态文章").Json(ctx)
		return
	}

	result.Page(articleService.PageByClassfily(searchArticle.Type, searchArticle.Tags, &model.Articles{
		Title:   searchArticle.Context,
		Content: searchArticle.Context,
		UserId:  searchUserId,
		State:   state,
	}, ginutils.GetPage(ctx), ginutils.GetOderBy(ctx), currentUserId)).Json(ctx)
}

func articleGet(c *gin.Context) {
	articleId, err := strconv.Atoi(c.Params.ByName("id"))
	if err != nil || articleId < 1 {
		log.Warnf("用户id: %d 未找到相关文章,文章id: %d err: %s", middleware.GetUserId(c), articleId, err.Error())
		result.Err("未找到相关文章").Json(c)
		return
	}
	result.Auto(articleService.GetArticleData(articleId, middleware.GetUserId(c))).ErrMsg("未找到相关文章").Json(c)
}

func articleDeleted(c *gin.Context) {
	id := c.Param("id")
	articleId, _ := strconv.Atoi(id)
	if err := articleService.DeleteByUserId(articleId, middleware.GetUserId(c)); err != nil {
		log.Warnf("用户id: %d 删除文章失败,文章id: %d ,err: %s", middleware.GetUserId(c), articleId, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(c)
}

func articleSave(c *gin.Context) {
	var o request.ReqArticle
	if err := c.ShouldBindJSON(&o); err != nil {
		log.Warnf("用户id: %d 保存文章解析文章失败 ,err: %s", middleware.GetUserId(c), err.Error())
		result.Err(err.Error()).Json(c)
		return
	}
	o.UserId = middleware.GetUserId(c)
	id, err := articleService.SaveArticle(o)
	if err != nil {
		log.Warnf("用户id: %d 保存文章失败,err: %s", middleware.GetUserId(c), err.Error())
		result.Err(utils.GetValidateErr(o, err)).Json(c)
		return
	}
	articleData, err := articleService.GetArticleData(id, o.UserId)
	if err != nil {
		log.Warnf("用户id: %d 获取文章失败,文章id: %d ,err: %s", middleware.GetUserId(c), id, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}
	result.OkWithMsg(articleData, constant.GetArticleMsg(o.State)).Json(c)
}

func articleLike(c *gin.Context) {
	v := c.Query("articleId")
	articleId, err := strconv.Atoi(v)
	if err != nil {
		log.Warnf("用户id: %d 点赞文章失败,文章id: %d ,err: %s", middleware.GetUserId(c), articleId, err.Error())
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
		log.Warnf("用户id: %d 获取文章点赞状态解析id失败,文章id: %d ,err: %s", middleware.GetUserId(c), articleId, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}
	userId := middleware.GetUserId(c)
	state := articleService.GetLikeState(articleId, userId)
	result.Ok(state, "").Json(c)
}

func articleTop(ctx *gin.Context) {
	types := ctx.Query("type")
	p, limit := page.GetPage(ctx)
	articles, count := articleService.PageTopArticle(types, p, limit)
	result.Page(articles, count, nil).Json(ctx)
}
