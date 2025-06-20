package routers

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
)

func InitSsoRouters(ctx *gin.Engine) {
	// 公开的SSO接口，不需要认证
	group := ctx.Group("/sso")
	group.GET("/login", SsoLogin)
	group.POST("/token", SsoToken)
	
	// 兼容其他路径格式
	apiGroup := ctx.Group("/api/sso/community")
	apiGroup.GET("/login", SsoLogin)
	apiGroup.POST("/token", SsoToken)
}

// SsoLogin SSO登录入口
func SsoLogin(c *gin.Context) {
	appKey := c.Query("app_key")
	redirectUrl := c.Query("redirect_url")

	if appKey == "" || redirectUrl == "" {
		result.Err("缺少必要参数").Json(c)
		return
	}

	var ssoService services.SsoService

	// 验证应用
	app, err := ssoService.GetApplicationByKey(appKey)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}

	// 验证回调地址
	if !ssoService.ValidateRedirectUrl(app, redirectUrl) {
		result.Err("回调地址未授权").Json(c)
		return
	}

	// 检查用户是否已登录
	token := c.GetHeader("Authorization")
	if len(token) == 0 {
		token, _ = c.Cookie("Authorization")
	}

	claims, err := middleware.ParseToken(token)
	if err != nil || claims.ID < 1 {
		// 未登录，重定向到登录页面，并携带回跳参数
		loginUrl := "/login?sso=1&app_key=" + url.QueryEscape(appKey) + "&redirect_url=" + url.QueryEscape(redirectUrl)
		c.Redirect(302, loginUrl)
		return
	}

	// 已登录，生成授权码
	authCode, err := ssoService.GenerateAuthCode(appKey, claims.ID, redirectUrl)
	if err != nil {
		result.Err("生成授权码失败").Json(c)
		return
	}

	// 重定向到第三方应用
	finalUrl := redirectUrl + "?code=" + authCode
	c.Redirect(302, finalUrl)
}

// SsoToken 获取用户信息
func SsoToken(c *gin.Context) {
	var req struct {
		AppKey    string `json:"app_key" binding:"required"`
		AppSecret string `json:"app_secret" binding:"required"`
		AuthCode  string `json:"auth_code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误").Json(c)
		return
	}

	var ssoService services.SsoService

	// 验证授权码并获取用户信息
	user, err := ssoService.ValidateAuthCode(req.AppKey, req.AppSecret, req.AuthCode)
	if err != nil {
		result.Err(err.Error()).Json(c)
		return
	}

	// 返回用户基本信息
	userInfo := ssoService.GetUserBasicInfo(user)
	result.Ok(userInfo, "").Json(c)
}
