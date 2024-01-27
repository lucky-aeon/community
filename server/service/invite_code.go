package services

import (
	"xhyovo.cn/community/pkg/utils"
)

func GenerateCode() uint16 {

	var flag bool
	var code1 uint16
	// generate code
	for flag {
		code1 = utils.GenerateCode(8)
		flag = code.Exist(code1)
	}

	return code1
}

func DestroyCode(code1 int) {

	code.Del(code1)
}

func SetState(id uint16) {
	code.SetState(id)
}
