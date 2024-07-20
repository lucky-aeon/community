package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var noteService services.RateService

func InitNoteRouters(r *gin.Engine) {
	group := r.Group("/community/rate")
	group.GET("", getRate)
	group.Use(middleware.OperLogger())
	group.POST("", commentRate)
	group.DELETE("/:id", deleteRate)
}

func getRate(ctx *gin.Context) {
	userId := middleware.GetUserId(ctx)

	result.Ok(noteService.GetById(userId), "").Json(ctx)
}

func commentRate(ctx *gin.Context) {
	var note model.Rates
	if err := ctx.ShouldBindJSON(&note); err != nil {
		msg := utils.GetValidateErr(note, err)
		log.Warnf("用户id: %d 保存留言解析失败 ,err: %s", middleware.GetUserId(ctx), msg)
		result.Err(msg).Json(ctx)
		return
	}
	note.UserId = middleware.GetUserId(ctx)
	noteService.Comment(note)
	result.OkWithMsg(nil, "保存成功").Json(ctx)
}

func deleteRate(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	noteService.Delete(id, userId)
	result.OkWithMsg(nil, "删除成功")
}
