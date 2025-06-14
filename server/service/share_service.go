package services

import (
	"crypto/md5"
	"errors"
	"fmt"
	"math/rand"
	"time"

	localTime "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
)

type ShareService struct{}

// CreateShareRequest 创建分享请求
type CreateShareRequest struct {
	BusinessType string `json:"business_type"`
	BusinessID   int    `json:"business_id"`
	CreatorID    int    `json:"creator_id"`
	ExpireDays   int    `json:"expire_days,omitempty"` // 过期天数，0表示永久
}

// ShareResponse 分享响应
type ShareResponse struct {
	ShareToken string               `json:"share_token"`
	ShareURL   string               `json:"share_url"`
	ExpireAt   *localTime.LocalTime `json:"expire_at,omitempty"`
}

// ShareStatistics 分享统计
type ShareStatistics struct {
	BusinessID int `json:"business_id"`
	ShareCount int `json:"share_count"`
	TotalViews int `json:"total_views"`
}

// VisitorInfo 访问者信息
type VisitorInfo struct {
	IP        string `json:"ip"`
	UserID    int    `json:"user_id,omitempty"`
	UserAgent string `json:"user_agent,omitempty"`
	Referer   string `json:"referer,omitempty"`
}

// CreateShare 创建分享
func (s *ShareService) CreateShare(req CreateShareRequest) (*ShareResponse, error) {
	// 验证业务类型是否允许分享
	if !IsShareableBusinessType(req.BusinessType) {
		return nil, errors.New("该业务类型不支持分享功能")
	}

	// 检查是否已经创建过分享
	var existingShare model.Share
	err := model.ShareDB().
		Where("business_type = ? AND business_id = ? AND creator_id = ? AND is_active = 1",
			req.BusinessType, req.BusinessID, req.CreatorID).
		First(&existingShare).Error

	if err == nil {
		// 已存在分享，返回现有分享
		return &ShareResponse{
			ShareToken: existingShare.ShareToken,
			ShareURL:   fmt.Sprintf("/community/share/%s", existingShare.ShareToken),
			ExpireAt:   existingShare.ExpireAt,
		}, nil
	}

	// 生成分享token
	shareToken := s.generateShareToken()

	// 计算过期时间
	var expireAt *localTime.LocalTime
	if req.ExpireDays > 0 {
		expireTime := time.Now().AddDate(0, 0, req.ExpireDays)
		expireAt = (*localTime.LocalTime)(&expireTime)
	}

	// 创建分享记录
	share := model.Share{
		BusinessType: req.BusinessType,
		BusinessID:   req.BusinessID,
		ShareToken:   shareToken,
		CreatorID:    req.CreatorID,
		ExpireAt:     expireAt,
		IsActive:     1,
	}

	err = model.ShareDB().Create(&share).Error
	if err != nil {
		return nil, err
	}

	return &ShareResponse{
		ShareToken: shareToken,
		ShareURL:   fmt.Sprintf("/community/share/%s", shareToken),
		ExpireAt:   expireAt,
	}, nil
}

// GetShareByToken 根据token获取分享信息
func (s *ShareService) GetShareByToken(token string) (*model.Share, error) {
	var share model.Share
	err := model.ShareDB().
		Where("share_token = ? AND is_active = 1", token).
		First(&share).Error

	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if share.ExpireAt != nil {
		expireTime := time.Time(*share.ExpireAt)
		if time.Now().After(expireTime) {
			return nil, errors.New("分享链接已过期")
		}
	}

	return &share, nil
}

// RecordView 记录访问
func (s *ShareService) RecordView(shareID int, visitor VisitorInfo) error {
	// 检查是否为重复访问（同一IP 5分钟内不重复计数）
	if s.isDuplicateView(shareID, visitor.IP) {
		return nil
	}

	// 记录访问
	shareView := model.ShareView{
		ShareID:   shareID,
		VisitorIP: visitor.IP,
		VisitorID: visitor.UserID,
		UserAgent: visitor.UserAgent,
		Referer:   visitor.Referer,
	}

	err := model.ShareViewDB().Create(&shareView).Error
	if err != nil {
		return err
	}

	// 更新分享记录的浏览次数
	return model.ShareDB().
		Where("id = ?", shareID).
		UpdateColumn("total_views", model.ShareDB().Raw("total_views + 1")).Error
}

// GetStatisticsByBusinessIDs 批量获取业务对象的分享统计
func (s *ShareService) GetStatisticsByBusinessIDs(businessType string, businessIDs []int) (map[int]ShareStatistics, error) {
	var results []struct {
		BusinessID int `json:"business_id"`
		ShareCount int `json:"share_count"`
		TotalViews int `json:"total_views"`
	}

	err := model.ShareDB().
		Select("business_id, COUNT(*) as share_count, SUM(total_views) as total_views").
		Where("business_type = ? AND business_id IN ? AND is_active = 1", businessType, businessIDs).
		Group("business_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// 转换为map
	statsMap := make(map[int]ShareStatistics)
	for _, result := range results {
		statsMap[result.BusinessID] = ShareStatistics{
			BusinessID: result.BusinessID,
			ShareCount: result.ShareCount,
			TotalViews: result.TotalViews,
		}
	}

	return statsMap, nil
}

// generateShareToken 生成分享token
func (s *ShareService) generateShareToken() string {
	// 使用时间戳和随机数生成token
	timestamp := time.Now().UnixNano()
	random := rand.Int63()

	source := fmt.Sprintf("%d_%d", timestamp, random)
	hash := md5.Sum([]byte(source))

	// 取前16个字符作为token
	return fmt.Sprintf("%x", hash)[:16]
}

// isDuplicateView 检查是否为重复访问
func (s *ShareService) isDuplicateView(shareID int, ip string) bool {
	var count int64
	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	model.ShareViewDB().
		Where("share_id = ? AND visitor_ip = ? AND visited_at > ?", shareID, ip, fiveMinutesAgo).
		Count(&count)

	return count > 0
}
