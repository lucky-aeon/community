package routers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// init router
func InitRouter() {
	r := gin.Default()
	store := cookie.NewStore([]byte("sadasdsadsadas"))
	r.Use(sessions.Sessions("community", store))
	r.POST("/login", Login)
	r.POST("/register", Register)
	r.GET("/test", GetInfo)
	r.Run()
}
