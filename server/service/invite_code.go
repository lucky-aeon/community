package services

import (
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/dao"
)

var InviteCode dao.InviteCode

func GenerateCode() uint16 {

	var flag bool
	var code uint16
	// generate code
	for flag {
		code = utils.GenerateCode(8)
		flag = InviteCode.Exist(code)
	}

	return code
}

func DestroyCode(code int) {

	InviteCode.Del(code)
}

func SetState(id uint16) {
	InviteCode.SetState(id)
}
