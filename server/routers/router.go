package routers

import (
	"github.com/gin-gonic/gin"
)

// init router
func InitRouter() {
	r := gin.Default()

	r.POST("/login", Login)
	r.POST("/register", Register)
	r.GET("/test", func(ctx *gin.Context) {
		ctx.JSON(200, R.Ok().setMsg("ok"))
	})
	r.Run()
}
