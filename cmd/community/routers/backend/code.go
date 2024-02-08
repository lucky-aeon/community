package backend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitCodeRouters(r *gin.Engine) {
	group := r.Group("/community/admin/code")
	group.GET("", listCode)
	group.POST("/generate", saveCode)
	group.DELETE("/:code", deleteCode)
}

func listCode(ctx *gin.Context) {
	var c services.CodeService
	p, limit := page.GetPage(ctx)
	codes, count := c.PageCodes(p, limit)
	result.Ok(page.New(codes, count), "").Json(ctx)
}

func saveCode(ctx *gin.Context) {
	var c services.CodeService
	var v model.GenerateCode
	if err := ctx.ShouldBindJSON(&v); err != nil {
		result.Err(utils.GetValidateErr(v, err)).Json(ctx)
		return
	}
	if err := c.GenerateCode(v); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.Ok(nil, "生成成功").Json(ctx)
}

func deleteCode(ctx *gin.Context) {
	code := ctx.Param("code")

	if code == "" {
		result.Err("删除的code不存在").Json(ctx)
		return
	}

	var c services.CodeService

	code1, _ := strconv.Atoi(code)
	if err := c.DestroyCode(code1); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.Ok(nil, "删除成功").Json(ctx)
}
