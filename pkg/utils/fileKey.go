package utils

import (
	"strconv"
	"xhyovo.cn/community/pkg/oss"

	"github.com/google/uuid"
)

func BuildFileKey(userId int) string {

	str := strconv.Itoa(userId)

	return str + "/" + (uuid.New().String())
}

func BuildFileUrl(fileKey string) string {
	return oss.GetEndpoint() + "/" + fileKey
}
