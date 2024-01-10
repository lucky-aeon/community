package services

import (
	"xhyovo.cn/community/dao"
	"xhyovo.cn/community/model"
	"xhyovo.cn/community/utils"
)

func GenerateCode() int {

	codeObject := &model.InviteCode{}
	// generate code
	for codeObject == nil {
		code := utils.GenerateCode(8)
		codeObject = dao.InviteCode.QuerySingle(&model.InviteCode{Code: code})
	}

	return codeObject.Code
}

func DestoryCode(code int) {

	dao.InviteCode.Del(code)
}
