package dao

import (
	"xhyovo.cn/community/server/model"
)

type Chat struct {
}

// QueryAIModels 获取所有启用的AI模型
func (c *Chat) QueryAIModels() ([]model.AIModels, error) {
	var models []model.AIModels
	err := model.AIModelDB().Where("status = ?", true).Find(&models).Error
	return models, err
}

// GetAIModelByID 根据ID获取AI模型
func (c *Chat) GetAIModelByID(id int64) (*model.AIModels, error) {
	var modelObject model.AIModels

	err := model.AIModelDB().Where("id = ? AND status = ?", id, true).First(&modelObject).Error
	if err != nil {
		return nil, err
	}
	return &modelObject, nil
}

// CreateChatGroup 创建对话分组
func (c *Chat) CreateChatGroup(group *model.ChatGroups) error {
	return model.ChatGroupDB().Create(group).Error
}

// GetChatGroup 根据ID获取对话分组
func (c *Chat) GetChatGroup(id int64) (*model.ChatGroups, error) {
	var group model.ChatGroups
	err := model.ChatGroupDB().Where("id = ? AND is_deleted = ?", id, false).First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetUserChatGroups 获取用户的所有对话分组
func (c *Chat) GetUserChatGroups(userID int64, page, limit int) ([]model.ChatGroups, error) {
	var groups []model.ChatGroups
	if limit < 1 {
		limit = 10
	}
	if page < 1 {
		page = 1
	}
	err := model.ChatGroupDB().
		Where("user_id = ? AND is_deleted = ?", userID, false).
		Offset((page - 1) * limit).
		Order("created_at desc").
		Limit(limit).
		Find(&groups).Error
	return groups, err
}

// CountUserChatGroups 获取用户的对话分组总数
func (c *Chat) CountUserChatGroups(userID int64) int64 {
	var count int64
	model.ChatGroupDB().Where("user_id = ? AND is_deleted = ?", userID, false).Count(&count)
	return count
}

// UpdateChatGroup 更新对话分组
func (c *Chat) UpdateChatGroup(group *model.ChatGroups) error {
	return model.ChatGroupDB().Model(group).Updates(map[string]interface{}{
		"title": group.Title,
	}).Error
}

// DeleteChatGroup 软删除对话分组
func (c *Chat) DeleteChatGroup(id int64) error {
	return model.ChatGroupDB().Where("id = ?", id).Update("is_deleted", true).Error
}

// CreateChatMessage 创建对话消息
func (c *Chat) CreateChatMessage(message *model.ChatMessages) error {
	return model.ChatMessageDB().Create(message).Error
}

// GetChatMessages 获取对话消息列表
func (c *Chat) GetChatMessages(groupID int64, page, pageSize int) ([]model.ChatMessages, error) {
	var messages []model.ChatMessages
	if pageSize < 1 {
		pageSize = 20
	}
	if page < 1 {
		page = 1
	}
	err := model.ChatMessageDB().
		Where("group_id = ?", groupID).
		Order("created_at desc").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&messages).Error
	return messages, err
}

// GetRecentMessages 获取最近的消息用于上下文
func (c *Chat) GetRecentMessages(groupID int64, limit int) ([]model.ChatMessages, error) {
	var messages []model.ChatMessages
	err := model.ChatMessageDB().
		Where("group_id = ?", groupID).
		Order("created_at desc").
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	// 反转消息顺序，使其按时间正序排列
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}
	return messages, nil
}

// GetTotalMessages 获取对话组的消息总数
func (c *Chat) GetTotalMessages(groupID int64) int64 {
	var count int64
	model.ChatMessageDB().Where("group_id = ?", groupID).Count(&count)
	return count
}

// ExistsByID 检查对话分组是否存在
func (c *Chat) ExistsByID(id int64) bool {
	var count int64
	model.ChatGroupDB().Where("id = ? AND is_deleted = ?", id, false).Count(&count)
	return count == 1
}
