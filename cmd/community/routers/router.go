package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/routers/backend"
	"xhyovo.cn/community/cmd/community/routers/frontend"
)

// init router

func InitFrontedRouter(r *gin.Engine) {

	InitLoginRegisterRouters(r)
	//r.Use(middleware.Auth) todo
	InitFileRouters(r)
	frontend.InitUserRouters(r)
	frontend.InitArticleRouter(r)
	frontend.InitTypeRouters(r)
	frontend.InitSubscriptionRouters(r)
	backend.InitTypeRouters(r)
	backend.InitCodeRouters(r)
	InitCommentRouters(r)

}
