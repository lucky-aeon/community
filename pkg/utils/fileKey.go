package utils

import (
	"strconv"

	"github.com/google/uuid"
	"xhyovo.cn/community/pkg/kodo"
)

func BuildFileKey(userId int) string {

	str := strconv.Itoa(userId)

	return str + "/" + (uuid.New().String())
}

func BuildFileUrl(fileKey string) string {
	return "http://" + kodo.GetDomain() + "/" + fileKey
}
