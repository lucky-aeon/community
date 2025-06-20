package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	xt "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

type registerForm struct {
	Code     string `binding:"required" form:"code" msg:"code不能为空" `
	Account  string `binding:"required,email" form:"account" msg:"邮箱格式不正确"`
	Name     string `binding:"required" form:"name" msg:"用户名不能为空"`
	Password string `binding:"required" form:"password" msg:"密码不能为空"`
}

func InitLoginRegisterRouters(ctx *gin.Engine) {
	group := ctx.Group("/community")
	group.POST("/login", Login)
	group.POST("/register", Register)
	
	// 添加登录页面路由（支持SSO参数）
	ctx.GET("/login", LoginPage)
}

func Login(c *gin.Context) {

	var login model.LoginForm
	if err := c.ShouldBindJSON(&login); err != nil {
		result.Err(utils.GetValidateErr(login, err)).Json(c)
		return
	}
	loginLog := model.LoginLogs{
		Account:   login.Account,
		Browser:   c.Request.UserAgent(),
		Equipment: c.GetHeader("Sec-Ch-Ua-Platform"),
		Ip:        utils.GetClientIP(c),
		CreatedAt: xt.Now(),
	}
	var logS services.LogServices
	user, err := services.Login(login)
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(err.Error()).Json(c)
		return
	}

	// 判断黑名单
	var userService services.UserService
	if userService.IsBlack(user.ID) {
		result.Err("你已涉嫌违规社区文化，已被纳入小黑屋，如误封请联系我：xhyQAQ250").Json(c)
		return
	}

	token, err := middleware.GenerateToken(user.ID, user.Name+uuid.New().String())
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(err.Error()).Json(c)
		return
	}

	c.SetCookie(middleware.AUTHORIZATION, token, int(constant.Token_TTl.Seconds()), "/", c.Request.Host, false, true)
	loginLog.State = "登录成功"
	logS.InsertLoginLog(loginLog)
	result.OkWithMsg(map[string]string{"token": token}, "登录成功").Json(c)
}

func Register(c *gin.Context) {
	var form registerForm

	err := c.ShouldBindJSON(&form)
	loginLog := model.LoginLogs{
		Account:   form.Account,
		Browser:   c.Request.UserAgent(),
		Equipment: c.GetHeader("Sec-Ch-Ua-Platform"),
		Ip:        utils.GetClientIP(c),
		CreatedAt: xt.Now(),
	}
	var logS services.LogServices
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(utils.GetValidateErr(form, err)).Json(c)
		return
	}

	if err != nil {
		log.Warnf("账户: %s 注册失败,获取加密密码错误,err %s", form.Account, err.Error())
		result.Err(err.Error()).Json(c)
		return
	}

	id, err := services.Register(form.Account, form.Password, form.Name, form.Code)
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(err.Error()).Json(c)
		return
	}
	var d services.Draft
	d.InitDraft(id)

	loginLog.State = "注册成功"
	logS.InsertLoginLog(loginLog)
	token, err := middleware.GenerateToken(id, form.Name)
	if err != nil {
		loginLog.State = err.Error()
		logS.InsertLoginLog(loginLog)
		result.Err(err.Error()).Json(c)
		return
	}
	c.SetCookie(middleware.AUTHORIZATION, token, int(constant.Token_TTl.Seconds()), "/", c.Request.Host, false, true)

	result.OkWithMsg(map[string]string{"token": token}, "注册成功").Json(c)
}

// LoginPage 登录页面（支持SSO参数）
func LoginPage(c *gin.Context) {
	sso := c.Query("sso")
	appKey := c.Query("app_key")
	redirectUrl := c.Query("redirect_url")
	
	// 检查用户是否已经登录
	token := c.GetHeader("Authorization")
	if len(token) == 0 {
		token, _ = c.Cookie("Authorization")
	}
	
	claims, err := middleware.ParseToken(token)
	if err == nil && claims.ID > 0 {
		// 用户已登录
		if sso == "1" && appKey != "" && redirectUrl != "" {
			// SSO场景：用户已登录，直接生成授权码并跳转
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
			
			// 生成授权码
			authCode, err := ssoService.GenerateAuthCode(appKey, claims.ID, redirectUrl)
			if err != nil {
				result.Err("生成授权码失败").Json(c)
				return
			}
			
			// 重定向到第三方应用
			finalUrl := redirectUrl + "?code=" + authCode
			c.Redirect(302, finalUrl)
			return
		}
		
		// 普通登录场景：用户已登录，重定向到首页
		c.Redirect(302, "/")
		return
	}
	
	// 用户未登录，显示登录页面
	if sso == "1" {
		// SSO登录页面，可以传递SSO参数到前端
		c.JSON(200, map[string]interface{}{
			"needLogin": true,
			"sso":       true,
			"appKey":    appKey,
			"redirectUrl": redirectUrl,
			"message":   "请登录以继续SSO认证",
		})
	} else {
		// 普通登录页面
		c.JSON(200, map[string]interface{}{
			"needLogin": true,
			"sso":       false,
			"message":   "请登录",
		})
	}
}
