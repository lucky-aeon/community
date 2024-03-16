package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"time"
	xt "xhyovo.cn/community/pkg/time"
	"xhyovo.cn/community/pkg/utils"
	"xhyovo.cn/community/server/model"
	services "xhyovo.cn/community/server/service"
)

var log services.LogServices

// 自定义一个结构体，实现 gin.ResponseWriter interface
type responseWriter struct {
	gin.ResponseWriter
	b *bytes.Buffer
}

// 重写 Write([]byte) (int, error) 方法
func (w responseWriter) Write(b []byte) (int, error) {
	//向一个bytes.buffer中写一份数据来为获取body使用
	w.b.Write(b)
	//完成gin.Context.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func OperLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer
		reqBytes, _ := c.GetRawData()

		// 请求包体写回。
		if len(reqBytes) > 0 {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(reqBytes))
		}
		c.Next()

		// 请求后
		latency := time.Since(t)
		execTime := latency
		body := writer.b.String()
		if len(body) > 200 {
			body = body[:200]
		}

		logs := model.OperLogs{
			ExecAt:        execTime.String(),
			RequestMethod: c.Request.Method,
			RequestInfo:   c.Request.URL.Path,
			RequestBody:   string(reqBytes),
			UserId:        GetUserId(c),
			Ip:            utils.GetClientIP(c.Request),
			ResponseData:  body,
			CreatedAt:     xt.Now(),
		}

		log.InsertOperLog(logs)
	}
}
