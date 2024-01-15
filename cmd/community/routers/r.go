package routers

import "github.com/gin-gonic/gin"

var R r

type r struct {
	Code int
	Data interface{}
	Msg  string
}

func (*r) Ok() *r {

	return &r{Code: 200}
}

func (*r) Error() *r {
	return &r{Code: 500}
}

func (r *r) setData(any interface{}) *r {
	r.Data = any
	return r
}

func (r *r) setMsg(msg string) *r {
	r.Msg = msg
	return r
}

func (r *r) Res(c *gin.Context) {

	c.JSONP(r.Code, r)

}
