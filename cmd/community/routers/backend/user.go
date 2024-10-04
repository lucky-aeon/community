package backend

import (
	"fmt"
	"strconv"
	"xhyovo.cn/community/pkg/cache"
	"xhyovo.cn/community/pkg/constant"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

func InitUserRouters(r *gin.Engine) {

	group := r.Group("/community/admin/user")
	group.GET("", listUser)
	group.GET("/black", blackListUser)
	group.Use(middleware.OperLogger())
	group.POST("", updateUser)
	group.POST("/renewal", renewal)
	group.DELETE("/:id", deleteUser)
	group.PUT("/reset/pwd", setRangePassword)

	group.DELETE("/black/ban", banUser)
	group.POST("/black/unBan", unBanUser)
}

func listUser(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)

	conditionUser := model.UserSimple{
		Account: fmt.Sprintf("%%%s%%", ctx.Query("account")),
	}
	userState, err := strconv.Atoi(ctx.Query("userState"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var u services.UserService

	users, count := u.PageUsers(p, limit, conditionUser, userState)
	result.Page(users, count, nil).Json(ctx)
}

type updateUserInfo struct {
	model.Users
	Tags []int `json:"tags"`
}

func updateUser(ctx *gin.Context) {

	user := updateUserInfo{}

	if err := ctx.ShouldBindBodyWith(&user, binding.JSON); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	var u services.UserService
	u.UpdateUser(&user.Users)
	user.Tags = userTagS.AssignUserLabel(user.ID, user.Tags)

	result.OkWithMsg(user, "修改成功").Json(ctx)
}

// 重置用户密码，随机生成
func setRangePassword(ctx *gin.Context) {
	account := ctx.Query("account")
	if account == "" {
		result.Err("用户名不能为空").Json(ctx)
		return
	}
	var u services.UserService
	if u.ResetPwd(account) {
		result.OkWithMsg(nil, "重置成功, 密码已发送至邮箱").Json(ctx)
		return
	}
	result.OkWithMsg(nil, "重置失败，用户不存在").Json(ctx)
}

func deleteUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	if userId == id {
		result.Err("不能删除自己").Json(ctx)
		return
	}
	var u services.UserService
	u.DeleteUser(id)
	log.Infof("用户id: %d,删除用户: %d", userId, id)
	result.OkWithMsg(nil, "删除成功").Json(ctx)
}

// 查询黑名单的用户
func blackListUser(ctx *gin.Context) {
	p, limit := page.GetPage(ctx)
	var u services.UserService
	user, count := u.ListBlackUser(p, limit)
	result.Page(user, count, nil).Json(ctx)
}

// 禁掉用户
func banUser(ctx *gin.Context) {
	account := ctx.Query("account")

	var u services.UserService
	if !u.ExistUserByAccount(account) {
		result.Err("用户不存在").Json(ctx)
		return
	}
	u.BanByUserAccount(account)
	result.OkWithMsg(nil, "已 ban 掉用户："+account).Json(ctx)
}

// 解封用户
func unBanUser(ctx *gin.Context) {
	id := ctx.Query("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	var u services.UserService
	u.UnBanByUserId(idInt)
	cache := cache.GetInstance()
	cache.Delete(constant.HEARTBEAT + id)
	cache.Delete(constant.BLACK_LIST + id)
	cache.Delete(constant.BLACK_LIST_COUNT + id)

	result.OkWithMsg(nil, "已解封用户："+id).Json(ctx)
}

// 续费
func renewal(ctx *gin.Context) {
	renewalValue := ctx.Query("renewalValue")
	userId := ctx.Query("userId")

	year, err := strconv.Atoi(renewalValue)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}

	var u services.UserService
	if err = u.Renewal(userIdInt, year); err != nil {
		result.Err(err.Error()).Json(ctx)
		return
	}
	result.OkWithMsg(nil, "续费成功").Json(ctx)

}
