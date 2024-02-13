package services

import "xhyovo.cn/community/server/model"

type FileService struct{}

func (*FileService) Save(file *model.Files) {
	fileDao.Save(file)
}

func (*FileService) Delete(userId, fileId, tenantId int) {
	fileDao.Delete(userId, fileId, tenantId)
}

func (*FileService) Deletes(userId, businessId, tenantId int) {

}
