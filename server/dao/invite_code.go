package dao

import "xhyovo.cn/community/server/model"

type InviteCode struct {
}

func (*InviteCode) SaveCode(code int) {

	model.InviteCode().Create(&model.InviteCodes{Code: code, State: false})

}

// 是否存在code
func (*InviteCode) Exist(code int) bool {

	var count int64
	object := &model.InviteCodes{}
	model.InviteCode().Where("code = ? and state = 0", code).Find(object).Count(&count)

	return count == 1
}

func (*InviteCode) Del(code int) {

	model.InviteCode().Where("code = ?", code).Delete(&model.InviteCodes{})
}

func (*InviteCode) SetState(id int) {

	model.InviteCode().Where("code = ?", id).Update("state", 1)
}

func (c *InviteCode) GetCount() int64 {
	var count int64
	model.InviteCode().Count(&count)
	return count
}

func (c *InviteCode) PageCodes(page int, limit int) []*model.InviteCodes {
	var codes []*model.InviteCodes
	model.InviteCode().Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&codes)
	return codes
}

func (c *InviteCode) SaveCodes(codeList []*model.InviteCodes) {
	model.InviteCode().Create(&codeList)
}
