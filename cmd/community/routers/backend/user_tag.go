package backend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/request"
	services "xhyovo.cn/community/server/service"
)

var userTagS services.UserTag

func InitUserTagRouters(r *gin.Engine) {
	group := r.Group("/community/admin/user/tag")
	group.GET("", listUserTags)
	group.POST("", saveUserTag)
	group.DELETE("/:id", deleteUserTag)
	group.POST("/assignUserLabel", assignUserLabel)
	group.GET("/:userId", getTagsByUserId)

}

func listUserTags(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	tags, count := userTagS.Page(p, limit)
	result.Page(tags, count, nil).Json(ctx)
}
func saveUserTag(ctx *gin.Context) {
	var userTag model.UserTags
	if err := ctx.ShouldBindJSON(&userTag); err != nil {
		log.Warnf("保护用户标签参数解析失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	userTagS.Save(userTag)
	result.OkWithMsg(nil, "保存成功").Json(ctx)
}

func deleteUserTag(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Warnf("删除用户标签参数解析失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	userTagS.DeleteById(id)
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}

func assignUserLabel(ctx *gin.Context) {
	var userTags request.UserTags
	if err := ctx.ShouldBindJSON(&userTags); err != nil {
		log.Warnf("分配用户标签参数解析失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	tagNames := userTagS.AssignUserLabel(userTags.UserId, userTags.TagsIds)
	result.OkWithMsg(tagNames, "分配成功").Json(ctx)
}
func getTagsByUserId(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		log.Warnf("获取用户标签参数解析失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	tagNames := userTagS.GetTagsByUserId(userId)
	result.Ok(tagNames, "").Json(ctx)
}
