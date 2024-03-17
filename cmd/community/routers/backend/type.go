package backend

import (
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitTypeRouters(r *gin.Engine) {
	group := r.Group("/community/admin/types")
	group.GET("", listType)
	group.POST("", saveType)
	group.PUT("", UpdateType)
	group.DELETE(":id", deleteType)
}

func listType(ctx *gin.Context) {
	parentId, err := strconv.Atoi(ctx.DefaultQuery("parentId", "0"))
	if err != nil {
		parentId = 0
	}
	var typeService services.TypeService
	result.Ok(typeService.List(parentId), "").Json(ctx)
}

func saveType(ctx *gin.Context) {

	var types model.Types
	if err := ctx.ShouldBindJSON(&types); err != nil {
		log.Warnf("用户id: %d 添加分类参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(utils.GetValidateErr(types, err)).Json(ctx)
		return
	}
	var typeService services.TypeService
	u, err := typeService.Save(&types)
	if err != nil {
		log.Warnf("用户id: %d 添加分类失败,err: %s", middleware.GetUserId(ctx), err.Error())
	}
	result.Auto(u, err).ErrMsg("添加失败").OkMsg("添加成功").Json(ctx)
}
func UpdateType(ctx *gin.Context) {
	var types model.Types
	if err := ctx.ShouldBindJSON(&types); err != nil {
		log.Warnf("用户id: %d 修改分类参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(utils.GetValidateErr(types, err)).Json(ctx)
		return
	}
	var typeService services.TypeService
	err := typeService.Update(&types)
	if err != nil {
		log.Warnf("用户id: %d 修改分类失败,err: %s", middleware.GetUserId(ctx), err.Error())
	}
	result.Auto(nil, err).OkMsg("分类更新成功").Json(ctx)
}

func deleteType(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warnf("用户id: %d 删除分类失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	var typeService services.TypeService
	var articleService services.ArticleService
	// 如果有文章
	count := articleService.CountByTypeId(id)
	if count > 0 {
		var msg = "删除失败,该分类下有文章"
		log.Warnf("用户id: %d 删除分类失败,err: %s", middleware.GetUserId(ctx), msg)
		result.Err(msg).Json(ctx)
		return
	}
	err = typeService.Delete(id)
	if err != nil {
		log.Warnf("用户id: %d 删除分类失败,err: %s", middleware.GetUserId(ctx), err.Error())
	}
	result.Auto(nil, err).Json(ctx)
}
