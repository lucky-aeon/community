package dao

import (
	"xhyovo.cn/community/server/model"
)

type CrawlerConfigDao struct{}

// AddConfig 添加爬虫配置
func (dao *CrawlerConfigDao) AddConfig(config *model.CrawlerConfig) error {
	return model.CrawlerConfigModel().Create(config).Error
}

// GetByID 根据ID获取配置
func (dao *CrawlerConfigDao) GetByID(id int) *model.CrawlerConfig {
	var config model.CrawlerConfig
	model.CrawlerConfigModel().Where("id = ?", id).First(&config)
	return &config
}

// ListAll 获取所有配置
func (dao *CrawlerConfigDao) ListAll() ([]*model.CrawlerConfig, error) {
	var configs []*model.CrawlerConfig
	err := model.CrawlerConfigModel().Order("created_at desc").Find(&configs).Error
	return configs, err
}

// ListByStatus 根据状态获取配置
func (dao *CrawlerConfigDao) ListByStatus(status int) ([]*model.CrawlerConfig, error) {
	var configs []*model.CrawlerConfig
	err := model.CrawlerConfigModel().Where("status = ?", status).Order("created_at desc").Find(&configs).Error
	return configs, err
}

// UpdateStatus 更新配置状态
func (dao *CrawlerConfigDao) UpdateStatus(id, status int) error {
	return model.CrawlerConfigModel().Where("id = ?", id).Update("status", status).Error
}

// Update 更新配置
func (dao *CrawlerConfigDao) Update(config *model.CrawlerConfig) error {
	return model.CrawlerConfigModel().Where("id = ?", config.ID).Updates(config).Error
}

// Delete 删除配置
func (dao *CrawlerConfigDao) Delete(id int) error {
	return model.CrawlerConfigModel().Where("id = ?", id).Delete(&model.CrawlerConfig{}).Error
}
