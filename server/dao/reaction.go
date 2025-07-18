package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/model"
)

type ReactionDao struct{}

var reactionDao = &ReactionDao{}

// GetReactionDao 获取通用表情回复DAO实例
func GetReactionDao() *ReactionDao {
	return reactionDao
}

// AddReaction 添加表情回复
func (d *ReactionDao) AddReaction(reaction *model.Reaction) error {
	// 首先检查是否已存在相同的表情回复（只检查未删除的）
	var existing model.Reaction
	err := model.ReactionDB().Model(&model.Reaction{}).Where("business_type = ? AND business_id = ? AND user_id = ? AND reaction_type = ? AND deleted_at IS NULL",
		reaction.BusinessType, reaction.BusinessId, reaction.UserId, reaction.ReactionType).First(&existing).Error

	if err == nil {
		// 已存在，返回错误
		return errors.New("用户已对此业务添加过此表情")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("查询表情回复失败: %v", err)
		return err
	}

	// 不存在，创建新的表情回复
	return model.ReactionDB().Model(&model.Reaction{}).Create(reaction).Error
}

// RemoveReaction 移除表情回复
func (d *ReactionDao) RemoveReaction(businessType, businessId, userId int, reactionType string) error {
	// 使用硬删除而不是软删除，避免统计问题
	return model.ReactionDB().Model(&model.Reaction{}).Unscoped().Where("business_type = ? AND business_id = ? AND user_id = ? AND reaction_type = ?",
		businessType, businessId, userId, reactionType).Delete(&model.Reaction{}).Error
}

// GetReactionsByBusiness 获取业务的所有表情回复
func (d *ReactionDao) GetReactionsByBusiness(businessType, businessId int) ([]model.Reaction, error) {
	var reactions []model.Reaction
	err := model.ReactionDB().Model(&model.Reaction{}).
		Joins("LEFT JOIN users ON reactions.user_id = users.id").
		Select("reactions.*, users.name as user_name, users.avatar as user_avatar").
		Where("reactions.business_type = ? AND reactions.business_id = ?", businessType, businessId).
		Order("reactions.created_at ASC").
		Find(&reactions).Error

	return reactions, err
}

// GetReactionSummaryByBusiness 获取业务表情统计
func (d *ReactionDao) GetReactionSummaryByBusiness(businessType, businessId int, currentUserId int) ([]model.ReactionSummary, error) {
	// 创建临时结构体，避免GORM字段映射问题
	type TempSummary struct {
		BusinessType int    `json:"businessType"`
		BusinessId   int    `json:"businessId"`
		ReactionType string `json:"reactionType"`
		Count        int    `json:"count"`
		UserReacted  bool   `json:"userReacted"`
	}

	var tempSummaries []TempSummary

	// 首先获取基础统计
	basicQuery := `
		SELECT 
			r.business_type,
			r.business_id,
			r.reaction_type,
			COUNT(*) as count,
			CASE WHEN user_reactions.user_id IS NOT NULL THEN 1 ELSE 0 END as user_reacted
		FROM reactions r
		LEFT JOIN (
			SELECT business_type, business_id, reaction_type, user_id 
			FROM reactions 
			WHERE user_id = ? AND business_type = ? AND business_id = ?
		) user_reactions ON r.business_type = user_reactions.business_type 
			AND r.business_id = user_reactions.business_id
			AND r.reaction_type = user_reactions.reaction_type
		WHERE r.business_type = ? AND r.business_id = ? AND r.deleted_at IS NULL
		GROUP BY r.business_type, r.business_id, r.reaction_type
		ORDER BY r.reaction_type
	`

	err := model.ReactionDB().Raw(basicQuery, currentUserId, businessType, businessId, businessType, businessId).Scan(&tempSummaries).Error
	if err != nil {
		return nil, err
	}

	// 转换为最终结果格式
	var summaries []model.ReactionSummary
	for _, temp := range tempSummaries {
		summary := model.ReactionSummary{
			BusinessType: temp.BusinessType,
			BusinessId:   temp.BusinessId,
			ReactionType: temp.ReactionType,
			Count:        temp.Count,
			UserReacted:  temp.UserReacted,
			Users:        []model.ReactionUser{},
		}

		// 获取该表情类型的用户信息
		var users []model.ReactionUser
		userQuery := `
			SELECT 
				r.user_id,
				u.name as user_name,
				u.avatar as user_avatar
			FROM reactions r
			LEFT JOIN users u ON r.user_id = u.id
			WHERE r.business_type = ? AND r.business_id = ? AND r.reaction_type = ? AND r.deleted_at IS NULL
			ORDER BY r.created_at ASC
		`

		err = model.ReactionDB().Raw(userQuery, businessType, businessId, summary.ReactionType).Scan(&users).Error
		if err != nil {
			log.Errorf("获取表情用户信息失败: %v", err)
			continue
		}

		summary.Users = users
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

// GetReactionSummaryByBusinessBatch 批量获取多个业务的表情统计
func (d *ReactionDao) GetReactionSummaryByBusinessBatch(businessType int, businessIds []int, currentUserId int) (map[int][]model.ReactionSummary, error) {
	log.Infof("DAO层批量获取表情统计，业务类型: %d, 业务ID: %v, 当前用户ID: %d", businessType, businessIds, currentUserId)
	if len(businessIds) == 0 {
		log.Infof("业务ID列表为空，返回空结果")
		return make(map[int][]model.ReactionSummary), nil
	}

	// 创建临时结构体，避免GORM字段映射问题
	type TempSummary struct {
		BusinessType int    `json:"businessType"`
		BusinessId   int    `json:"businessId"`
		ReactionType string `json:"reactionType"`
		Count        int    `json:"count"`
		UserReacted  bool   `json:"userReacted"`
	}

	var tempSummaries []TempSummary

	// 批量查询多个业务的表情统计
	query := `
		SELECT 
			r.business_type,
			r.business_id,
			r.reaction_type,
			COUNT(*) as count,
			CASE WHEN user_reactions.user_id IS NOT NULL THEN 1 ELSE 0 END as user_reacted
		FROM reactions r
		LEFT JOIN (
			SELECT business_type, business_id, reaction_type, user_id 
			FROM reactions 
			WHERE user_id = ? AND business_type = ? AND business_id IN ?
		) user_reactions ON r.business_type = user_reactions.business_type 
			AND r.business_id = user_reactions.business_id
			AND r.reaction_type = user_reactions.reaction_type
		WHERE r.business_type = ? AND r.business_id IN ? AND r.deleted_at IS NULL
		GROUP BY r.business_type, r.business_id, r.reaction_type
		ORDER BY r.business_id, r.reaction_type
	`

	log.Infof("执行批量查询SQL")
	err := model.ReactionDB().Raw(query, currentUserId, businessType, businessIds, businessType, businessIds).Scan(&tempSummaries).Error
	if err != nil {
		log.Errorf("批量查询表情统计失败: %v", err)
		return nil, err
	}
	log.Infof("批量查询表情统计成功，获得 %d 条记录", len(tempSummaries))

	// 批量获取所有用户信息
	type UserReactionInfo struct {
		BusinessId   int    `json:"businessId"`
		ReactionType string `json:"reactionType"`
		UserId       int    `json:"userId"`
		UserName     string `json:"userName"`
		UserAvatar   string `json:"userAvatar"`
	}

	var userReactions []UserReactionInfo
	userQuery := `
		SELECT 
			r.business_id,
			r.reaction_type,
			r.user_id,
			u.name as user_name,
			u.avatar as user_avatar
		FROM reactions r
		LEFT JOIN users u ON r.user_id = u.id
		WHERE r.business_type = ? AND r.business_id IN ? AND r.deleted_at IS NULL
		ORDER BY r.business_id, r.reaction_type, r.created_at ASC
	`

	log.Infof("执行用户信息查询SQL")
	err = model.ReactionDB().Raw(userQuery, businessType, businessIds).Scan(&userReactions).Error
	if err != nil {
		log.Errorf("批量获取表情用户信息失败: %v", err)
	}
	log.Infof("批量获取用户信息成功，获得 %d 条记录", len(userReactions))

	// 组织用户信息到map中
	userMap := make(map[string][]model.ReactionUser)
	for _, ur := range userReactions {
		key := fmt.Sprintf("%d-%s", ur.BusinessId, ur.ReactionType)
		userMap[key] = append(userMap[key], model.ReactionUser{
			UserId:     ur.UserId,
			UserName:   ur.UserName,
			UserAvatar: ur.UserAvatar,
		})
	}
	log.Infof("组织用户信息到map，包含 %d 个键", len(userMap))

	// 转换为最终结果格式
	var summaries []model.ReactionSummary
	for _, temp := range tempSummaries {
		summary := model.ReactionSummary{
			BusinessType: temp.BusinessType,
			BusinessId:   temp.BusinessId,
			ReactionType: temp.ReactionType,
			Count:        temp.Count,
			UserReacted:  temp.UserReacted,
			Users:        []model.ReactionUser{},
		}

		key := fmt.Sprintf("%d-%s", temp.BusinessId, temp.ReactionType)
		if users, exists := userMap[key]; exists {
			summary.Users = users
		}

		summaries = append(summaries, summary)
	}

	// 组织结果为map格式
	result := make(map[int][]model.ReactionSummary)
	for _, summary := range summaries {
		result[summary.BusinessId] = append(result[summary.BusinessId], summary)
	}

	log.Infof("DAO层批量获取表情统计完成，返回 %d 个业务的统计信息", len(result))
	return result, nil
}

// CheckUserReaction 检查用户是否已对业务添加指定表情
func (d *ReactionDao) CheckUserReaction(businessType, businessId, userId int, reactionType string) (bool, error) {
	var count int64
	err := model.ReactionDB().Model(&model.Reaction{}).Where("business_type = ? AND business_id = ? AND user_id = ? AND reaction_type = ? AND deleted_at IS NULL",
		businessType, businessId, userId, reactionType).Count(&count).Error
	return count > 0, err
}