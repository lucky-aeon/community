package services

import (
	"xhyovo.cn/community/pkg/kodo"
	"xhyovo.cn/community/server/model"
)

type FileService struct{}

// save fileDao
func (*FileService) Save(userId, businessId int, fileKey string) error {

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

	fileDao.Save(&f)

	return nil

}

func (*FileService) Delete(userId, fileId, tenantId int) {

	fileDao.Delete(userId, fileId, tenantId)

	// ignore err because this is not important
	kodo.Delete(fileDao.GetFileInfo(fileId, tenantId).FileKey)

}

func (*FileService) Deletes(userId, businessId, tenantId int) {

	// 获取所有文件的key
	fileKeys := fileDao.GetFileKeys(businessId)
	fileDao.Deletes(userId, businessId, tenantId)
	for _, fileKey := range fileKeys {
		go kodo.Delete(fileKey)
	}
}

// fileDao get
func (*FileService) GetFileKey(fileId, tenantId int) string {
	return fileDao.GetFileInfo(fileId, tenantId).FileKey

}
