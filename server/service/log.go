package services

import (
	mapset "github.com/deckarep/golang-set/v2"
	"xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/server/model"
)

type LogServices struct {
}

func (*LogServices) GetPageOperLog(page, limit int, logSearch model.LogSearch, flag bool) (logs []model.OperLogs, count int64) {
	db := model.OperLog()
	if logSearch.RequestMethod != "" {
		db.Where("request_method = ?", logSearch.RequestMethod)
	}
	if logSearch.RequestInfo != "" {
		db.Where("request_info like ?", "%"+logSearch.RequestInfo+"%")
	}
	if logSearch.Ip != "" {
		db.Where("ip like ?", "%"+logSearch.Ip+"%")
	}
	if logSearch.StartTime != "" {
		db.Where("created_at <= ? and ? >= created_at", logSearch.StartTime, logSearch.EndTime)
	}
	if logSearch.UserName != "" {
		var userS UserService
		ids := userS.SearchNameSelectId(logSearch.UserName)
		db.Where("user_id in ?", ids)
	}
	if flag {
		db.Where("request_info != '/community/file/singUrl'")
	} else {
		db.Where("request_info = '/community/file/singUrl'")
	}
	db.Count(&count)
	if count == 0 {
		return
	}
	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)

	set := mapset.NewSet[int]()

	for _, lo := range logs {
		set.Add(lo.UserId)
	}
	var u UserService
	userMap := u.ListByIdsToMap(set.ToSlice())
	for i := range logs {
		logs[i].UserName = userMap[logs[i].UserId].Name
	}
	return
}

func (*LogServices) InsertOperLog(log model.OperLogs) {
	go func(log model.OperLogs) {
		model.OperLog().Create(&log)
	}(log)
}

func (*LogServices) InsertLoginLog(log model.LoginLogs) {
	go func(log model.LoginLogs) {
		model.LoginLog().Create(&log)
	}(log)
}

func (s *LogServices) GetPageLoginPage(page, limit int, logSearch model.LogSearch) (logs []model.LoginLogs, count int64) {
	db := model.LoginLog()
	if logSearch.Account != "" {
		db.Where("account like ?", "%"+logSearch.Account+"%")
	}
	if logSearch.Ip != "" {
		db.Where("ip like ?", "%"+logSearch.Ip+"%")
	}

	db.Count(&count)
	if count == 0 {
		return
	}
	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)
	return
}

func (s *LogServices) GetPageQuestionLogs(page int, limit int) ([]struct {
	ID        int            `json:"id"`
	Question  string         `json:"question"`
	UserName  string         `json:"userName"`
	CreatedAt time.LocalTime `json:"createdAt"`
}, int64) {
	var count int64
	var logs []struct {
		ID        int            `json:"id"`
		Question  string         `json:"question"`
		UserName  string         `json:"userName"`
		CreatedAt time.LocalTime `json:"createdAt"`
	}

	db := model.OperLog()

	// 查找 request_info 以 "/community/knowledge/query" 开头的记录，并提取 question 参数、关联的 user name 和提问时间
	db.Select("oper_logs.id, SUBSTRING_INDEX(SUBSTRING(oper_logs.request_info, LOCATE('?question=', oper_logs.request_info) + 10), '&', 1) AS question, users.name AS user_name, oper_logs.created_at").
		Joins("JOIN users ON oper_logs.user_id = users.id").
		Where("oper_logs.request_info LIKE ?", "/community/knowledge/query%").
		Count(&count)

	if count == 0 {
		return nil, 0
	}

	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)

	return logs, count
}
