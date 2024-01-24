package routers

import (
	"github.com/gin-gonic/gin"
)

// init router

func InitFrontedRouter(r *gin.Engine) {

	InitLoginRegisterRouter(r)
	InitFileRouter(r)
	InitUserRouters(r)
	InitArticleRouter(r)
}
