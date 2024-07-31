package frontend

import (
	"strconv"
	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/utils/page"

	"xhyovo.cn/community/pkg/log"

	"xhyovo.cn/community/cmd/community/middleware"

	services "xhyovo.cn/community/server/service"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
)

var (
	userService services.UserService
)

type editUserForm struct {
	Name      string `binding:"required" form:"name" msg:"用户名不可为空"`
	Desc      string `form:"desc"`
	Subscribe int    `form:"subscribe"`
}

type editPasswordForm struct {
	OldPassword     string `form:"oldPassword" binding:"required" msg:"旧密码不能为空"`
	NewPassword     string `form:"newPassword" binding:"required" msg:"新密码不能为空"`
	ConfirmPassword string `form:"confirmPassword" binding:"required" msg:"确认密码不能为空"`
}

func InitUserRouters(r *gin.Engine) {
	group := r.Group("/community/user")
	group.GET("/info", getUserInfo)
	group.GET("/menu", getUserMenu)
	group.GET("/statistics", statistics)
	group.GET("", listUsers)
	group.GET("/tags/:userId", getTagsByUserId)
	group.GET("/active", activeUsers)
	group.GET("/all", listAllUsers)
	group.GET("/heart", heart)
	group.Use(middleware.OperLogger())
	group.POST("/edit/:tab", updateUser)
}
func activeUsers(ctx *gin.Context) {
	var u services.UserService
	p, limit := page.GetPage(ctx)
	users, count := u.ActiveUsers(p, limit)
	result.Page(users, count, nil).Json(ctx)

}

func getUserMenu(ctx *gin.Context) {
	result.Ok(userService.GetUserMenu(), "ok").Json(ctx)
}

// 获取用户信息
func getUserInfo(ctx *gin.Context) {

	var userService services.UserService
	userId := ctx.Query("userId")
	uId, err := strconv.Atoi(userId)
	if err != nil {
		uId = middleware.GetUserId(ctx)
	}
	user := userService.GetUserSimpleById(uId)

	result.Ok(user, "").Json(ctx)
}

func updateUser(ctx *gin.Context) {

	userId := middleware.GetUserId(ctx)
	t := ctx.Param("tab")
	switch t {
	case "info":
		form := editUserForm{}
		err := ctx.ShouldBind(&form)

		if err != nil {
			log.Warnf("用户id: %d 修改信息参数解析失败,err: %s", userId, err.Error())
			result.Err(utils.GetValidateErr(form, err)).Json(ctx)
			return
		}
		if len(form.Desc) > 200 {
			var msg = "描述长度不可超过200字"
			log.Warnf("用户id: %d 修改信息失败,err: %s", userId, msg)
			result.Err(msg).Json(ctx)
			return
		}
		userService.UpdateUser(&model.Users{Name: form.Name, Desc: form.Desc, ID: userId, Subscribe: form.Subscribe})
	case "pass":
		form := editPasswordForm{}
		err := ctx.ShouldBind(&form)
		if err != nil {
			log.Warnf("用户id: %d 修改密码参数解析失败,err: %s", userId, err.Error())
			result.Err(utils.GetValidateErr(form, err)).Json(ctx)
			return
		}
		// check 旧密码

		if !services.ComparePswd(userService.GetUserById(userId).Password, form.OldPassword) {
			var msg = "旧密码不一致"
			log.Warnf("用户id: %d 修改密码失败,err: %s", userId, msg)
			result.Err(msg).Json(ctx)
			return
		}
		// check 新密码
		if form.NewPassword != form.ConfirmPassword {
			var msg = "两次新密码不一致"
			log.Warnf("用户id: %d 修改密码失败,err: %s", userId, msg)
			result.Err(msg).Json(ctx)
			return
		}
		pwd, err := services.GetPwd(form.ConfirmPassword)
		if err != nil {
			var msg = "加密密码错误"
			log.Warnf("用户id: %d 修改密码失败,err: %s", userId, msg)
			result.Err(msg).Json(ctx)
			return
		}
		userService.UpdateUser(&model.Users{Password: string(pwd), ID: userId})
	case "avatar":
		type avatar struct {
			Avatar string `json:"avatar" binding:"required" msg:"头像不能为空"`
		}
		object := &avatar{}
		if err := ctx.ShouldBindJSON(&object); err != nil {
			log.Warnf("用户id: %d 修改头像参数解析失败,err: %s", userId, err.Error())
			result.Err(utils.GetValidateErr(object, err)).Json(ctx)
			return
		}
		// 更改用户信息
		userService.UpdateUser(&model.Users{ID: userId, Avatar: object.Avatar})
	}
	result.OkWithMsg(nil, "修改成功").Json(ctx)
}

// 数据统计
func statistics(ctx *gin.Context) {
	types, err := strconv.Atoi(ctx.Query("type"))
	if err != nil {
		types = 1
	}
	userId := middleware.GetUserId(ctx)
	m := userService.Statistics(userId, types)
	result.Ok(m, "").Json(ctx)
}

func listUsers(ctx *gin.Context) {
	name := ctx.Query("name")
	users := userService.ListUsers(name)
	result.Ok(users, "").Json(ctx)
}

var userTagS services.UserTag

func getTagsByUserId(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		log.Warnf("获取用户标签参数解析失败,err: %s", err.Error())
		result.Err(err.Error()).Json(ctx)
		return
	}
	tagNames := userTagS.GetTagsByUserId(userId)
	result.Ok(tagNames, "").Json(ctx)
}

// 查询所有用户
func listAllUsers(ctx *gin.Context) {
	var data []model.Users
	model.User().Select("id", "name").Find(&data)

	users := []map[string]interface{}{}

	for _, item := range data {
		users = append(users, map[string]interface{}{
			"name": item.Name,
			"code": item.ID,
		})
	}
	result.Ok(users, "").Json(ctx)
}

// 心跳
func heart(ctx *gin.Context) {
	// 获取用户 id，ip
	userId := middleware.GetUserId(ctx)
	// 获取ip
	ip := utils.GetClientIP(ctx)
	ip = ctx.ClientIP()
	cache := cache.GetInstance()
	token := ctx.GetHeader(middleware.AUTHORIZATION)
	if len(token) == 0 {
		token, _ = ctx.Cookie(middleware.AUTHORIZATION)
	}
	key := constant.HEARTBEAT + strconv.Itoa(userId)
	valueIp, b := cache.Get(key)
	// ip 不一致则说明同一时间内多个用户使用一个账号

	if b {
		if valueIp != ip {
			// 将 token 设置为无效并且清空用户登陆状态并且加入黑名单中
			// 1.前端清空 store
			// 2.后端将 token 存入黑名单
			// 3.对黑名单进行计数：可能用户不服
			var blackService = services.BlacklistService{}
			blackService.Add(userId, token)
			result.Err("你已涉嫌同一账号多人使用，请注意你的行为").Json(ctx)
			return
		}
	}
	cache.Set(key, ip, constant.HEARTBEAT_TTL)
	result.Ok(nil, "").Json(ctx)
	return
}

// todo 后台 黑名单界面
// todo 前端全局心跳
