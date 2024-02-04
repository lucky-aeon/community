package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
	"xhyovo.cn/community/server/service/event"
)

func InitSubscriptionRouters(r *gin.Engine) {
	group := r.Group("/community")
	group.GET("/subscription", listSubscription)
	group.GET("/event", eventList)
	group.POST("/subscription/state", subscriptionState)
	group.POST("/subscribe", subscribe)
}

// 查看订阅列表
func listSubscription(ctx *gin.Context) {
	var su services.SubscriptionService
	userId := middleware.GetUserId(ctx)
	eventId, _ := strconv.Atoi(ctx.DefaultQuery("eventId", "0"))
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "15"))
	result.Ok(su.ListSubscription(userId, eventId, page, limit), "").Json(ctx)
}

// 查看对应事件订阅状态
func subscriptionState(ctx *gin.Context) {
	var su services.SubscriptionService
	var subscription model.Subscriptions
	if err := ctx.ShouldBindJSON(&subscription); err != nil {
		result.Err(utils.GetValidateErr(subscription, err)).Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	subscription.UserId = userId
	result.Ok(su.SubscriptionState(&subscription), "").Json(ctx)
}

// 订阅/取消订阅
func subscribe(ctx *gin.Context) {
	var su services.SubscriptionService
	var subscription model.Subscriptions
	if err := ctx.ShouldBindJSON(&subscription); err != nil {
		result.Err(utils.GetValidateErr(subscription, err)).Json(ctx)
	}
	userId := middleware.GetUserId(ctx)
	subscription.UserId = userId
	msg := ""
	if su.Subscribe(&subscription) {
		msg = "订阅成功"
	} else {
		msg = "取消订阅"
	}
	result.Ok(msg, "").Json(ctx)
}

func eventList(ctx *gin.Context) {
	result.Ok(event.List(), "").Json(ctx)
}
