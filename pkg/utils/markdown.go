package utils

import (
	"regexp"
	"strings"
)

func GetFirstImage(markdown string) string {
	// 使用正则表达式匹配 Markdown 中的图片链接
	re := regexp.MustCompile(`\[.*?]\((.*?)\)`)
	matches := re.FindStringSubmatch(markdown)

	// 如果匹配到图片链接，则返回第一个匹配结果
	if len(matches) > 1 {
		// 提取链接中的 fileKey 参数值
		params := strings.Split(matches[1], "?")
		for _, param := range params {
			if strings.HasPrefix(param, "fileKey=") {
				return strings.TrimPrefix(param, "fileKey=")
			}
		}
	}
	// 否则返回 null
	return ""
}
