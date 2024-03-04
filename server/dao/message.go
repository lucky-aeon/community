package dao

import (
	"xhyovo.cn/community/pkg/mysql"
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

func (*MessageDao) SaveMessageTemplate(template model.MessageTemplates) {
	mysql.GetInstance().Save(&template)
}

func (*MessageDao) DeleteMessageTemplate(id int) {
	model.MessageState().Where("id = ?", id).Delete(nil)
}

// 消息日志crud
func (*MessageDao) ListMessageLogs(page, limit int) []*model.MessageLogs {
	var messageLogs []*model.MessageLogs
	model.MessageLog().Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&messageLogs)
	return messageLogs
}

// 添加记录
func (*MessageDao) SaveMessageLogs(messageLog []*model.MessageLogs) {
	model.MessageLog().Create(&messageLog)
}

func (*MessageDao) DeleteMessageLogs(id []int) {
	model.MessageLog().Delete(&id)
}

// 保存消息
func (*MessageDao) SaveMessage(from, types, businessId int, to []int, content string) {
	var msgs []*model.MessageStates
	for i := range to {
		state := &model.MessageStates{
			From:      from,
			To:        to[i],
			Content:   content,
			Type:      types,
			State:     1,
			ArticleId: businessId,
		}
		msgs = append(msgs, state)
	}

	model.MessageState().Create(&msgs)
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

func (d *MessageDao) CountMessage(userId, types, state int) int64 {
	var count int64

	model.MessageState().Where(model.MessageStates{To: userId, Type: types, State: state}).Count(&count)
	return count
}
