package result

import "github.com/gin-gonic/gin"

type R struct {
	Code int
	Data interface{}
	Msg  string
}

func Ok(data any, msg string) *R {
	return &R{Code: 200, Data: data, Msg: msg}
}
func Err(msg string) *R {
	return &R{Code: 200, Data: nil, Msg: msg}
}
func (r *R) Json(c *gin.Context) {
	c.JSON(r.Code, r)
}
