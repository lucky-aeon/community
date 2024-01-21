package dao

import "xhyovo.cn/community/server/model"

type InviteCode struct {
}

func (*InviteCode) SaveCode(code uint16) {

	model.InviteCode().Create(&model.InviteCodes{Code: code, State: false})

}

// 是否存在code
func (*InviteCode) Exist(code uint16) bool {

	var count int64
	object := &model.InviteCodes{}
	model.InviteCode().Where("code = ? and state = 0", code).Find(object).Count(&count)

	return count == 1
}

func (*InviteCode) Del(id int) {

	model.InviteCode().Where("id = ?", id).Delete(&model.InviteCodes{})
}

func (*InviteCode) SetState(id uint16) {

	model.InviteCode().Model(&model.InviteCodes{}).Where("code = ?", id).Update("state", 1)

}
