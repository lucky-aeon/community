package dao

import (
	"xhyovo.cn/community/server/model"
)

type CommentSummaryDao struct{}

var CommentSummaryDaoInstance = &CommentSummaryDao{}

// GetByBusinessAndTenant 根据业务ID和租户ID获取总结
func (d *CommentSummaryDao) GetByBusinessAndTenant(businessId, tenantId int) (*model.CommentSummary, error) {
	var summary model.CommentSummary
	err := model.CommentSummaryModel().
		Where("business_id = ? AND tenant_id = ?", businessId, tenantId).
		First(&summary).Error
	return &summary, err
}

// CreateOrUpdate 创建或更新总结
func (d *CommentSummaryDao) CreateOrUpdate(summary *model.CommentSummary) error {
	var existing model.CommentSummary
	err := model.CommentSummaryModel().
		Where("business_id = ? AND tenant_id = ?", summary.BusinessId, summary.TenantId).
		First(&existing).Error
	
	if err != nil {
		// 记录不存在，创建新记录
		return model.CommentSummaryModel().Create(summary).Error
	} else {
		// 记录存在，更新现有记录
		return model.CommentSummaryModel().
			Where("id = ?", existing.ID).
			Updates(map[string]interface{}{
				"summary":         summary.Summary,
				"comment_count":   summary.CommentCount,
				"last_comment_id": summary.LastCommentId,
			}).Error
	}
}

// NeedUpdate 判断是否需要更新总结
func (d *CommentSummaryDao) NeedUpdate(businessId, tenantId int) (bool, int, error) {
	var summary model.CommentSummary
	err := model.CommentSummaryModel().
		Where("business_id = ? AND tenant_id = ?", businessId, tenantId).
		First(&summary).Error
	
	if err != nil {
		// 记录不存在，需要创建
		return true, 0, nil
	}
	
	// 获取当前评论总数
	var currentCount int64
	model.Comment().Where("business_id = ? AND tenant_id = ?", businessId, tenantId).Count(&currentCount)
	
	// 如果评论数量增加了3条或以上，则需要更新
	needUpdate := int(currentCount) > summary.CommentCount+2
	return needUpdate, summary.LastCommentId, nil
}