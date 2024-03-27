package middleware

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

func AdminAuth(ctx *gin.Context) {
	userId := GetUserId(ctx)
	var uS services.UserService
	flag, err := uS.IsAdmin(userId)
	if err != nil {
		result.Err("判断 admin 失败").Json(ctx)
		ctx.Abort()
		return
	}
	if !flag {
		result.Err("无权限").Json(ctx)
		ctx.Abort()
	}
	ctx.Next()
}
