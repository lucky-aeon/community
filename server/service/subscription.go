package services

import (
	"regexp"
	"strings"
	"xhyovo.cn/community/pkg/constant"
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
func (*SubscriptionService) ListSubscriptionUserId(event, businessId int) []model.Subscriptions {
	return subscriptionDao.ListSubscriptions(event, businessId)
}

// 触发订阅事件
func (s *SubscriptionService) Do(eventId, businessId, triggerId int, content string) {

	go func(eventId, businessId int, content string) {

		subscriptions := s.ListSubscriptionUserId(eventId, businessId)
		var m MessageService

		if len(subscriptions) > 0 {
			var userIds []int
			for i := range subscriptions {
				userIds = append(userIds, subscriptions[i].SubscriberId)
			}
			sendId := subscriptions[0].SendId
			to := userDao.ListByIds(userIds...)

			messageTemplate := messageDao.GetMessageTemplate(eventId)
			msg := m.GetMsg(messageTemplate, eventId, businessId)
			m.SendMessages(sendId, constant.NOTICE, userIds, content)
			email.Send(to, msg, "技术鸭社区")

			// @的情况再给@发送
			if strings.Contains(content, "@") {
				re := regexp.MustCompile(`@(\w+)`)
				matches := re.FindAllStringSubmatch(content, -1)

				var usernames []string
				for _, match := range matches {
					usernames = append(usernames, match[1])
				}
				var userS UserService
				emails, ids := userS.ListByNameSelectEmailAndId(usernames)
				m.SendMessages(triggerId, constant.MENTION, ids, content)
				email.Send(emails, content, "技术鸭社区")
			}
		}

	}(eventId, businessId, content)

}
