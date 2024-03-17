package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitCodeRouters(r *gin.Engine) {
	group := r.Group("/community/admin/code")
	group.Use(middleware.OperLogger())
	group.GET("", listCode)
	group.POST("/generate", generate)
	group.DELETE("/:code", deleteCode)
}

func listCode(ctx *gin.Context) {
	var c services.CodeService
	p, limit := page.GetPage(ctx)
	codes, count := c.PageCodes(p, limit)
	result.Ok(page.New(codes, count), "").Json(ctx)
}

func generate(ctx *gin.Context) {
	var c services.CodeService
	var v model.GenerateCode
	if err := ctx.ShouldBindJSON(&v); err != nil {
		log.Warnf("用户id: %d 生成邀请码解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(utils.GetValidateErr(v, err)).Json(ctx)
		return
	}
	if err := c.GenerateCode(v); err != nil {
		log.Warn("用户id: %d 生成邀请码失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "生成成功").Json(ctx)
}

func deleteCode(ctx *gin.Context) {
	code := ctx.Param("code")

	if code == "" {
		log.Warnln("用户id: %d 删除邀请码不存在: %s", middleware.GetUserId(ctx), code)
		result.Err("删除的code不存在").Json(ctx)
		return
	}

	var c services.CodeService

	code1, _ := strconv.Atoi(code)
	if err := c.DestroyCode(code1); err != nil {
		log.Warnf("用户id: %d 删除邀请码失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}
