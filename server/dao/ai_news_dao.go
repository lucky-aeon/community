package dao

import (
	"crypto/md5"
	"fmt"

	"xhyovo.cn/community/server/model"
)

type AiNewsDao struct{}

// AddNews 添加AI新闻
func (dao *AiNewsDao) AddNews(news *model.AiNews) error {
	// 生成内容哈希用于去重
	news.Hash = dao.GenerateHash(news.Title + news.Content)
	return model.AiNewsModel().Create(news).Error
}

// GetByHash 根据哈希值查询新闻（用于去重）
func (dao *AiNewsDao) GetByHash(hash string) *model.AiNews {
	var news model.AiNews
	model.AiNewsModel().Where("hash = ?", hash).First(&news)
	return &news
}

// ExistsByHash 检查哈希值是否已存在
func (dao *AiNewsDao) ExistsByHash(hash string) bool {
	var count int64
	model.AiNewsModel().Where("hash = ?", hash).Count(&count)
	return count > 0
}

// UpdateStatus 更新新闻状态
func (dao *AiNewsDao) UpdateStatus(id, status int) error {
	return model.AiNewsModel().Where("id = ?", id).Update("status", status).Error
}

// GetByID 根据ID获取新闻
func (dao *AiNewsDao) GetByID(id int) *model.AiNews {
	var news model.AiNews
	model.AiNewsModel().Where("id = ?", id).First(&news)
	return &news
}

// ListByStatus 根据状态分页查询新闻
func (dao *AiNewsDao) ListByStatus(status, page, limit int) ([]*model.AiNews, int64) {
	var news []*model.AiNews
	var count int64

	db := model.AiNewsModel().Where("status = ?", status)
	db.Count(&count)

	if count == 0 {
		return news, 0
	}

	db.Order("created_at desc").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&news)

	return news, count
}

// ListAll 分页查询所有新闻（管理端用）
func (dao *AiNewsDao) ListAll(page, limit int) ([]*model.AiNews, int64) {
	var news []*model.AiNews
	var count int64

	db := model.AiNewsModel()
	db.Count(&count)

	if count == 0 {
		return news, 0
	}

	db.Order("created_at desc").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&news)

	return news, count
}

// ListByDateRange 根据日期范围查询新闻
func (dao *AiNewsDao) ListByDateRange(startDate, endDate string, status int) ([]*model.AiNews, error) {
	var news []*model.AiNews

	db := model.AiNewsModel().Where("status = ?", status)
	if startDate != "" && endDate != "" {
		db = db.Where("DATE(created_at) BETWEEN ? AND ?", startDate, endDate)
	}

	err := db.Order("created_at desc").Find(&news).Error
	return news, err
}

// UpdateAITags 更新AI标注的分类和标签
func (dao *AiNewsDao) UpdateAITags(id int, category, tags string) error {
	return model.AiNewsModel().Where("id = ?", id).
		Updates(map[string]interface{}{
			"category": category,
			"tags":     tags,
		}).Error
}

// Delete 删除新闻
func (dao *AiNewsDao) Delete(id int) error {
	return model.AiNewsModel().Where("id = ?", id).Delete(&model.AiNews{}).Error
}

// Create 创建AI新闻（别名方法）
func (dao *AiNewsDao) Create(news *model.AiNews) error {
	return dao.AddNews(news)
}

// ExistsByURL 检查URL是否已存在
func (dao *AiNewsDao) ExistsByURL(url string) (bool, error) {
	var count int64
	err := model.AiNewsModel().Where("source_url = ?", url).Count(&count).Error
	return count > 0, err
}

// GetMaxSourceID 获取指定来源的最大ID
func (dao *AiNewsDao) GetMaxSourceID(sourceName string) (int, error) {
	// 从source_url中提取ID，假设URL格式为 https://www.aibase.com/zh/daily/{id}
	var result struct {
		MaxID int `gorm:"column:max_id"`
	}

	err := model.AiNewsModel().
		Select("COALESCE(MAX(CAST(SUBSTRING_INDEX(source_url, '/', -1) AS UNSIGNED)), 0) as max_id").
		Where("source_name = ? AND source_url REGEXP '[0-9]+$'", sourceName).
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.MaxID, nil
}

// GenerateHash 生成内容哈希（公共方法）
func (dao *AiNewsDao) GenerateHash(content string) string {
	hash := md5.Sum([]byte(content))
	return fmt.Sprintf("%x", hash)
}
