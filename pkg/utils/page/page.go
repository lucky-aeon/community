package page

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type Page struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}

func GetPage(ctx *gin.Context) (p, limit int) {
	p, _ = strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ = strconv.Atoi(ctx.DefaultQuery("limit", "15"))
	return p, limit
}

func New(data interface{}, count int64) *Page {
	return &Page{Data: data, Count: count}
}
