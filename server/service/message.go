package services

import (
	"xhyovo.cn/community/server/model"
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

func (*MessageService) AddMessageLogs(messageLog *model.MessageLogs) {
	messageDao.SaveMessageLogs(messageLog)
}

func (*MessageService) DeleteMessageLogs(id []int) {
	messageDao.DeleteMessageLogs(id)
}

func (m *MessageService) SendMessage(from, to, contentType uint) {

	// 不会搞，后续再说 todo
	messageDao.SendMessage(from, to, "")
	// 添加记录
	m.AddMessageLogs(&model.MessageLogs{From: from, To: to, Content: ""})
}

func (*MessageService) ListNoReadMessage(page, limit, userId int) []*model.MessageStates {
	return messageDao.ListNoReadMessage(page, limit, userId)
}

func (*MessageService) ReadMessage(id []int, userId int) {
	messageDao.DeleteMessage(id, userId)
}
