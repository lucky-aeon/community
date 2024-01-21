package routers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/service_context"
	services "xhyovo.cn/community/server/service"
)

var typeService services.TypeService

func register(ctx *gin.Context) {
	context := service_context.DataContext(ctx)
	context.View("login.register", nil)

}
func login(ctx *gin.Context) {
	service_context.DataContext(ctx).View("login.login", nil)

}
func logout(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "login.login", nil)
}

func index(ctx *gin.Context) {
	context := service_context.DataContext(ctx)
	if !context.Check() {
		context.To("/login").Redirect()
		return
	}

	// 获取所有分类
	types := typeService.List()

	// 获取所有文章
	var articleService services.ArticleService
	data := articleService.Page(context)

	if data != nil {
		data["types"] = types
	}
	// todo 获取额外信息
	context.View("home.index", data)
}

func edit(ctx *gin.Context) {
	context := service_context.DataContext(ctx)
	f := ctx.DefaultQuery("tab", "info")
	context.View("user.edit", gin.H{"tab": f})
}

func InitViewsRouters(r *gin.Engine) {

	r.GET("/register", register)
	r.GET("/login", login)
	//r.GET("/logout", logout)
	r.GET("/", index)

	group := r.Group("/community")
	group.GET("/user/edit", edit)
	group.Use(middleware.Auth)
}
