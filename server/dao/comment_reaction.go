package dao

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/model"
)

type CommentReactionDao struct{}

var commentReactionDao = &CommentReactionDao{}

// GetCommentReactionDao 获取评论表情回复DAO实例
func GetCommentReactionDao() *CommentReactionDao {
	return commentReactionDao
}

// AddReaction 添加表情回复
func (d *CommentReactionDao) AddReaction(reaction *model.CommentReaction) error {
	// 首先检查是否已存在相同的表情回复（只检查未删除的）
	var existing model.CommentReaction
	err := model.CommentReactionDB().Where("comment_id = ? AND user_id = ? AND reaction_type = ? AND deleted_at IS NULL", 
		reaction.CommentId, reaction.UserId, reaction.ReactionType).First(&existing).Error
	
	if err == nil {
		// 已存在，返回错误
		return errors.New("用户已对此评论添加过此表情")
	}
	
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorf("查询表情回复失败: %v", err)
		return err
	}
	
	// 不存在，创建新的表情回复
	return model.CommentReactionDB().Create(reaction).Error
}

// RemoveReaction 移除表情回复
func (d *CommentReactionDao) RemoveReaction(commentId, userId int, reactionType string) error {
	// 使用硬删除而不是软删除，避免统计问题
	return model.CommentReactionDB().Unscoped().Where("comment_id = ? AND user_id = ? AND reaction_type = ?", 
		commentId, userId, reactionType).Delete(&model.CommentReaction{}).Error
}

// GetReactionsByCommentId 获取评论的所有表情回复
func (d *CommentReactionDao) GetReactionsByCommentId(commentId int) ([]model.CommentReaction, error) {
	var reactions []model.CommentReaction
	err := model.CommentReactionDB().
		Joins("LEFT JOIN users ON comment_reactions.user_id = users.id").
		Select("comment_reactions.*, users.name as user_name, users.avatar as user_avatar").
		Where("comment_reactions.comment_id = ?", commentId).
		Order("comment_reactions.created_at ASC").
		Find(&reactions).Error
	
	return reactions, err
}

// GetReactionSummaryByCommentId 获取评论表情统计
func (d *CommentReactionDao) GetReactionSummaryByCommentId(commentId int, currentUserId int) ([]model.CommentReactionSummary, error) {
	var summaries []model.CommentReactionSummary
	
	// 首先获取基础统计
	basicQuery := `
		SELECT 
			cr.comment_id,
			cr.reaction_type,
			COUNT(*) as count,
			CASE WHEN user_reactions.user_id IS NOT NULL THEN 1 ELSE 0 END as user_reacted
		FROM comment_reactions cr
		LEFT JOIN (
			SELECT comment_id, reaction_type, user_id 
			FROM comment_reactions 
			WHERE user_id = ? AND comment_id = ?
		) user_reactions ON cr.comment_id = user_reactions.comment_id 
			AND cr.reaction_type = user_reactions.reaction_type
		WHERE cr.comment_id = ? AND cr.deleted_at IS NULL
		GROUP BY cr.comment_id, cr.reaction_type
		ORDER BY cr.reaction_type
	`
	
	err := model.CommentReactionDB().Raw(basicQuery, currentUserId, commentId, commentId).Scan(&summaries).Error
	if err != nil {
		return nil, err
	}
	
	// 为每个表情类型获取用户信息
	for i := range summaries {
		summary := &summaries[i]
		
		// 获取该表情类型的用户信息
		var users []model.ReactionUser
		userQuery := `
			SELECT 
				cr.user_id,
				u.name as user_name,
				u.avatar as user_avatar
			FROM comment_reactions cr
			LEFT JOIN users u ON cr.user_id = u.id
			WHERE cr.comment_id = ? AND cr.reaction_type = ? AND cr.deleted_at IS NULL
			ORDER BY cr.created_at ASC
		`
		
		err = model.CommentReactionDB().Raw(userQuery, commentId, summary.ReactionType).Scan(&users).Error
		if err != nil {
			log.Errorf("获取表情用户信息失败: %v", err)
			continue
		}
		
		summary.Users = users
	}
	
	return summaries, err
}

// GetReactionSummaryByCommentIds 批量获取多个评论的表情统计
func (d *CommentReactionDao) GetReactionSummaryByCommentIds(commentIds []int, currentUserId int) (map[int][]model.CommentReactionSummary, error) {
	log.Infof("DAO层批量获取表情统计，评论ID: %v, 当前用户ID: %d", commentIds, currentUserId)
	if len(commentIds) == 0 {
		log.Infof("评论ID列表为空，返回空结果")
		return make(map[int][]model.CommentReactionSummary), nil
	}
	
	// 创建临时结构体，避免GORM字段映射问题
	type TempSummary struct {
		CommentId    int  `json:"commentId"`
		ReactionType string `json:"reactionType"`
		Count        int    `json:"count"`
		UserReacted  bool   `json:"userReacted"`
	}
	
	var tempSummaries []TempSummary
	
	// 批量查询多个评论的表情统计
	query := `
		SELECT 
			cr.comment_id,
			cr.reaction_type,
			COUNT(*) as count,
			CASE WHEN user_reactions.user_id IS NOT NULL THEN 1 ELSE 0 END as user_reacted
		FROM comment_reactions cr
		LEFT JOIN (
			SELECT comment_id, reaction_type, user_id 
			FROM comment_reactions 
			WHERE user_id = ? AND comment_id IN ?
		) user_reactions ON cr.comment_id = user_reactions.comment_id 
			AND cr.reaction_type = user_reactions.reaction_type
		WHERE cr.comment_id IN ? AND cr.deleted_at IS NULL
		GROUP BY cr.comment_id, cr.reaction_type
		ORDER BY cr.comment_id, cr.reaction_type
	`
	
	log.Infof("执行批量查询SQL")
	err := model.CommentReactionDB().Raw(query, currentUserId, commentIds, commentIds).Scan(&tempSummaries).Error
	if err != nil {
		log.Errorf("批量查询表情统计失败: %v", err)
		return nil, err
	}
	log.Infof("批量查询表情统计成功，获得 %d 条记录", len(tempSummaries))
	
	// 批量获取所有用户信息
	type UserReactionInfo struct {
		CommentId    int    `json:"commentId"`
		ReactionType string `json:"reactionType"`
		UserId       int    `json:"userId"`
		UserName     string `json:"userName"`
		UserAvatar   string `json:"userAvatar"`
	}
	
	var userReactions []UserReactionInfo
	userQuery := `
		SELECT 
			cr.comment_id,
			cr.reaction_type,
			cr.user_id,
			u.name as user_name,
			u.avatar as user_avatar
		FROM comment_reactions cr
		LEFT JOIN users u ON cr.user_id = u.id
		WHERE cr.comment_id IN ? AND cr.deleted_at IS NULL
		ORDER BY cr.comment_id, cr.reaction_type, cr.created_at ASC
	`
	
	log.Infof("执行用户信息查询SQL")
	err = model.CommentReactionDB().Raw(userQuery, commentIds).Scan(&userReactions).Error
	if err != nil {
		log.Errorf("批量获取表情用户信息失败: %v", err)
	}
	log.Infof("批量获取用户信息成功，获得 %d 条记录", len(userReactions))
	
	// 组织用户信息到map中
	userMap := make(map[string][]model.ReactionUser)
	for _, ur := range userReactions {
		key := fmt.Sprintf("%d-%s", ur.CommentId, ur.ReactionType)
		userMap[key] = append(userMap[key], model.ReactionUser{
			UserId:     ur.UserId,
			UserName:   ur.UserName,
			UserAvatar: ur.UserAvatar,
		})
	}
	log.Infof("组织用户信息到map，包含 %d 个键", len(userMap))
	
	// 转换为最终结果格式
	var summaries []model.CommentReactionSummary
	for _, temp := range tempSummaries {
		summary := model.CommentReactionSummary{
			CommentId:    temp.CommentId,
			ReactionType: temp.ReactionType,
			Count:        temp.Count,
			UserReacted:  temp.UserReacted,
			Users:        []model.ReactionUser{},
		}
		
		key := fmt.Sprintf("%d-%s", temp.CommentId, temp.ReactionType)
		if users, exists := userMap[key]; exists {
			summary.Users = users
		}
		
		summaries = append(summaries, summary)
	}
	
	// 组织结果为map格式
	result := make(map[int][]model.CommentReactionSummary)
	for _, summary := range summaries {
		result[summary.CommentId] = append(result[summary.CommentId], summary)
	}
	
	log.Infof("DAO层批量获取表情统计完成，返回 %d 个评论的统计信息", len(result))
	return result, nil
}

// GetAllExpressionTypes 获取所有表情类型
func (d *CommentReactionDao) GetAllExpressionTypes() ([]model.ExpressionType, error) {
	var types []model.ExpressionType
	err := model.ExpressionTypeDB().Where("is_active = ?", true).
		Order("sort_order ASC").Find(&types).Error
	return types, err
}

// GetExpressionTypeByCode 根据代码获取表情类型
func (d *CommentReactionDao) GetExpressionTypeByCode(code string) (*model.ExpressionType, error) {
	var expressionType model.ExpressionType
	err := model.ExpressionTypeDB().Where("code = ? AND is_active = ?", code, true).First(&expressionType).Error
	if err != nil {
		return nil, err
	}
	return &expressionType, nil
}

// CheckUserReaction 检查用户是否已对评论添加指定表情
func (d *CommentReactionDao) CheckUserReaction(commentId, userId int, reactionType string) (bool, error) {
	var count int64
	err := model.CommentReactionDB().Where("comment_id = ? AND user_id = ? AND reaction_type = ? AND deleted_at IS NULL", 
		commentId, userId, reactionType).Count(&count).Error
	return count > 0, err
}