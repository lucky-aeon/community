package services

import (
	"xhyovo.cn/community/pkg/utils"
)

func GenerateCode() uint16 {

	var flag bool
	var code1 uint16
	// generate codeDao
	for flag {
		code1 = utils.GenerateCode(8)
		flag = codeDao.Exist(code1)
	}

	return code1
}

func DestroyCode(code1 int) {

	codeDao.Del(code1)
}

func SetState(id uint16) {
	codeDao.SetState(id)
}
