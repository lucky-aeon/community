package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/cmd/community/routers/backend"
	"xhyovo.cn/community/cmd/community/routers/frontend"
)

// init router

func InitFrontedRouter(r *gin.Engine) {

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
