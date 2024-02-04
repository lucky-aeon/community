package services

import (
	"xhyovo.cn/community/pkg/email"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type SubscriptionService struct {
}

// 查看订阅列表
func (*SubscriptionService) ListSubscription(userId, eventId, page, limit int) []model.Subscriptions {
	subscriptions := subscriptionDao.ListSubscription(userId, eventId, page, limit)

	var userIds []int
	var articleIds []int

	for i := range subscriptions {
		v := subscriptions[i]
		if v.EventId == event.CommentUpdateEvent {
			articleIds = append(articleIds, v.BusinessId)
		} else if v.EventId == event.UserFollowingEvent {
			userIds = append(userIds, v.BusinessId)
		}
	}
	var articleService ArticleService
	articleMap := articleService.ListByIdsSelectIdTitleMap(articleIds)

	var userService UserService
	nameMap := userService.ListByIdsSelectIdNameMap(userIds)

	for i := range subscriptions {
		v := &subscriptions[i]
		businessId := v.BusinessId
		if v.EventId == event.CommentUpdateEvent {
			v.BusinessName = articleMap[businessId]
		} else if v.EventId == event.UserFollowingEvent {
			v.BusinessName = nameMap[businessId]
		}
		v.EventName = event.GetMsg(v.EventId)
	}

	return subscriptions
}

// 查看对应事件订阅状态
func (*SubscriptionService) SubscriptionState(subscriptions *model.Subscriptions) bool {

	return subscriptionDao.SubscriptionState(subscriptions)
}

// 订阅/取消订阅
func (*SubscriptionService) Subscribe(subscription *model.Subscriptions) bool {
	return subscriptionDao.Subscribe(subscription)
}

// 查询事件订阅的userid
func (*SubscriptionService) ListSubscriptionUserId(event, businessId int) []int {
	return subscriptionDao.ListSubscriptionUserId(event, businessId)
}

// 触发订阅事件
func (s *SubscriptionService) Do(subscription *model.Subscriptions) {

	go func(subscription *model.Subscriptions) {
		userIds := s.ListSubscriptionUserId(subscription.EventId, subscription.BusinessId)
		if len(userIds) > 0 {
			to := userDao.ListByIds(userIds...)

			messageTemplate := messageDao.GetMessageTemplate(subscription.EventId)

			email.Send(to, messageTemplate, "技术鸭社区")
		}
	}(subscription)

}
