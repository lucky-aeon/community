package dao

import (
	"xhyovo.cn/community/server/model"
)

type File struct {
}

func (*File) Save(file *model.File) {
	db.Create(file)
}

func (*File) GetFileInfo(fileId, tenantId uint) *model.File {
	fileInfo := &model.File{}
	db.Where(&model.File{ID: fileId, TenantId: tenantId}).Find(fileInfo)

	return fileInfo
}

func (*File) Delete(userId, fileId, tenantId uint) {
	db.Where("id = ? and user_id = ? and tenant_id = ?", fileId, userId, tenantId).Delete(&model.File{})
}

func (*File) Deletes(userId, businessId, tenantId uint) {
	db.Where("business_id = ? and user_id = ? and tenant_id = ?", businessId, userId, tenantId).Delete(&model.File{})

}

func (*File) GetFileKeys(businessId uint) []string {
	var files []model.File
	db.Where("business_id = ?", businessId).Find(files)
	var fileKeys []string
	for _, file := range files {
		fileKeys = append(fileKeys, file.FileKey)
	}
	return fileKeys
}
