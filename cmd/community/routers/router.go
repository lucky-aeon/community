package routers

import (
	"github.com/gin-gonic/gin"
	services "xhyovo.cn/community/server/service"
)

// init router
var (
	file           services.FileService
	articleService services.ArticleService
	typeService    services.TypeService
	userService    services.UserService

	fileService services.FileService
)

func InitFrontedRouter(r *gin.Engine) {

	InitLoginRegisterRouter(r)
	InitFileRouter(r)
	InitUserRouters(r)
	InitArticleRouter(r)
	InitCommentRouter(r)
}
