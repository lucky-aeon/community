package services

import (
	"xhyovo.cn/community/server/dao"
	"xhyovo.cn/community/server/model"
)

type File struct{}

var file dao.File

var k Kodo

// save file
func (*File) Save(userId, businessId uint, fileKey string) error {

	// get fileInfo
	fileInfo, err := k.GetFileInfo(fileKey)
	if err != nil {
		return err
	}

	f := model.File{
		FileKey:    fileKey,
		Size:       fileInfo.Fsize,
		Format:     fileInfo.MimeType,
		BusinessId: businessId,
		UserId:     userId,
	}

	file.Save(&f)

	return nil

}

func (*File) Delete(userId, fileId, tenantId uint) {

	file.Delete(userId, fileId, tenantId)

	// ignore err because this is not important
	k.Delete(file.GetFileInfo(fileId, tenantId).FileKey)

}

func (*File) Deletes(userId, businessId, tenantId uint) {

	// 获取所有文件的key
	fileKeys := file.GetFileKeys(businessId)
	file.Deletes(userId, businessId, tenantId)
	for _, fileKey := range fileKeys {
		go k.Delete(fileKey)
	}
}

// file get
func (*File) GetFileKey(fileId, tenantId uint) string {
	return file.GetFileInfo(fileId, tenantId).FileKey

}
