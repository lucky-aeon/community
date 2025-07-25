package services

import (
	mapset "github.com/deckarep/golang-set/v2"
	"regexp"
	"strconv"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/email"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type SubscriptionService struct {
}

// 查看订阅列表
func (*SubscriptionService) ListSubscription(userId, eventId, page, limit int) ([]model.Subscriptions, int64) {
	subscriptions, count := subscriptionDao.ListSubscription(userId, eventId, page, limit)

	var userIds []int
	var articleIds []int
	var courseIds []int
	for i := range subscriptions {
		v := subscriptions[i]
		if v.EventId == event.CommentUpdateEvent {
			articleIds = append(articleIds, v.BusinessId)
		} else if v.EventId == event.UserFollowingEvent {
			userIds = append(userIds, v.BusinessId)
		} else if v.EventId == event.CourseUpdate {
			courseIds = append(courseIds, v.BusinessId)
		}
	}
	var articleService ArticleService
	articleMap := articleService.ListByIdsSelectIdTitleMap(articleIds)

	var userService UserService
	nameMap := userService.ListByIdsToMap(userIds)

	var courseService CourseService
	courseMap := courseService.ListByIdsSelectIdTitleMap(courseIds)
	for i := range subscriptions {
		v := &subscriptions[i]
		businessId := v.BusinessId
		if v.EventId == event.CommentUpdateEvent {
			v.BusinessName = articleMap[businessId]
		} else if v.EventId == event.UserFollowingEvent {
			v.BusinessName = nameMap[businessId].Name
		} else if v.EventId == event.CourseUpdate {
			v.BusinessName = courseMap[businessId]
		}
		v.EventName = event.GetMsg(v.EventId)
	}

	return subscriptions, count
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
	} else if event.CourseUpdate == subscription.EventId {
		model.Course().Where("id = ?", businessId).Select("user_id").First(&subscription.SendId)

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
			send(userIds, eventId, constant.NOTICE, sendId, b, "")
		}

	}(eventId, b)
}

func (s *SubscriptionService) DoWithMessageTempl(eventId int, b SubscribeData, messageTempl string) {
	go func(eventId int, b SubscribeData) {
		subscriptions := s.ListSubscriptionUserId(eventId, b.SubscribeId)
		if len(subscriptions) > 0 {
			var userIds []int
			for i := range subscriptions {
				userIds = append(userIds, subscriptions[i].SubscriberId)
			}
			sendId := subscriptions[0].SendId
			send(userIds, eventId, constant.NOTICE, sendId, b, messageTempl)
		}

	}(eventId, b)
}

// 触发 @ 事件
func (s *SubscriptionService) ConstantAtSend(eventId, triggerId int, content string, b SubscribeData) {
	go func(eventId, triggerId int, content string, b SubscribeData) {

		ids := findAtUser(content)
		send(ids, eventId, constant.MENTION, triggerId, b, "")
	}(eventId, triggerId, content, b)

}

// 触发 @ 事件，直接通知用户
func (s *SubscriptionService) NoticeUsers(eventId, triggerId int, userIds []int, b SubscribeData) {
	go func(eventId, triggerId int, userIds []int, b SubscribeData) {
		send(userIds, eventId, constant.MENTION, triggerId, b, "")
	}(eventId, triggerId, userIds, b)

}

func (s *SubscriptionService) Send(eventId, eventType, fromId, toId int, b SubscribeData) {
	go func(eventId, fromId, toId int, b SubscribeData) {
		send([]int{toId}, eventId, eventType, fromId, b, "")
	}(eventId, fromId, toId, b)
}

func findAtUser(content string) []int {

	// 使用正则表达式匹配文本中的 @ 符号
	re := regexp.MustCompile(`@\((.*?)\)\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(content, -1)
	ids := mapset.NewSet[int]()
	for _, match := range matches {
		id, err := strconv.Atoi(match[2])
		if err != nil {
			log.Warnf("解析 @ 失败,err: %s", err.Error())
			continue
		}
		ids.Add(id)
	}
	return ids.ToSlice()
}

func send(userIds []int, eventId, eventType, sendId int, b SubscribeData, messageTemp string) {
	if len(userIds) == 0 {
		return
	}
	var m MessageService

	var users []model.Users
	model.User().Where("id in ?", userIds).Select("account", "id", "subscribe").Find(&users)

	var ids []int
	var emails []string
	for i := range users {
		ids = append(ids, users[i].ID)
		if users[i].Subscribe == 2 {
			emails = append(emails, users[i].Account)
		}
	}
	if messageTemp == "" {
		messageTemp = messageDao.GetMessageTemplate(eventId)
	}
	
	msg := m.GetMsgWithEventId(messageTemp, b, eventId)
	m.SendMessages(sendId, eventType, eventId, b.CurrentBusinessId, ids, msg)
	email.Send(emails, msg, "敲鸭社区")
}

/*
主动触发订阅事件
userId： 发送者
eventId：事件
messageType：消息类型
subscribeId：业务id
message：消息
*/
func (s *SubscriptionService) SendMsg(userId, eventId, messageType, subscribeId int, message string) {
	// 查出所有订阅人
	subscriptions := s.ListSubscriptionUserId(eventId, subscribeId)
	if len(subscriptions) == 0 {
		return
	}
	var userIds []int
	for _, subscription := range subscriptions {
		userIds = append(userIds, subscription.SubscriberId)
	}
	var userS UserService
	userMap := userS.ListByIdsToMap(userIds)
	var emails []string
	for _, v := range userMap {
		emails = append(emails, v.Account)
	}
	var m MessageService

	m.SendMessages(userId, messageType, eventId, subscribeId, userIds, message)
	email.Send(emails, message, "敲鸭社区")
}

func (s *SubscriptionService) SendMsgByToIds(userId, eventId, messageType, subscribeId int, toUserIds []int, message string) {

	var userS UserService
	userMap := userS.ListByIdsToMap(toUserIds)
	var emails []string
	for _, v := range userMap {
		emails = append(emails, v.Account)
	}
	var m MessageService

	m.SendMessages(userId, messageType, eventId, subscribeId, toUserIds, message)
	email.Send(emails, message, "敲鸭社区")
}
