package utils

import "github.com/gin-gonic/gin"

type R struct {
	code int
	data interface{}
	msg  string
}

func Ok() *R {

	return &R{code: 200}
}

func Error() *R {
	return &R{code: 500}
}

func (r *R) Data(any interface{}) *R {
	r.data = any
	return r
}

func (r *R) Msg(msg string) *R {
	r.msg = msg
	return r
}

func (r *R) Res(c *gin.Context) {

	c.PureJSON(r.code, r)

}
