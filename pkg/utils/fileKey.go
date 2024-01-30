package utils

import (
	"github.com/google/uuid"
	"strconv"
	"xhyovo.cn/community/pkg/kodo"
)

func BuildFileKey(userId uint) string {

	str := strconv.Itoa(int(userId))

	return str + "/" + (uuid.New().String())
}

func BuildFileUrl(fileKey string) string {
	return "http://" + kodo.GetDomain() + "/" + fileKey
}
