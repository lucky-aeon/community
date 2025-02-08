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
	InitIndexRouters(r)
	frontend.InitFileRouters(r)
	r.Use(middleware.Auth)
	frontend.InitUserRouters(r)
	frontend.InitChatRouter(r)
	frontend.InitArticleRouter(r)
	frontend.InitTypeRouters(r)
	frontend.InitSubscriptionRouters(r)
	frontend.InitMessageRouters(r)
	frontend.InitCommentRouters(r)
	frontend.InitArticleTagRouter(r)
	frontend.InitDraftRouters(r)
	frontend.InitCourseRouters(r)
	frontend.InitNoteRouters(r)
	frontend.InitMeetingRouters(r)
	frontend.InitKnowledgeRouters(r)

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
	backend.InitDashboardRouters(r)
	backend.InitOrderRouters(r)
	backend.InitMonitRouters(r)
	backend.InitMeetingRouters(r)
	backend.InitActivityRouters(r)

}
