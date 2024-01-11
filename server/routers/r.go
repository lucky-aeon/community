package routers

import "github.com/gin-gonic/gin"

var R r

type r struct {
	code int
	data interface{}
	msg  string
}

func (r *r) Ok() *r {
	r.code = 200
	return r
}

func (r *r) Error() *r {
	r.code = 500
	return r
}

func (r *r) Data(any interface{}) *r {
	r.data = any
	return r
}

func (r *r) Msg(msg string) *r {
	r.msg = msg
	return r
}

func (r *r) Res(c *gin.Context) {

	c.PureJSON(r.code, r)

}
