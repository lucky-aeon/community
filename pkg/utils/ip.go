package utils

import (
	"github.com/gin-gonic/gin"
)

func GetClientIP(ctx *gin.Context) string {

	return ctx.GetHeader("X-Real-IP")

}
