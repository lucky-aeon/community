package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

// GitHubActivateRequest GitHub激活请求结构
type GitHubActivateRequest struct {
	InviteCode string `json:"invite_code" binding:"required" msg:"邀请码不能为空"`
	GitHubUser struct {
		ID        int64  `json:"id" binding:"required"`
		Username  string `json:"username" binding:"required"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Name      string `json:"name"`
	} `json:"github_user" binding:"required"`
}

// InitGitHubAuthRouters 初始化GitHub认证路由
func InitGitHubAuthRouters(ctx *gin.Engine) {
	group := ctx.Group("/community/auth/github")

	// GitHub OAuth 登录入口
	group.GET("/login", GitHubLogin)

	// GitHub OAuth 回调处理
	group.GET("/callback", GitHubCallback)

	// GitHub 账号激活（需要邀请码）
	group.POST("/activate", GitHubActivate)

	// 需要认证的路由
	authGroup := group.Use(middleware.Auth)

	// 获取当前用户的GitHub绑定信息
	authGroup.GET("/binding", GetGitHubBinding)

	// 解绑GitHub账号
	authGroup.DELETE("/unbind", UnbindGitHub)
}

// GitHubLogin GitHub登录入口
func GitHubLogin(c *gin.Context) {
	var githubAuthService services.GitHubAuthService

	// 生成GitHub OAuth授权URL
	authURL := githubAuthService.GetGitHubOAuthURL()

	log.Info("GitHub OAuth 授权URL: " + authURL)

	// 重定向到GitHub授权页面
	c.Redirect(http.StatusFound, authURL)
}

// GitHubCallback GitHub OAuth回调处理
func GitHubCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	errorParam := c.Query("error")

	// 检查是否有错误
	if errorParam != "" {
		log.Error("GitHub OAuth 授权失败: " + errorParam)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "GitHub授权失败: " + errorParam,
		})
		return
	}

	// 检查必要参数
	if code == "" || state == "" {
		log.Error("GitHub OAuth 回调缺少必要参数")
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   true,
			"message": "授权回调参数不完整",
		})
		return
	}

	var githubAuthService services.GitHubAuthService

	// 处理OAuth回调，获取GitHub用户信息
	githubUser, err := githubAuthService.HandleOAuthCallback(code, state)
	if err != nil {
		log.Error("处理GitHub OAuth回调失败: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   true,
			"message": "GitHub授权处理失败: " + err.Error(),
		})
		return
	}

	// 检查该GitHub账号是否已绑定用户
	existingUser, err := githubAuthService.FindUserByGitHubID(githubUser.ID)
	if err == nil && existingUser != nil {
		// 检查用户是否被禁用
		var userService services.UserService
		if userService.IsBlack(existingUser.ID) {
			log.Warn("被禁用用户尝试通过GitHub登录: " + existingUser.Name)
			c.JSON(http.StatusForbidden, gin.H{
				"error":   true,
				"message": "账号已被禁用，无法登录",
			})
			return
		}

		// 已绑定用户，直接登录
		token, err := middleware.GenerateToken(existingUser.ID, existingUser.Name)
		if err != nil {
			log.Error("生成JWT token失败: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   true,
				"message": "登录失败",
			})
			return
		}

		log.Info("GitHub用户 " + githubUser.Login + " 登录成功")

		// 返回登录成功结果
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "登录成功",
			"token":   token,
			"user": gin.H{
				"id":     existingUser.ID,
				"name":   existingUser.Name,
				"avatar": existingUser.Avatar,
			},
		})
		return
	}

	// 新用户，需要邀请码激活
	log.Info("新GitHub用户 " + githubUser.Login + " 需要邀请码激活")

	// 返回需要激活的响应
	c.JSON(http.StatusOK, gin.H{
		"need_activation": true,
		"message":         "请输入邀请码激活账号",
		"github_user": gin.H{
			"id":         githubUser.ID,
			"username":   githubUser.Login,
			"email":      githubUser.Email,
			"avatar_url": githubUser.AvatarURL,
			"name":       githubUser.Name,
		},
	})
}

// GitHubActivate GitHub账号激活
func GitHubActivate(c *gin.Context) {
	var req GitHubActivateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.Err("参数错误: " + err.Error()).Json(c)
		return
	}

	var userService services.UserService
	var githubAuthService services.GitHubAuthService
	var codeService services.CodeService

	// 验证邀请码是否存在
	if !codeService.Exist(req.InviteCode) {
		result.Err("邀请码不存在或已失效").Json(c)
		return
	}

	// 检查该GitHub账号是否已经绑定其他用户
	existingGitHubUser, err := githubAuthService.FindUserByGitHubID(req.GitHubUser.ID)
	if err == nil && existingGitHubUser != nil {
		result.Err("该GitHub账号已绑定其他用户").Json(c)
		return
	}

	// 构建GitHub用户信息
	githubUser := &model.GitHubUser{
		ID:        req.GitHubUser.ID,
		Login:     req.GitHubUser.Username,
		Email:     req.GitHubUser.Email,
		AvatarURL: req.GitHubUser.AvatarURL,
		Name:      req.GitHubUser.Name,
	}

	// 检查邀请码是否已被使用
	if codeService.IsCodeUsed(req.InviteCode) {
		// 场景A：旧用户绑定GitHub账号
		existingUser, err := codeService.FindUserByInviteCode(req.InviteCode)
		if err != nil {
			log.Error("查找邀请码对应用户失败: " + err.Error())
			result.Err("邀请码对应的用户不存在").Json(c)
			return
		}

		// 检查用户是否被禁用
		if userService.IsBlack(existingUser.ID) {
			result.Err("账号已被禁用，无法绑定").Json(c)
			return
		}

		// 检查用户是否已经绑定了其他GitHub账号
		_, err = githubAuthService.GetGitHubBinding(existingUser.ID)
		if err == nil {
			result.Err("该账号已绑定其他GitHub账号").Json(c)
			return
		}

		// 绑定GitHub账号到现有用户
		err = githubAuthService.BindGitHubAccount(existingUser.ID, githubUser)
		if err != nil {
			log.Error("绑定GitHub账号失败: " + err.Error())
			result.Err("绑定失败: " + err.Error()).Json(c)
			return
		}

		// 生成JWT token
		token, err := middleware.GenerateToken(existingUser.ID, existingUser.Name)
		if err != nil {
			log.Error("生成JWT token失败: " + err.Error())
			result.Err("登录失败").Json(c)
			return
		}

		log.Info("用户 " + existingUser.Name + " 成功绑定GitHub账号")

		// 返回绑定成功结果
		result.Ok(gin.H{
			"message": "GitHub账号绑定成功",
			"token":   token,
			"user": gin.H{
				"id":     existingUser.ID,
				"name":   existingUser.Name,
				"avatar": existingUser.Avatar,
			},
		}, "").Json(c)

	} else {
		// 场景B：新用户GitHub注册
		newUser, err := userService.CreateUserFromGitHub(githubUser, req.InviteCode)
		if err != nil {
			log.Error("创建GitHub用户失败: " + err.Error())
			result.Err("账号创建失败: " + err.Error()).Json(c)
			return
		}

		// 绑定GitHub账号
		err = githubAuthService.BindGitHubAccount(newUser.ID, githubUser)
		if err != nil {
			log.Error("绑定GitHub账号失败: " + err.Error())
			result.Err("绑定失败: " + err.Error()).Json(c)
			return
		}

		// 生成JWT token
		token, err := middleware.GenerateToken(newUser.ID, newUser.Name)
		if err != nil {
			log.Error("生成JWT token失败: " + err.Error())
			result.Err("登录失败").Json(c)
			return
		}

		log.Info("GitHub新用户注册成功: " + newUser.Name)

		// 返回注册成功结果
		result.Ok(gin.H{
			"message": "账号注册成功",
			"token":   token,
			"user": gin.H{
				"id":     newUser.ID,
				"name":   newUser.Name,
				"avatar": newUser.Avatar,
			},
		}, "").Json(c)
	}
}

// GetGitHubBinding 获取当前用户的GitHub绑定信息
func GetGitHubBinding(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		result.Err("用户未登录").Json(c)
		return
	}

	userClaims := claims.(*middleware.JwtCustomClaims)
	var userService services.UserService

	binding, err := userService.GetGitHubBinding(userClaims.ID)
	if err != nil {
		result.Err("未绑定GitHub账号").Json(c)
		return
	}

	// 返回脱敏的绑定信息
	result.Ok(gin.H{
		"github_id":       binding.GitHubID,
		"github_username": binding.GitHubUsername,
		"github_email":    binding.GitHubEmail,
		"github_avatar":   binding.GitHubAvatar,
		"bound_at":        binding.BoundAt,
	}, "").Json(c)
}

// UnbindGitHub 解绑GitHub账号
func UnbindGitHub(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		result.Err("用户未登录").Json(c)
		return
	}

	userClaims := claims.(*middleware.JwtCustomClaims)
	var githubAuthService services.GitHubAuthService

	err := githubAuthService.UnbindGitHubAccount(userClaims.ID)
	if err != nil {
		log.Error("解绑GitHub账号失败: " + err.Error())
		result.Err("解绑失败").Json(c)
		return
	}

	log.Info("用户 " + userClaims.Name + " 解绑GitHub账号成功")
	result.Ok(nil, "解绑成功").Json(c)
}
