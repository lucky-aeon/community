package frontend

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitDraftRouters(r *gin.Engine) {
	group := r.Group("/community/draft")
	group.GET("", getDraft)
	group.POST("", saveDraft)
	//group.GET("/init", initDraft)
}

func initDraft(ctx *gin.Context) {

	// 获取所有用户
	var userIds []int
	model.User().Select("id").Find(&userIds)

	var drafts = make([]model.Drafts, 0, len(userIds))
	for i := range userIds {
		drafts = append(drafts, model.Drafts{UserId: userIds[i], State: 2})
	}
	model.Draft().Create(drafts)
	result.Ok(nil, "").Json(ctx)
}

func getDraft(ctx *gin.Context) {
	var d services.Draft
	userId := middleware.GetUserId(ctx)
	draft := d.Get(userId)
	result.Ok(draft, "").Json(ctx)
}

// 如何实现临时保存文本，和文章相关联
func saveDraft(ctx *gin.Context) {
	var draft model.Drafts
	userId := middleware.GetUserId(ctx)
	if err := ctx.ShouldBindJSON(&draft); err != nil {
		log.Warnf("用户id: %d,临时存储文章参数解析错误,err %s", userId, err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	draft.UserId = userId
	var d services.Draft
	d.Save(draft)
	result.Ok(nil, "").Json(ctx)

}
