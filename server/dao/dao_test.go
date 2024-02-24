package dao_test

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"strconv"
	"testing"
	"xhyovo.cn/community/server/model"
)

func getDb() *gorm.DB {
	d := "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"

	dsn := fmt.Sprintf(d, "root", "123", "127.0.0.1:3306", "community")
	instance, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
		Logger:         logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return instance
}
func Test1(t *testing.T) {

	var subscriptions []model.Subscriptions
	userId := 1
	event := 0
	getDb().Where(&model.Subscriptions{SubscriberId: userId, EventId: event}).Find(&subscriptions)
	for i := range subscriptions {
		m := subscriptions[i]
		fmt.Println("subscription event:", m.EventId)
	}
}

func Test2(t *testing.T) {
	subscription := &model.Subscriptions{SubscriberId: 1, EventId: 1, BusinessId: 1}
	subscription.IndexKey = strconv.Itoa(subscription.SubscriberId) + strconv.Itoa(subscription.EventId) + strconv.Itoa(subscription.BusinessId)

	err := getDb().Save(subscription).Error
	if err != nil {
		fmt.Println("存在咯")
	} else {
		fmt.Println("不")
	}
}

func Test3(t *testing.T) {
	getDb().Where(&model.Subscriptions{ID: 1, SubscriberId: 1, EventId: 1, BusinessId: 1}).Delete(&model.Subscriptions{})
}

func Test4(t *testing.T) {
	var count int64

	getDb().Model(&model.Subscriptions{}).Where(&model.Subscriptions{SubscriberId: 1, EventId: 1, BusinessId: 1}).Count(&count)
	fmt.Println(count)
	fmt.Println(count == 1)
}
