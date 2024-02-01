package utils

import (
	"github.com/google/uuid"
	"strconv"
	"xhyovo.cn/community/pkg/kodo"
)

func BuildFileKey(userId int) string {

	str := strconv.Itoa(userId)

	return str + "/" + (uuid.New().String())
}

func BuildFileUrl(fileKey string) string {
	return "http://" + kodo.GetDomain() + "/" + fileKey
}
