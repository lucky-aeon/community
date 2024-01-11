package services

import (
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/dao"
)

func GenerateCode() int {

	var resCode int
	// generate code
	for resCode == 0 {
		code := utils.GenerateCode(8)
		resCode = dao.InviteCode.QueryCode(code)
	}

	return resCode
}

func DestoryCode(code int) error {

	return dao.InviteCode.Del(code)
}
