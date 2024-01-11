package routers

import (
	"github.com/gin-gonic/gin"
)

// init router
func InitRouter() {
	r := gin.Default()
	r.GET("/login", Login)
	r.GET("/register", Register)
	r.Run()
}
