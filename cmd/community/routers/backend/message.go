package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitMessageRouters(r *gin.Engine) {
	group := r.Group("/community/admin/message")
	group.GET("/template/var", listMsgVar)
	group.GET("/template", listMsgTemp)
	group.POST("/template", saveMsgTemp)
	group.DELETE("/template", deleteMsgTemp)
}

// 获取消息模板中的变量
func listMsgVar(ctx *gin.Context) {

	var mS services.MessageService

	result.Ok(mS.GetMessageTemplateVar(), "").Json(ctx)
}

func listMsgTemp(ctx *gin.Context) {
	var mS services.MessageService
	template, count := mS.ListMessageTemplate(page.GetPage(ctx))
	result.Page(template, count, nil).Json(ctx)
}

func saveMsgTemp(ctx *gin.Context) {
	var mS services.MessageService
	var template model.MessageTemplates
	if err := ctx.ShouldBindJSON(&template); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	mS.SaveMessageTemplate(template)
	result.OkWithMsg(nil, "保存成功").Json(ctx)
}

func deleteMsgTemp(ctx *gin.Context) {
	var mS services.MessageService
	id := ctx.Param("id")
	atoi, err := strconv.Atoi(id)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	mS.DeleteMessageTemplate(atoi)
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}
