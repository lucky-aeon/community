package utils

import (
	"strings"
)

func IsImage(fileName string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, ext := range extensions {
		if strings.HasSuffix(strings.ToLower(fileName), ext) {
			return true
		}
	}

	return false
}
