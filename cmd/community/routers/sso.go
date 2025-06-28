package routers

import (
	"github.com/gin-gonic/gin"
	"net/url"
	"strings"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
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
	log.Info("触发 sso")

	if appKey == "" || redirectUrl == "" {
		result.Err("缺少必要参数").Json(c)
		return
	}

	// 调用通用的SSO处理逻辑
	handleSsoFlow(c, appKey, redirectUrl)
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

// handleSsoFlow 通用SSO流程处理
func handleSsoFlow(c *gin.Context, appKey, redirectUrl string) {
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
		log.Info("用户没有登录")

		// 用户未登录，根据请求类型决定响应方式
		acceptHeader := c.GetHeader("Accept")
		contentType := c.GetHeader("Content-Type")

		// 判断是否为API请求（JSON格式）
		isApiRequest := strings.Contains(acceptHeader, "application/json") ||
			strings.Contains(contentType, "application/json") ||
			strings.HasPrefix(c.Request.URL.Path, "/api/")
		log.Info("判断是否为API请求")

		if isApiRequest {
			log.Info("是 api 请求")
			// API请求返回JSON
			c.JSON(200, map[string]interface{}{
				"needLogin":   true,
				"sso":         true,
				"appKey":      appKey,
				"redirectUrl": redirectUrl,
				"message":     "请登录以继续SSO认证",
			})
		} else {
			// 浏览器请求重定向到前端登录页面
			loginUrl := "http://127.0.0.1:5173/login?sso=1&app_key=" + url.QueryEscape(appKey) + "&redirect_url=" + url.QueryEscape(redirectUrl)
			log.Info("重定向地址社区登录地址：" + loginUrl)
			c.Redirect(302, loginUrl)
		}
		return
	}

	// 用户已登录，生成授权码并跳转
	authCode, err := ssoService.GenerateAuthCode(appKey, claims.ID, redirectUrl)
	if err != nil {
		result.Err("生成授权码失败").Json(c)
		return
	}

	// 重定向到第三方应用
	finalUrl := redirectUrl + "?code=" + authCode
	c.Redirect(302, finalUrl)
}
