package services

import "xhyovo.cn/community/server/model"

type Log struct {
}

func (*Log) GetPageOperLog(page, limit int, begin, end string) (logs []model.OperLogs, count int64) {
	db := model.OperLog()
	db.Limit(limit).Where("created_at >= ? and created_at <= ?", begin, end).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)
	db.Count(&count)
	return
}

func (*Log) InsertOperLog(log model.OperLogs) {
	model.OperLog().Create(&log)
}

func (*Log) DeletesOperLogs(id []int) {
	model.OperLog().Where("id ? in").Delete(model.OperLogs{})
}
