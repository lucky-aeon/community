package middleware

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/service_context"
)

func Auth(ctx *gin.Context) {
	context := service_context.DataContext(ctx)
	if !context.Check() {
		context.To("/login").WithError("请登录以后再访问").Redirect()
		return
	}
}
