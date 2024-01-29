package result

import "github.com/gin-gonic/gin"

type R struct {
	Code int
	Data interface{}
	Msg  string
	ok   bool
}

func Ok(data any, msg string) *R {
	return &R{Code: 200, Data: data, Msg: msg, ok: true}
}
func Err(msg string) *R {
	return &R{Code: 500, Data: nil, Msg: msg}
}
func Auto(data any, err error) *R {
	if err == nil {
		return Ok(data, "成功")
	}
	return Err(err.Error())
}
func (t *R) OkMsg(msg string) *R {

	if t.ok {
		t.Msg = msg
	}
	return t
}

func (t *R) ErrMsg(msg string) *R {
	if !t.ok {
		t.Msg = msg
	}
	return t
}

func (r *R) Json(c *gin.Context) {
	c.JSON(r.Code, r)
}

func (r *R) Xml(c *gin.Context) {
	c.XML(r.Code, r)
}
