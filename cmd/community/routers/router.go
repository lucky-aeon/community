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
	frontend.InitFileRouters(r)
	r.Use(middleware.Auth)
	frontend.InitUserRouters(r)
	frontend.InitArticleRouter(r)
	frontend.InitTypeRouters(r)
	frontend.InitSubscriptionRouters(r)
	frontend.InitMessageRouters(r)
	frontend.InitCommentRouters(r)
	frontend.InitArticleTagRouter(r)
	frontend.InitDraftRouters(r)
	frontend.InitCourseRouters(r)

	r.Use(middleware.AdminAuth)
	backend.InitTypeRouters(r)
	backend.InitCodeRouters(r)
	backend.InitFileRouters(r)
	backend.InitCommentRouters(r)
	backend.InitUserRouters(r)
	backend.InitMemberRouters(r)
	backend.InitMessageRouters(r)
	backend.InitLogRouters(r)
	backend.InitArticleRouters(r)
	backend.InitUserTagRouters(r)
	backend.InitCourseRouters(r)

}
