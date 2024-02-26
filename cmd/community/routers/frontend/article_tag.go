package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	ginutils "xhyovo.cn/community/pkg/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var (
	articleTagService = new(services.ArticleTagService)
)

func InitArticleTagRouter(r *gin.Engine) {
	rg := r.Group("/community/tags")
	rg.GET("", getArticleTags)
	rg.GET("/hot", getHotTags)
	rg.POST("", saveArticleTags)
	rg.DELETE("/:tagId", deleteArticleTags)
	rg.GET("/getTagArticleCount", getTagArticleCount)
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

func saveArticleTags(c *gin.Context) {
	var articleTag model.ArticleTags
	if err := c.ShouldBindJSON(&articleTag); err != nil {
		result.Err(utils.GetValidateErr(articleTag, err)).Json(c)
		return
	}
	articleTag.UserId = middleware.GetUserId(c)
	tag, err := articleTagService.CreateTag(articleTag)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	result.Ok(tag, "").Json(c)
}

func deleteArticleTags(c *gin.Context) {
	tagId := c.Param("tagId")
	atoi, _ := strconv.Atoi(tagId)
	userId := middleware.GetUserId(c)

	if err := articleTagService.DeleteTag(atoi, userId); err != nil {
		result.Err(err.Error()).Json(c)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(c)
}

// 获取标签引用的文章
func getTagArticleCount(c *gin.Context) {
	userId := middleware.GetUserId(c)
	tagArticleCount := articleTagService.GetTagArticleCount(userId)
	result.Ok(tagArticleCount, "").Json(c)
}
