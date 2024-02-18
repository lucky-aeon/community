package services

import (
	"fmt"
	"reflect"
	"strings"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type MessageService struct {
}

func (*MessageService) ListMessageTemplate(page, limit int) []*model.MessageTemplates {

	return messageDao.ListMessageTemplate(page, limit)
}

func (*MessageService) SaveMessageTemplate(template *model.MessageTemplates) {
	messageDao.SaveMessageTemplate(template)
}

func (*MessageService) DeleteMessageTemplate(id []int) {
	messageDao.DeleteMessageTemplate(id)
}

func (*MessageService) ListMessageLogs(page, limit int) []*model.MessageLogs {
	return messageDao.ListMessageLogs(page, limit)
}

func (*MessageService) AddMessageLogs(from int, to []int, content string) {
	var messageLogs []*model.MessageLogs
	for i := range to {
		log := &model.MessageLogs{
			From:    from,
			To:      i,
			Content: content,
		}
		messageLogs = append(messageLogs, log)
	}
	messageDao.SaveMessageLogs(messageLogs)
}

func (*MessageService) DeleteMessageLogs(id []int) {
	messageDao.DeleteMessageLogs(id)
}

func (m *MessageService) SendMessage(from, to, types int, content string) {

	messageDao.SaveMessage(from, types, []int{to}, content)
	// 添加记录
	m.AddMessageLogs(from, []int{to}, content)
}

func (m *MessageService) SendMessages(from, types int, to []int, content string) {

	messageDao.SaveMessage(from, types, to, content)
	// 添加记录
	m.AddMessageLogs(from, to, content)
}

func (*MessageService) ReadMessage(id []int, userId int) int64 {
	return messageDao.ReadMessage(id, userId)
}

func (m *MessageService) PageMessage(page, limit, userId, types, state int) (msgs []*model.MessageStates, count int64) {
	msgs = messageDao.ListMessage(page, limit, userId, types, state)
	count = messageDao.CountMessage(userId, types)
	return msgs, count
}

// 人：你订阅的 xxx 用户发布了文章，文章标题：xxx
// 文章：你订阅的 xxx 文章，被xxx评论了，评论内容：xxx
func (m *MessageService) GetMsg(template string, eventId, businessId int) string {
	// 根据事件类型获取对象
	var object any
	if eventId == event.CommentUpdateEvent {
		// 查询文章
		var articleS ArticleService
		object = articleS.GetById(businessId)

	} else if eventId == event.UserFollowingEvent {
		// 查询用户
		object = userDao.GetById(businessId)
	}
	v := reflect.ValueOf(object)
	return setValue(template, v)
}

func setValue(template string, v reflect.Value) string {
	var msg = template
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := t.Field(i).Name

		if field.Kind() == reflect.Struct {
			// 如果字段是结构体类型，则递归调用处理
			msg = setValue(msg, field)
		} else {
			if fieldName[0] >= 'a' && fieldName[0] <= 'z' {
				continue
			}
			fieldValue := fmt.Sprintf("%v", field.Interface())
			fieldName = fmt.Sprintf("${%s}", strings.ToLower(fieldName))
			msg = strings.Replace(msg, fieldName, fieldValue, -1)
		}
	}
	return msg
}
