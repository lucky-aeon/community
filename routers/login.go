package routers

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/services"
	"xhyovo.cn/community/utils"
)

// todo: feat login
func Login(c *gin.Context) {

	user := services.Login(c.Query("account"), c.Query("password"))
	user.Password = ""
	utils.Ok().Data(user).Res(c)
}

// todo: feat register
func Register(c *gin.Context) {
	// services.
}
