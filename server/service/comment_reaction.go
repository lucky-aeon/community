package services

import (
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

// AddReaction 添加表情回复（使用通用系统）
func (s *CommentReactionService) AddReaction(commentId, userId int, reactionType string) error {
	// 使用通用表情回复系统
	reactionService := NewReactionService(s.ctx)
	_, err := reactionService.ToggleReaction(model.BusinessTypeComment, commentId, userId, reactionType)
	if err != nil {
		log.Errorf("添加评论表情回复失败: %v", err)
		return err
	}
	log.Infof("用户 %d 对评论 %d 添加了表情回复: %s", userId, commentId, reactionType)
	return nil
}

// RemoveReaction 移除表情回复（使用通用系统）
func (s *CommentReactionService) RemoveReaction(commentId, userId int, reactionType string) error {
	// 使用通用表情回复系统
	reactionService := NewReactionService(s.ctx)
	_, err := reactionService.ToggleReaction(model.BusinessTypeComment, commentId, userId, reactionType)
	if err != nil {
		log.Errorf("移除评论表情回复失败: %v", err)
		return err
	}
	log.Infof("用户 %d 移除了对评论 %d 的表情回复: %s", userId, commentId, reactionType)
	return nil
}

// GetReactionSummary 获取评论表情统计（使用通用系统）
func (s *CommentReactionService) GetReactionSummary(commentId, currentUserId int) ([]model.CommentReactionSummary, error) {
	// 使用通用表情回复系统
	reactionService := NewReactionService(s.ctx)
	summaries, err := reactionService.GetReactionSummary(model.BusinessTypeComment, commentId, currentUserId)
	if err != nil {
		log.Errorf("获取评论表情统计失败: %v", err)
		return nil, err
	}

	// 转换为评论表情统计格式
	commentSummaries := make([]model.CommentReactionSummary, len(summaries))
	for i, summary := range summaries {
		commentSummaries[i] = model.CommentReactionSummary{
			ReactionType: summary.ReactionType,
			Count:        summary.Count,
			UserReacted:  summary.UserReacted,
			Users:        summary.Users,
		}
	}

	return commentSummaries, nil
}

// GetReactionSummaryBatch 批量获取多个评论的表情统计（使用通用系统）
func (s *CommentReactionService) GetReactionSummaryBatch(commentIds []int, currentUserId int) (map[int][]model.CommentReactionSummary, error) {
	log.Infof("批量获取表情统计，评论ID: %v, 当前用户ID: %d", commentIds, currentUserId)
	
	// 使用通用表情回复系统
	reactionService := NewReactionService(s.ctx)
	summaryMap, err := reactionService.GetReactionSummaryBatch(model.BusinessTypeComment, commentIds, currentUserId)
	if err != nil {
		log.Errorf("批量获取评论表情统计失败: %v", err)
		return nil, err
	}

	// 转换为评论表情统计格式
	commentSummaryMap := make(map[int][]model.CommentReactionSummary)
	for commentId, summaries := range summaryMap {
		commentSummaries := make([]model.CommentReactionSummary, len(summaries))
		for i, summary := range summaries {
			commentSummaries[i] = model.CommentReactionSummary{
				ReactionType: summary.ReactionType,
				Count:        summary.Count,
				UserReacted:  summary.UserReacted,
				Users:        summary.Users,
			}
		}
		commentSummaryMap[commentId] = commentSummaries
	}

	log.Infof("批量获取表情统计成功，结果数量: %d", len(commentSummaryMap))
	for commentId, reactions := range commentSummaryMap {
		log.Infof("评论ID %d 有 %d 个表情统计", commentId, len(reactions))
	}

	return commentSummaryMap, nil
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

// ToggleReaction 切换表情回复状态（使用通用系统）
func (s *CommentReactionService) ToggleReaction(commentId, userId int, reactionType string) (bool, error) {
	log.Infof("切换表情回复状态，评论ID: %d, 用户ID: %d, 表情类型: %s", commentId, userId, reactionType)
	
	// 使用通用表情回复系统
	reactionService := NewReactionService(s.ctx)
	added, err := reactionService.ToggleReaction(model.BusinessTypeComment, commentId, userId, reactionType)
	if err != nil {
		log.Errorf("切换评论表情回复失败: %v", err)
		return false, err
	}
	
	if added {
		log.Infof("添加表情回复成功")
	} else {
		log.Infof("移除表情回复成功")
	}
	
	return added, nil
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
