package dao

import (
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var kodo services.Kodo

type File struct {
}

func (*File) Save(fileKey string) error {
	// get fileInfo
	info, err := kodo.GetFileInfo(fileKey)
	if err != nil {
		return err
	}
	db.Create(&model.File{FileKey: fileKey, Size: info.Fsize, Format: info.MimeType})
	return nil
}

func (*File) GetFileInfo(fileKey string) (*model.File, error) {
	info, err := kodo.GetFileInfo(fileKey)
	if err != nil {
		return nil, err
	}
	fileInfo := &model.File{
		FileKey: fileKey,
		Size:    info.Fsize,
		Format:  info.MimeType,
	}
	return fileInfo, nil
}

func (*File) Del(fileKey string) {
	db.Where("file_key = ?", fileKey).Delete(&model.File{})
}
