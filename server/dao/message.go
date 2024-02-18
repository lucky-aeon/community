package dao

import (
	"xhyovo.cn/community/server/model"
)

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

// 删除用户收到的消息(确认消息),
func (*MessageDao) ReadMessage(id []int, userId int) int64 {
	tx := model.MessageState().Where("id in ? and `to` = ?", id, userId).Updates(&model.MessageStates{State: 1})
	return tx.RowsAffected
}

func (d *MessageDao) ListMessage(page, limit, userId, types, state int) []*model.MessageStates {
	m := model.MessageStates{
		To:    userId,
		Type:  types,
		State: state,
	}
	var message []*model.MessageStates
	model.MessageState().Where(&m).Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&message)
	return message
}

func (d *MessageDao) CountMessage(userId int, types int) int64 {
	var count int64
	model.MessageState().Where("`to` = ? and type = ?", userId, types).Count(&count)
	return count
}
