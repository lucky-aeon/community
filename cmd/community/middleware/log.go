package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"time"
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
	//完成gin.Record.Writer.Write()原有功能
	return w.ResponseWriter.Write(b)
}

func OperLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		//请求体 body
		requestBody := ""
		b, err := c.GetRawData()
		if err != nil {
			requestBody = "failed to get request body"
		} else {
			requestBody = string(b)
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(b))

		writer := responseWriter{
			c.Writer,
			bytes.NewBuffer([]byte{}),
		}
		c.Writer = writer
		c.Next()

		// 请求后
		latency := time.Since(t)
		execTime := latency
		body := writer.b.String()
		if len(body) > 100 {
			body = body[:100]
		}

		logs := model.OperLogs{
			ExecAt:        execTime.String(),
			RequestMethod: c.Request.Method,
			RequestInfo:   c.Request.URL.Path,
			RequestBody:   requestBody,
			UserId:        GetUserId(c),
			Ip:            utils.GetClientIP(c),
			Platform:      c.GetHeader("sec-ch-ua-platform"),
			UserAgent:     c.GetHeader("user-agent"),
			ResponseData:  body,
		}

		log.InsertOperLog(logs)
	}
}
