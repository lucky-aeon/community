package services

import (
	"encoding/json"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/dustin/go-humanize"
	"xhyovo.cn/community/pkg/log"
	"xhyovo.cn/community/server/model"
)

type FileService struct{}

func (*FileService) Save(file *model.Files) {
	file.TenantId = 1
	jsonBody, _ := json.Marshal(file)
	log.Infof("用户id: %d,上传文件: %s", file.UserId, jsonBody)
	fileDao.Save(file)
}

func (*FileService) Delete(userId, fileId, tenantId int) {
	fileDao.Delete(userId, fileId, tenantId)
}

func (*FileService) Deletes(userId, businessId, tenantcdId int) {

}

func (s *FileService) PageFiles(p, limit, userId int) (files []model.Files, count int64) {

	files = fileDao.PageFiles(p, limit, userId)
	if len(files) == 0 {
		return []model.Files{}, 0
	}
	var uS UserService
	userIds := mapset.NewSet[int]()
	for i := range files {
		userIds.Add(files[i].UserId)
	}
	nameMap := uS.ListByIdsToMap(userIds.ToSlice())
	for i := range files {
		files[i].UserName = nameMap[files[i].UserId].Name
		files[i].SizeName = humanize.Bytes(uint64(files[i].Size))
	}
	count = fileDao.Count()

	return files, count
}

func (s *FileService) ExistFile(key string) bool {
	var count int64
	model.File().Where("file_key = ?", key).Count(&count)
	return count == 1
}
