package services

import (
	mapset "github.com/deckarep/golang-set/v2"
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
	nameMap := userService.ListByIdsToMap(userIds)

	for i := range subscriptions {
		v := &subscriptions[i]
		businessId := v.BusinessId
		if v.EventId == event.CommentUpdateEvent {
			v.BusinessName = articleMap[businessId]
		} else if v.EventId == event.UserFollowingEvent {
			v.BusinessName = nameMap[businessId].Name
		}
		v.EventName = event.GetMsg(v.EventId)
	}

	return subscriptions
}

// 查看对应事件订阅状态
func (*SubscriptionService) SubscriptionState(subscriptions *model.SubscriptionState) bool {

	return subscriptionDao.SubscriptionState(subscriptions)
}

// 订阅/取消订阅
func (*SubscriptionService) Subscribe(subscription *model.Subscriptions) bool {
	businessId := subscription.BusinessId
	if event.CommentUpdateEvent == subscription.EventId {
		model.Article().Where("id = ?", businessId).Select("user_id").First(&subscription.SendId)
	} else if event.UserFollowingEvent == subscription.EventId {
		model.User().Where("id = ?", businessId).Select("id").First(&subscription.SendId)
	}

	return subscriptionDao.Subscribe(subscription)
}

// 查询事件订阅的userid
func (*SubscriptionService) ListSubscriptionUserId(event, businessId int) []model.Subscriptions {
	return subscriptionDao.ListSubscriptions(event, businessId)
}

// 触发订阅事件
func (s *SubscriptionService) Do(eventId int, b SubscribeData) {

	go func(eventId int, b SubscribeData) {

		subscriptions := s.ListSubscriptionUserId(eventId, b.SubscribeId)
		if len(subscriptions) > 0 {
			var userIds []int
			for i := range subscriptions {
				userIds = append(userIds, subscriptions[i].SubscriberId)
			}
			sendId := subscriptions[0].SendId
			var users []model.Users
			model.User().Where("id in ?", userIds).Select("account", "id", "subscribe").Find(&users)
			send(users, eventId, constant.NOTICE, sendId, b)
		}

	}(eventId, b)
}

// 触发 @ 事件
func (s *SubscriptionService) ConstantAtSend(eventId, triggerId int, content string, b SubscribeData) {
	go func(eventId, triggerId int, content string, b SubscribeData) {

		userNames := findAtUser(content)
		if len(userNames) > 0 {
			var users []model.Users
			model.User().Where("name in ?", userNames).Select("account", "id", "subscribe").Find(&users)
			send(users, eventId, constant.MENTION, triggerId, b)
		}
	}(eventId, triggerId, content, b)

}

func (s *SubscriptionService) Send(eventId, eventType, fromId, toId int, b SubscribeData) {
	go func(eventId, fromId, toId int, b SubscribeData) {
		var users []model.Users
		users = append(users, model.Users{ID: toId})
		send(users, eventId, eventType, fromId, b)
	}(eventId, fromId, toId, b)
}

func findAtUser(content string) []string {

	// 使用正则表达式匹配文本中的 @ 符号
	re := regexp.MustCompile(`\s@(\w+)\s`)
	matches := re.FindAllString(content, -1)
	names := mapset.NewSet[string]()
	for i := range matches {
		names.Add(strings.ReplaceAll(strings.TrimSpace(matches[i]), "@", ""))
	}
	return names.ToSlice()
}

func send(users []model.Users, eventId, eventType, sendId int, b SubscribeData) {
	var m MessageService

	var ids []int
	var emails []string
	for i := range users {
		ids = append(ids, users[i].ID)
		if users[i].Subscribe {
			emails = append(emails, users[i].Account)
		}
	}
	messageTemplate := messageDao.GetMessageTemplate(eventId)
	msg := m.GetMsg(messageTemplate, b)
	m.SendMessages(sendId, eventType, b.CurrentBusinessId, ids, msg)
	email.Send(emails, msg, "技术鸭社区")
}
