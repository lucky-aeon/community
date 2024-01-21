package services

import (
	"xhyovo.cn/community/pkg/kodo"
	"xhyovo.cn/community/server/model"
)

type FileService struct{}

// save file
func (*FileService) Save(userId, businessId uint, fileKey string) error {

	// get fileInfo
	fileInfo, err := kodo.GetFileInfo(fileKey)
	if err != nil {
		return err
	}

	f := model.Files{
		FileKey:    fileKey,
		Size:       fileInfo.Fsize,
		Format:     fileInfo.MimeType,
		BusinessId: businessId,
		UserId:     userId,
	}

	file.Save(&f)

	return nil

}

func (*FileService) Delete(userId, fileId, tenantId uint) {

	file.Delete(userId, fileId, tenantId)

	// ignore err because this is not important
	kodo.Delete(file.GetFileInfo(fileId, tenantId).FileKey)

}

func (*FileService) Deletes(userId, businessId, tenantId uint) {

	// 获取所有文件的key
	fileKeys := file.GetFileKeys(businessId)
	file.Deletes(userId, businessId, tenantId)
	for _, fileKey := range fileKeys {
		go kodo.Delete(fileKey)
	}
}

// file get
func (*FileService) GetFileKey(fileId, tenantId uint) string {
	return file.GetFileInfo(fileId, tenantId).FileKey

}
