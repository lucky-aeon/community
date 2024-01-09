package utils

import (
	"fmt"
	"math/rand"
)

func GenerateCode(digits int) string {

	code := ""
	for i := 0; i < digits; i++ {
		digit := rand.Int()

		code += fmt.Sprint(digit)
	}

	return code
}
