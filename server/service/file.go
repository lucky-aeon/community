package services

import (
	"github.com/dustin/go-humanize"
	"xhyovo.cn/community/server/model"
)

type FileService struct{}

func (*FileService) Save(file *model.Files) {
	file.TenantId = 1
	fileDao.Save(file)
}

func (*FileService) Delete(userId, fileId, tenantId int) {
	fileDao.Delete(userId, fileId, tenantId)
}

func (*FileService) Deletes(userId, businessId, tenantId int) {

}

func (s *FileService) PageFiles(p, limit, userId int) (files []model.Files, count int64) {

	files = fileDao.PageFiles(p, limit, userId)
	if len(files) == 0 {
		return []model.Files{}, 0
	}
	var uS UserService
	var uIds = make([]int, 0)
	for i := range files {
		uIds = append(uIds, files[i].UserId)
	}
	nameMap := uS.ListByIdsToMap(uIds)
	for i := range files {
		files[i].UserName = nameMap[files[i].UserId].Name
		files[i].SizeName = humanize.Bytes(uint64(files[i].Size))
	}
	count = fileDao.Count()

	return files, count
}
