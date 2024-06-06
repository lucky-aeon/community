package frontend

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"xhyovo.cn/community/cmd/community/middleware"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/pkg/result"
	"xhyovo.cn/community/pkg/utils"
	page2 "xhyovo.cn/community/pkg/utils/page"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
	"xhyovo.cn/community/server/service/event"
)

func InitSubscriptionRouters(r *gin.Engine) {
	group := r.Group("/community")
	group.GET("/subscription", listSubscription)
	group.GET("/event", eventList)
	group.POST("/subscription/state", subscriptionState)
	group.Use(middleware.OperLogger())
	group.POST("/subscribe", subscribe)
}

// 查看订阅列表
func listSubscription(ctx *gin.Context) {
	var su services.SubscriptionService
	userId := middleware.GetUserId(ctx)
	eventId, _ := strconv.Atoi(ctx.DefaultQuery("eventId", "0"))
	page, limit := page2.GetPage(ctx)
	subscriptions, count := su.ListSubscription(userId, eventId, page, limit)
	result.Page(subscriptions, count, nil).Json(ctx)
}

// 查看对应事件订阅状态
func subscriptionState(ctx *gin.Context) {
	var su services.SubscriptionService
	var subscription model.SubscriptionState
	if err := ctx.ShouldBindJSON(&subscription); err != nil {
		log.Warnf("用户id: %d 查看事件订阅状态参数解析失败,err: %s", middleware.GetUserId(ctx), err.Error())
		result.Err(utils.GetValidateErr(subscription, err)).Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	subscription.SubscriberId = userId
	result.Ok(su.SubscriptionState(&subscription), "").Json(ctx)
}

// 订阅/取消订阅
func subscribe(ctx *gin.Context) {

	var subscription model.Subscriptions
	if err := ctx.ShouldBindJSON(&subscription); err != nil {
		msg := utils.GetValidateErr(comment, err)
		log.Warnf("用户id: %d 订阅事件参数解析失败,err: %s", middleware.GetUserId(ctx), msg)
		result.Err(msg).Json(ctx)
		return
	}
	userId := middleware.GetUserId(ctx)
	subscription.SubscriberId = userId
	businessId := subscription.BusinessId

	// 如果是订阅文章则找文章信息
	var a services.ArticleService
	var flag bool = false
	if subscription.EventId == event.CommentUpdateEvent {
		if a.GetById(businessId).UserId == userId {
			flag = true
		}
	}

	if subscription.EventId == event.UserFollowingEvent && businessId == userId {
		result.Err("关注用户不能是自己").Json(ctx)
		return
	} else if flag {
		result.Err("订阅文章不能是自己所发布").Json(ctx)
		return
	}

	var su services.SubscriptionService
	msg := "取消订阅"
	flag = false
	if su.Subscribe(&subscription) {
		msg = "订阅成功"
		flag = true
	}
	result.OkWithMsg(flag, msg).Json(ctx)
}

func eventList(ctx *gin.Context) {
	result.Ok(event.List(), "").Json(ctx)
}
