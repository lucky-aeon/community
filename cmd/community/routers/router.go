package routers

import (
	"os"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/cmd/community/routers/backend"
	"xhyovo.cn/community/cmd/community/routers/frontend"
)

// init router

func InitFrontedRouter(r *gin.Engine) {
	fileInfo, err := os.Stat("./web/assets")
	if err == nil && fileInfo.IsDir() {
		r.Static("/assets", "./web/assets")
	}
	fileInfo, err = os.Stat("./web/index.html")
	if err == nil && !fileInfo.IsDir() {
		r.LoadHTMLFiles("./web/index.html")
		r.GET("/", func(ctx *gin.Context) {
			ctx.HTML(200, "index.html", nil)
		})
	}

	InitLoginRegisterRouters(r)
	r.Use(middleware.Auth)
	r.Use(middleware.OperLogger())
	InitFileRouters(r)
	frontend.InitUserRouters(r)
	frontend.InitArticleRouter(r)
	frontend.InitTypeRouters(r)
	frontend.InitSubscriptionRouters(r)
	frontend.InitMessageRouters(r)
	backend.InitTypeRouters(r)
	backend.InitCodeRouters(r)
	backend.InitFileRouters(r)
	backend.InitCommentRouters(r)
	backend.InitUserRouters(r)
	backend.InitMemberRouters(r)
	backend.InitMessageRouters(r)
	backend.InitLogRouters(r)
	InitCommentRouters(r)
	frontend.InitArticleTagRouter(r)
}
