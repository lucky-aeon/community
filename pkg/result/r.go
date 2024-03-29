package result

import "github.com/gin-gonic/gin"

type R struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Ok   bool        `json:"ok"`
}

func OkWithMsg(data any, msg string) *R {
	return &R{Code: 2000, Data: data, Msg: msg, Ok: true}
}

func Ok(data any, msg string) *R {
	return &R{Code: 200, Data: data, Msg: msg, Ok: true}
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

func Page(data any, total int64, err error) *R {
	if err != nil {
		return Err(err.Error())
	}
	return Ok(map[string]any{
		"list":  data,
		"total": total,
	}, "成功")
}

func (t *R) OkMsg(msg string) *R {

	if t.Ok {
		t.Msg = msg
	}
	return t
}

func (t *R) ErrMsg(msg string) *R {
	if !t.Ok {
		t.Msg = msg
	}
	return t
}

func (r *R) Json(c *gin.Context) {

	c.JSON(200, r)
}

func (r *R) Xml(c *gin.Context) {
	c.XML(r.Code, r)
}
