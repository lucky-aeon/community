package dao

import (
	"xhyovo.cn/community/server/model"
)

type File struct {
}

func (*File) Save(file *model.Files) {
	model.File().Create(&file)
}

func (*File) GetFileInfo(fileId, tenantId int) *model.Files {
	fileInfo := &model.Files{}
	model.File().Where(&model.Files{ID: fileId, TenantId: tenantId}).Find(fileInfo)

	return fileInfo
}

func (*File) Delete(userId, fileId, tenantId int) {
	model.File().Where("id = ? and user_id = ? and tenant_id = ?", fileId, userId, tenantId).Delete(&model.Files{})
}

func (*File) Deletes(userId, businessId, tenantId int) {
	model.File().Where("business_id = ? and user_id = ? and tenant_id = ?", businessId, userId, tenantId).Delete(&model.Files{})

}

func (*File) GetFileKeys(businessId int) []string {
	var files []model.Files
	model.File().Where("business_id = ?", businessId).Find(files)
	var fileKeys []string
	for _, file := range files {
		fileKeys = append(fileKeys, file.FileKey)
	}
	return fileKeys
}

func (f *File) PageFiles(p, limit int) []model.Files {

	var files []model.Files
	model.File().Offset((p - 1) * limit).Limit(limit).Order("created_at desc").Find(&files)

	return files
}

func (f *File) Count() int64 {

	var count int64
	model.File().Count(&count)

	return count
}
