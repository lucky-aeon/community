package services

import (
	"time"

	"xhyovo.cn/community/server/model"
)

type AiNewsService struct{}

// HistoryDateItem 历史日期项
type HistoryDateItem struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

// GetHistoryDates 获取有AI新闻的历史日期列表
func (s *AiNewsService) GetHistoryDates() ([]HistoryDateItem, error) {
	// 先查询所有文章，然后在Go代码中按日期分组
	var articles []*model.AiNews
	err := model.AiNewsModel().
		Select("publish_date").
		Where("deleted_at IS NULL").
		Find(&articles).Error

	if err != nil {
		return nil, err
	}

	// 按日期分组统计
	dateCount := make(map[string]int)
	for _, article := range articles {
		// 将LocalTime转换为标准时间，然后格式化为日期字符串
		publishTime := time.Time(article.PublishDate)
		if !publishTime.IsZero() {
			dateStr := publishTime.Format("2006-01-02")
			dateCount[dateStr]++
		}
	}

	// 转换为结果数组并排序
	var results []HistoryDateItem
	for date, count := range dateCount {
		results = append(results, HistoryDateItem{
			Date:  date,
			Count: count,
		})
	}

	// 按日期降序排序
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Date < results[j].Date {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}

// GetLatestDate 获取最新有数据的日期
func (s *AiNewsService) GetLatestDate() (string, error) {
	var article model.AiNews
	err := model.AiNewsModel().
		Select("publish_date").
		Where("deleted_at IS NULL").
		Order("publish_date DESC").
		First(&article).Error

	if err != nil {
		return "", err
	}

	// 将LocalTime转换为标准时间，然后格式化为日期字符串
	publishTime := time.Time(article.PublishDate)
	return publishTime.Format("2006-01-02"), nil
}

// GetDailyNews 获取指定日期的AI新闻列表
func (s *AiNewsService) GetDailyNews(date string, withContent bool) ([]*model.AiNews, error) {
	var articles []*model.AiNews

	query := model.AiNewsModel().
		Where("DATE(publish_date) = ? AND deleted_at IS NULL", date).
		Order("id DESC")

	// 如果不需要内容，则不查询content字段（优化性能）
	if !withContent {
		query = query.Select("id, title, summary, source_url, source_name, publish_date, category, tags")
	}

	err := query.Find(&articles).Error
	return articles, err
}

// GetNewsById 根据ID获取文章详情
func (s *AiNewsService) GetNewsById(id int) (*model.AiNews, error) {
	var article model.AiNews
	err := model.AiNewsModel().
		Where("id = ? AND deleted_at IS NULL", id).
		First(&article).Error
	return &article, err
}

// GetNewsByDate 获取指定日期范围内的新闻
func (s *AiNewsService) GetNewsByDate(startDate, endDate string, page, limit int) ([]*model.AiNews, int64, error) {
	var articles []*model.AiNews
	var total int64

	query := model.AiNewsModel().
		Where("DATE(publish_date) BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate)

	// 计算总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	err = query.
		Order("publish_date DESC, id DESC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&articles).Error

	return articles, total, err
}
