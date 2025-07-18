package services

import (
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type ReactionService struct {
	ctx *gin.Context
}

// NewReactionService 创建通用表情回复服务实例
func NewReactionService(ctx *gin.Context) *ReactionService {
	return &ReactionService{ctx: ctx}
}

// ToggleReaction 切换表情回复状态
func (s *ReactionService) ToggleReaction(businessType, businessId, userId int, reactionType string) (bool, error) {
	log.Infof("用户ID %d 切换表情回复，业务类型: %d, 业务ID: %d, 表情类型: %s", userId, businessType, businessId, reactionType)
	
	reactionDao := dao.GetReactionDao()
	
	// 检查用户是否已经添加过此表情
	exists, err := reactionDao.CheckUserReaction(businessType, businessId, userId, reactionType)
	if err != nil {
		log.Errorf("检查用户表情回复状态失败: %v", err)
		return false, err
	}
	
	if exists {
		// 已存在，移除表情回复
		err = reactionDao.RemoveReaction(businessType, businessId, userId, reactionType)
		if err != nil {
			log.Errorf("移除表情回复失败: %v", err)
			return false, err
		}
		log.Infof("用户ID %d 移除表情回复成功", userId)
		return false, nil
	} else {
		// 不存在，添加表情回复
		reaction := &model.Reaction{
			BusinessType: businessType,
			BusinessId:   businessId,
			UserId:       userId,
			ReactionType: reactionType,
		}
		err = reactionDao.AddReaction(reaction)
		if err != nil {
			log.Errorf("添加表情回复失败: %v", err)
			return false, err
		}
		log.Infof("用户ID %d 添加表情回复成功", userId)
		return true, nil
	}
}

// GetReactionSummary 获取单个业务的表情统计
func (s *ReactionService) GetReactionSummary(businessType, businessId, currentUserId int) ([]model.ReactionSummary, error) {
	log.Infof("获取表情统计，业务类型: %d, 业务ID: %d, 当前用户ID: %d", businessType, businessId, currentUserId)
	
	reactionDao := dao.GetReactionDao()
	summaries, err := reactionDao.GetReactionSummaryByBusiness(businessType, businessId, currentUserId)
	if err != nil {
		log.Errorf("获取表情统计失败: %v", err)
		return nil, err
	}
	
	log.Infof("获取表情统计成功，返回 %d 个表情类型", len(summaries))
	return summaries, nil
}

// GetReactionSummaryBatch 批量获取多个业务的表情统计
func (s *ReactionService) GetReactionSummaryBatch(businessType int, businessIds []int, currentUserId int) (map[int][]model.ReactionSummary, error) {
	log.Infof("批量获取表情统计，业务类型: %d, 业务ID数量: %d, 当前用户ID: %d", businessType, len(businessIds), currentUserId)
	
	if len(businessIds) == 0 {
		return make(map[int][]model.ReactionSummary), nil
	}
	
	reactionDao := dao.GetReactionDao()
	summaryMap, err := reactionDao.GetReactionSummaryByBusinessBatch(businessType, businessIds, currentUserId)
	if err != nil {
		log.Errorf("批量获取表情统计失败: %v", err)
		return nil, err
	}
	
	log.Infof("批量获取表情统计成功，返回 %d 个业务的统计信息", len(summaryMap))
	return summaryMap, nil
}

// ValidateBusinessType 验证业务类型是否有效
func (s *ReactionService) ValidateBusinessType(businessType int) bool {
	return businessType >= model.BusinessTypeArticle && businessType <= model.BusinessTypeAINews
}

// ValidateReactionType 验证表情类型是否有效
func (s *ReactionService) ValidateReactionType(reactionType string) bool {
	// 先进行基本的非空验证
	if reactionType == "" {
		return false
	}
	
	// 从数据库查询表情类型是否存在且启用
	var count int64
	err := model.ReactionDB().Model(&model.ExpressionType{}).
		Where("code = ? AND is_active = ?", reactionType, true).
		Count(&count).Error
	
	if err != nil {
		log.Errorf("验证表情类型失败: %v", err)
		return false
	}
	
	return count > 0
}