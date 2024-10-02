package dao

import (
	"xhyovo.cn/community/pkg/constant"
	model "xhyovo.cn/community/server/model/knowledge"
)

type DocumentDao struct {
}

// Create 创建文档
func (d *DocumentDao) Create(doc *model.Documents) (int, error) {
	if err := model.Document().Create(doc).Error; err != nil {
		return 0, err
	}
	return doc.ID, nil
}

// Get 获取文档
func (d *DocumentDao) Get(id int) (model.Documents, error) {
	var doc model.Documents
	if err := model.Document().First(&doc, id).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

/*
*
返回 err 说明不存在，则不用理会
*/
func (d *DocumentDao) GetByLink(businessId int, typee constant.ContentType) (*model.Documents, error) {
	var doc model.Documents
	if err := model.Document().Where("business_id", businessId).Where("type", typee).First(&doc).Error; err != nil {
		return nil, err
	}
	return &doc, nil
}

func (d *DocumentDao) Exist(businessId int, typee constant.ContentType) bool {
	var count int64
	model.Document().Where("business_id", businessId).Where("type", typee).Count(&count)
	return count == 1
}

// Update 更新文档
func (d *DocumentDao) Update(doc *model.Documents) error {
	if err := model.Document().Save(doc).Error; err != nil {
		return err
	}
	return nil
}

// Delete 删除文档
func (d *DocumentDao) Delete(businessId int, typee constant.ContentType) error {
	if err := model.Document().Where("business_id", businessId).Where("type", typee).Delete(&model.Documents{}).Error; err != nil {
		return err
	}
	return nil
}

// List 获取所有文档
func (d *DocumentDao) List() ([]model.Documents, error) {
	var docs []model.Documents
	if err := model.Document().Find(&docs).Error; err != nil {
		return nil, err
	}
	return docs, nil
}

func (d *DocumentDao) ListById(ids []int) ([]model.Documents, error) {
	var documents []model.Documents
	// 假设 d.DB 是 GORM 的 *gorm.DB 实例
	if err := model.Document().Where("id IN (?)", ids).Find(&documents).Error; err != nil {
		return nil, err
	}
	return documents, nil
}
