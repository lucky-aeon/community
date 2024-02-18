package dao

import (
	"strconv"
	"xhyovo.cn/community/server/model"
)

type SubscriptionDao struct {
}

// 查看订阅列表
func (*SubscriptionDao) ListSubscription(userId, event, page, limit int) []model.Subscriptions {
	var subscriptions []model.Subscriptions

	model.Subscription().Where(&model.Subscriptions{SubscriberId: userId, EventId: event}).Offset((page - 1) * limit).Limit(limit).Find(&subscriptions)
	return subscriptions
}

// 查看对应事件订阅状态
func (*SubscriptionDao) SubscriptionState(subscriptions *model.Subscriptions) bool {
	var count int64
	model.Subscription().Where(subscriptions).Count(&count)
	return count == 1
}

// 订阅/取消订阅
func (s *SubscriptionDao) Subscribe(subscription *model.Subscriptions) bool {
	subscription.IndexKey = strconv.Itoa(subscription.SubscriberId) + strconv.Itoa(subscription.EventId) + strconv.Itoa(subscription.BusinessId)
	tx := model.Subscription().Save(&subscription)
	if tx.Error != nil {
		s.cancelSubscribe(subscription)
		return false
	}
	return true
}

// 取消订阅
func (*SubscriptionDao) cancelSubscribe(subscription *model.Subscriptions) {
	model.Subscription().Where("index_key = ? ", subscription.IndexKey).Delete(&subscription)
}

func (s *SubscriptionDao) ListSubscriptions(event, businessId int) []model.Subscriptions {
	var sub []model.Subscriptions
	model.Subscription().Where("event_id = ? and business_id = ?", event, businessId).Find(&sub)
	return sub
}
