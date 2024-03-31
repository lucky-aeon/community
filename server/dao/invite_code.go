package dao

import "xhyovo.cn/community/server/model"

type InviteCode struct {
}

// 是否存在code
func (*InviteCode) Exist(code string) bool {

	var count int64
	object := &model.InviteCodes{}
	model.InviteCode().Where("code = ?", code).Find(object).Count(&count)

	return count == 1
}

func (*InviteCode) Del(code int) int64 {

	tx := model.InviteCode().Where("code = ? and state = 0", code).Delete(&model.InviteCodes{})
	return tx.RowsAffected
}

func (*InviteCode) SetState(code string) {

	model.InviteCode().Where("code = ?", code).Update("state", 1)
}

func (c *InviteCode) GetCount() int64 {
	var count int64
	model.InviteCode().Count(&count)
	return count
}

func (c *InviteCode) PageCodes(page int, limit int) []*model.InviteCodes {
	var codes []*model.InviteCodes
	model.InviteCode().Limit(limit).Offset((page - 1) * limit).Order("id desc").Find(&codes)
	return codes
}

func (c *InviteCode) SaveCodes(codeList []*model.InviteCodes) {
	model.InviteCode().Create(&codeList)
}
