package log

import (
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var log *zap.SugaredLogger

/*
setJSONEncoder 设置logger编码
*/
func setJSONEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   //转换编码的时间戳
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder //编码级别调整为大写的级别输出
	return zapcore.NewConsoleEncoder(encoderConfig)
}

/*
setLoggerWrite 设置logger写入文件
*/
func setLoggerWrite() zapcore.WriteSyncer {
	l := &lumberjack.Logger{
		Filename:   "./test1.log", //Filename 是要写入日志的文件。
		MaxSize:    1,             //MaxSize 是日志文件在轮换之前的最大大小（以兆字节为单位）。它默认为 100 兆字节
		MaxBackups: 1,             //MaxBackups 是要保留的最大旧日志文件数。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除。）
		MaxAge:     30,            //MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
		Compress:   true,          //压缩
		LocalTime:  true,          //LocalTime 确定用于格式化备份文件中的时间戳的时间是否是计算机的本地时间。默认是使用 UTC 时间。
	}
	return zapcore.AddSync(l)
}

func Init() {
	core := zapcore.NewCore(setJSONEncoder(), zapcore.NewMultiWriteSyncer(setLoggerWrite()), zap.InfoLevel)
	log = zap.New(core, zap.AddCaller()).Sugar()
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

func Info(args ...interface{}) {
	log.Info(args)
}
func Infof(template string, args ...interface{}) {
	log.Infof(template, args)
}
func Logln(lvl zapcore.Level, args ...interface{}) {
	log.Logln(lvl, args)
}
func Warn(args ...interface{}) {
	log.Warn(args)
}
func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args)
}

func Warnln(args ...interface{}) {
	log.Warnln(args)
}

func Error(args ...interface{}) {
	log.Error(args)
}

func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args)
}

func Errorln(args ...interface{}) {
	log.Errorln(args)
}
