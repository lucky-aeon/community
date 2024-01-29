package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitTypeRouters(r *gin.Engine) {
	group := r.Group("/community/admin/types")
	group.GET("/", list)
	group.POST("/", save)
	group.PUT("/", update)
	group.DELETE("/:id", delete)
}

func list(ctx *gin.Context) {
	parentId, err := strconv.Atoi(ctx.DefaultQuery("parentId", "0"))
	if err != nil {
		parentId = 0
	}
	var typeService services.TypeService
	result.Ok(typeService.List(uint(parentId)), "").Json(ctx)
}

func save(ctx *gin.Context) {

	var types model.Types
	if err := ctx.ShouldBindJSON(&types); err != nil {
		result.Err(utils.GetValidateErr(types, err)).Json(ctx)
		return
	}
	var typeService services.TypeService
	typeService.Save(&types)
	result.Ok(nil, "添加成功").Json(ctx)
}
func update(ctx *gin.Context) {
	var types model.Types
	if err := ctx.ShouldBindJSON(&types); err != nil {
		result.Err(utils.GetValidateErr(types, err)).Json(ctx)
		return
	}
	var typeService services.TypeService
	typeService.Update(&types)
	result.Ok(nil, "修改成功").Json(ctx)
}

func delete(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var typeService services.TypeService
	var articleService services.ArticleService
	// 如果有文章
	count := articleService.CountByTypeId(id)
	if count > 0 {
		result.Err("删除失败,该分类下有文章").Json(ctx)
		return
	}
	typeService.Delete(uint(id))
	result.Ok(nil, "删除成功").Json(ctx)
}
