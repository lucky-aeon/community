package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
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
}

func getDraft(ctx *gin.Context) {
	var d services.Draft
	articleId, err := strconv.Atoi(ctx.Query("articleId"))
	if err != nil {
		articleId = 0
	}
	userId := middleware.GetUserId(ctx)
	draft := d.Get(userId, articleId)
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
