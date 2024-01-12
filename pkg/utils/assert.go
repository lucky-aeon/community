package utils

import (
	"errors"
	"fmt"
)

func NotBlank(any ...interface{}) error {
	for _, v := range any {
		if v == nil || v == "" {
			return errors.New(fmt.Sprintf("%s,is null or is empty", v))
		}
	}
	return nil
}
