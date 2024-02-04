package dao

import "xhyovo.cn/community/server/model"

type MessageDao struct {
}

// 消息模板crud
func (*MessageDao) ListMessageTemplate(page, limit int) []*model.MessageTemplates {
	var templates []*model.MessageTemplates
	model.MessageTemplate().Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&templates)
	return templates
}

func (*MessageDao) GetMessageTemplate(id int) string {
	var messageTemplate string
	model.MessageTemplate().Where("event_id = ?", id).Select("content").Find(&messageTemplate)
	return messageTemplate
}

// todo 模板填充

func (*MessageDao) SaveMessageTemplate(template *model.MessageTemplates) {
	model.MessageTemplate().Save(&template)
}

func (*MessageDao) DeleteMessageTemplate(id []int) {
	model.MessageState().Delete(&id)
}

// 消息日志crud
func (*MessageDao) ListMessageLogs(page, limit int) []*model.MessageLogs {
	var messageLogs []*model.MessageLogs
	model.MessageLog().Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&messageLogs)
	return messageLogs
}

// 添加记录
func (*MessageDao) SaveMessageLogs(messageLog *model.MessageLogs) {
	model.MessageLog().Save(messageLog)
}

func (*MessageDao) DeleteMessageLogs(id []int) {
	model.MessageLog().Delete(&id)
}

// 发送消息
func (*MessageDao) SendMessage(from, to int, content string) {
	state := &model.MessageStates{
		From:    from,
		To:      to,
		Content: content,
	}
	model.MessageState().Save(&state)
}

// 查看用户未读消息
func (*MessageDao) ListNoReadMessage(page, limit, userId int) []*model.MessageStates {
	var message []*model.MessageStates
	model.MessageState().Where("to", userId).Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&message)
	return message
}

// 删除用户收到的消息(确认消息),
func (*MessageDao) DeleteMessage(id []int, userId int) {

	model.MessageState().Where("to", userId).Delete(&model.MessageStates{}, id)
}
