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

func (s *FileService) PageFiles(p, limit, userId int) (files []model.Files, count int64) {

	files = fileDao.PageFiles(p, limit, userId)
	count = fileDao.Count()
	return files, count
}
