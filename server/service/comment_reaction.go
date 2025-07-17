package services

import (
	"errors"
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type CommentReactionService struct {
	ctx *gin.Context
}

func NewCommentReactionService(ctx *gin.Context) *CommentReactionService {
	return &CommentReactionService{ctx: ctx}
}

// AddReaction 添加表情回复
func (s *CommentReactionService) AddReaction(commentId, userId int, reactionType string) error {
	// 验证表情类型是否有效
	if !s.isValidReactionType(reactionType) {
		return errors.New("无效的表情类型")
	}

	// 检查评论是否存在
	if !s.isCommentExists(commentId) {
		return errors.New("评论不存在")
	}

	// 创建表情回复
	reaction := &model.CommentReaction{
		CommentId:    commentId,
		UserId:       userId,
		ReactionType: reactionType,
	}

	err := dao.GetCommentReactionDao().AddReaction(reaction)
	if err != nil {
		log.Errorf("添加表情回复失败: %v", err)
		return err
	}

	log.Infof("用户 %d 对评论 %d 添加了表情回复: %s", userId, commentId, reactionType)
	return nil
}

// RemoveReaction 移除表情回复
func (s *CommentReactionService) RemoveReaction(commentId, userId int, reactionType string) error {
	// 验证表情类型是否有效
	if !s.isValidReactionType(reactionType) {
		return errors.New("无效的表情类型")
	}

	// 检查用户是否已添加此表情
	hasReaction, err := dao.GetCommentReactionDao().CheckUserReaction(commentId, userId, reactionType)
	if err != nil {
		log.Errorf("检查用户表情回复失败: %v", err)
		return err
	}

	if !hasReaction {
		return errors.New("用户未对此评论添加过此表情")
	}

	err = dao.GetCommentReactionDao().RemoveReaction(commentId, userId, reactionType)
	if err != nil {
		log.Errorf("移除表情回复失败: %v", err)
		return err
	}

	log.Infof("用户 %d 移除了对评论 %d 的表情回复: %s", userId, commentId, reactionType)
	return nil
}

// GetReactionSummary 获取评论表情统计
func (s *CommentReactionService) GetReactionSummary(commentId, currentUserId int) ([]model.CommentReactionSummary, error) {
	summaries, err := dao.GetCommentReactionDao().GetReactionSummaryByCommentId(commentId, currentUserId)
	if err != nil {
		log.Errorf("获取评论表情统计失败: %v", err)
		return nil, err
	}

	return summaries, nil
}

// GetReactionSummaryBatch 批量获取多个评论的表情统计
func (s *CommentReactionService) GetReactionSummaryBatch(commentIds []int, currentUserId int) (map[int][]model.CommentReactionSummary, error) {
	log.Infof("批量获取表情统计，评论ID: %v, 当前用户ID: %d", commentIds, currentUserId)
	summaries, err := dao.GetCommentReactionDao().GetReactionSummaryByCommentIds(commentIds, currentUserId)
	if err != nil {
		log.Errorf("批量获取评论表情统计失败: %v", err)
		return nil, err
	}
	log.Infof("批量获取表情统计成功，结果数量: %d", len(summaries))
	for commentId, reactions := range summaries {
		log.Infof("评论ID %d 有 %d 个表情统计", commentId, len(reactions))
	}

	return summaries, nil
}

// GetReactionDetails 获取评论表情详情
func (s *CommentReactionService) GetReactionDetails(commentId int) ([]model.CommentReaction, error) {
	reactions, err := dao.GetCommentReactionDao().GetReactionsByCommentId(commentId)
	if err != nil {
		log.Errorf("获取评论表情详情失败: %v", err)
		return nil, err
	}

	return reactions, nil
}

// GetAllExpressionTypes 获取所有表情类型
func (s *CommentReactionService) GetAllExpressionTypes() ([]model.ExpressionType, error) {
	types, err := dao.GetCommentReactionDao().GetAllExpressionTypes()
	if err != nil {
		log.Errorf("获取表情类型失败: %v", err)
		return nil, err
	}

	return types, nil
}

// ToggleReaction 切换表情回复状态（添加或移除）
func (s *CommentReactionService) ToggleReaction(commentId, userId int, reactionType string) (bool, error) {
	log.Infof("切换表情回复状态，评论ID: %d, 用户ID: %d, 表情类型: %s", commentId, userId, reactionType)
	
	// 检查用户是否已添加此表情
	hasReaction, err := dao.GetCommentReactionDao().CheckUserReaction(commentId, userId, reactionType)
	if err != nil {
		log.Errorf("检查用户表情回复状态失败: %v", err)
		return false, err
	}
	
	log.Infof("用户当前表情状态: %t", hasReaction)

	if hasReaction {
		// 已存在，移除
		log.Infof("移除表情回复")
		err = s.RemoveReaction(commentId, userId, reactionType)
		if err != nil {
			log.Errorf("移除表情回复失败: %v", err)
		} else {
			log.Infof("移除表情回复成功")
		}
		return false, err
	} else {
		// 不存在，添加
		log.Infof("添加表情回复")
		err = s.AddReaction(commentId, userId, reactionType)
		if err != nil {
			log.Errorf("添加表情回复失败: %v", err)
		} else {
			log.Infof("添加表情回复成功")
		}
		return true, err
	}
}

// isValidReactionType 验证表情类型是否有效
func (s *CommentReactionService) isValidReactionType(reactionType string) bool {
	_, err := dao.GetCommentReactionDao().GetExpressionTypeByCode(reactionType)
	return err == nil
}

// isCommentExists 检查评论是否存在
func (s *CommentReactionService) isCommentExists(commentId int) bool {
	var commentService CommentsService
	comment := commentService.GetById(commentId)
	return comment.ID > 0
}

// GetCommentReactionDao 获取评论表情回复DAO实例
func GetCommentReactionDao() *dao.CommentReactionDao {
	return &dao.CommentReactionDao{}
}
