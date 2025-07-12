package services

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
	llmService "xhyovo.cn/community/server/service/llm"
)

type CommentSummaryService struct {
	ctx *gin.Context
}

func NewCommentSummaryService(ctx *gin.Context) *CommentSummaryService {
	return &CommentSummaryService{ctx: ctx}
}

// GetSummary 获取评论总结
func (s *CommentSummaryService) GetSummary(businessId, tenantId int) (*model.CommentSummary, error) {
	summaryDao := dao.CommentSummaryDaoInstance

	// 检查是否需要更新
	needUpdate, lastCommentId, err := summaryDao.NeedUpdate(businessId, tenantId)
	if err != nil {
		log.Errorf("检查总结更新状态失败: %v", err)
	}

	// 如果需要更新，异步执行更新
	if needUpdate {
		go func() {
			if err := s.generateSummary(businessId, tenantId, lastCommentId); err != nil {
				log.Errorf("异步生成总结失败: %v", err)
			}
		}()
	}

	// 返回现有总结（如果存在）
	summary, err := summaryDao.GetByBusinessAndTenant(businessId, tenantId)
	if err != nil {
		// 如果没有总结，检查是否有评论
		comments := s.getCommentsForSummary(businessId, tenantId)
		if len(comments) == 0 {
			// 没有评论，返回空结果而不是错误
			return nil, nil
		}
		
		// 有评论但没有总结，尝试立即生成一个
		if err := s.generateSummary(businessId, tenantId, 0); err != nil {
			return nil, fmt.Errorf("生成总结失败: %v", err)
		}
		// 重新获取刚生成的总结
		summary, err = summaryDao.GetByBusinessAndTenant(businessId, tenantId)
	}

	return summary, err
}

// UpdateSummaryIfNeeded 检查并更新总结（用于评论发布后调用）
func (s *CommentSummaryService) UpdateSummaryIfNeeded(businessId, tenantId int) {
	summaryDao := dao.CommentSummaryDaoInstance
	needUpdate, lastCommentId, err := summaryDao.NeedUpdate(businessId, tenantId)
	if err != nil {
		log.Errorf("检查总结更新状态失败: %v", err)
		return
	}

	if needUpdate {
		if err := s.generateSummary(businessId, tenantId, lastCommentId); err != nil {
			log.Errorf("更新总结失败: %v", err)
		}
	}
}

// generateSummary 生成总结
func (s *CommentSummaryService) generateSummary(businessId, tenantId, lastCommentId int) error {
	// 获取所有评论内容
	comments := s.getCommentsForSummary(businessId, tenantId)
	if len(comments) == 0 {
		return fmt.Errorf("没有找到评论内容")
	}

	// 构建评论文本
	commentsText := s.buildCommentsText(comments)

	// 调用LLM生成总结
	summary, err := s.callLLMForSummary(commentsText, tenantId)
	if err != nil {
		return fmt.Errorf("LLM调用失败: %v", err)
	}

	// 保存总结
	commentSummary := &model.CommentSummary{
		BusinessId:    businessId,
		TenantId:      tenantId,
		Summary:       summary,
		CommentCount:  len(comments),
		LastCommentId: s.getMaxCommentId(comments),
	}

	summaryDao := dao.CommentSummaryDaoInstance
	return summaryDao.CreateOrUpdate(commentSummary)
}

// getCommentsForSummary 获取用于总结的评论
func (s *CommentSummaryService) getCommentsForSummary(businessId, tenantId int) []*model.Comments {
	var comments []*model.Comments
	model.Comment().
		Where("business_id = ? AND tenant_id = ?", businessId, tenantId).
		Order("created_at ASC").
		Find(&comments)
	return comments
}

// buildCommentsText 构建评论文本
func (s *CommentSummaryService) buildCommentsText(comments []*model.Comments) string {
	var builder strings.Builder

	for i, comment := range comments {
		builder.WriteString(fmt.Sprintf("评论%d: %s\n", i+1, comment.Content))
	}

	return builder.String()
}

// callLLMForSummary 调用LLM生成总结
func (s *CommentSummaryService) callLLMForSummary(commentsText string, tenantId int) (string, error) {
	llm := &llmService.LLMService{}

	// 根据租户类型确定业务类型名称
	businessType := s.getBusinessTypeName(tenantId)

	systemPrompt := fmt.Sprintf(`你是一个专业的内容分析师，专门帮助用户快速理解%s评论区的核心内容。

用户面临的问题：评论太多，难以快速找到有用信息。

请按以下结构整理评论内容：

**📋 核心问题与解答**
- 提取用户提出的主要问题
- 标注对应的回答和解决方案
- 格式：Q: [问题] A: [答案/解决方案]

**💡 有价值的信息**
- 实用技巧和建议
- 经验分享和最佳实践
- 工具推荐和资源链接

**🔥 热门讨论点**
- 被多次提及的话题
- 有争议但有建设性的讨论
- 需要注意的问题和坑点

**📝 补充信息**
- 其他有用的补充说明
- 相关的扩展讨论

要求：
- 重点提取问题和对应的答案
- 突出实用性和可操作性
- 如果没有明显的Q&A，则重点提取有价值的信息点
- 保持结构清晰，便于快速阅读
- 总结长度控制在300-600字
- **不要引用具体评论序号**（如"评论1"、"评论2"等），直接描述内容即可
- 用中文回复`, businessType)

	userPrompt := fmt.Sprintf("请总结以下评论内容：\n\n%s", commentsText)

	return llm.Chat(systemPrompt, userPrompt)
}

// getBusinessTypeName 获取业务类型名称
func (s *CommentSummaryService) getBusinessTypeName(tenantId int) string {
	switch tenantId {
	case 0:
		return "文章"
	case 1:
		return "课程章节"
	case 2:
		return "课程"
	case 3:
		return "分享会"
	case 4:
		return "AI日报"
	default:
		return "内容"
	}
}

// getMaxCommentId 获取评论中的最大ID
func (s *CommentSummaryService) getMaxCommentId(comments []*model.Comments) int {
	maxId := 0
	for _, comment := range comments {
		if comment.ID > maxId {
			maxId = comment.ID
		}
	}
	return maxId
}
