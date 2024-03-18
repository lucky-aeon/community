package services

import (
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
func (s *SubscriptionService) Do(eventId int, b BusinessId) {

	go func(eventId int, b BusinessId) {

		subscriptions := s.ListSubscriptionUserId(eventId, b.CurrentBusinessId)
		var m MessageService

		if len(subscriptions) > 0 {
			var userIds []int
			for i := range subscriptions {
				userIds = append(userIds, subscriptions[i].SubscriberId)
			}
			sendId := subscriptions[0].SendId
			var users []model.Users
			model.User().Where("id in ?", userIds).Select("account", "id", "subscribe").Find(&users)
			var emails []string
			for i := range users {
				if users[i].Subscribe {
					emails = append(emails, users[i].Account)
				}
			}
			messageTemplate := messageDao.GetMessageTemplate(eventId)
			msg := m.GetMsg(messageTemplate, b)
			m.SendMessages(sendId, constant.NOTICE, b.ArticleId, userIds, msg) // todo确定内容
			email.Send(emails, msg, "技术鸭社区")
		}

	}(eventId, b)
}

// 触发 @ 事件
func (s *SubscriptionService) ConstantAtSend(eventId, triggerId int, content string, b BusinessId) {
	go func(eventId, triggerId int, content string, b BusinessId) {
		var m MessageService
		if strings.Contains(content, "@") {
			mentionedUsers := extractMentionedUsers(content)
			if len(mentionedUsers) == 0 {
				return
			}

			var users []model.Users
			model.User().Where("name in ?", mentionedUsers).Select("account", "id", "subscribe").Find(&users)

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
			m.SendMessages(triggerId, constant.MENTION, b.UserId, ids, msg)
			email.Send(emails, msg, "技术鸭社区")
		}
	}(eventId, triggerId, content, b)

}

func extractMentionedUsers(s string) []string {
	var mentionedUsers []string

	// 搜索字符串中的 "@" 符号
	startIndex := 0
	for {
		atIndex := strings.Index(s[startIndex:], "@")
		if atIndex == -1 {
			break
		}

		atIndex += startIndex
		startIndex = atIndex + 1

		// 寻找下一个空格或句号，作为结束索引
		endIndex := len(s)
		spaceIndex := strings.Index(s[startIndex:], " ")
		dotIndex := strings.Index(s[startIndex:], ".")

		if spaceIndex != -1 && spaceIndex < endIndex {
			endIndex = startIndex + spaceIndex
		}
		if dotIndex != -1 && dotIndex < endIndex {
			endIndex = startIndex + dotIndex
		}

		mentionedUser := s[startIndex:endIndex]
		mentionedUsers = append(mentionedUsers, mentionedUser)
	}

	return mentionedUsers
}
