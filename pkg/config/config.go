package config

// 读取配置信息并且给各个配置类赋值

import (
	"os"
	"strconv"
)

type AppConfig struct {
	ServerBind  string      `yaml:"server-bind" default:":8080"`
	DbConfig    DbConfig    `yaml:"db"`
	OssConfig   OssConfig   `yaml:"oss"`
	EmailConfig EmailConfig `yaml:"email"`
}

type DbConfig struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type OssConfig struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Bucket    string `yaml:"bucket"`
	Cdn       string `yaml:"cdn"`
	Callback  string `json:"callback"`
	Endpoint  string `yaml:"endpoint"`
}

type EmailConfig struct {
	Address   string `yaml:"address"`
	PollCount int    `yaml:"pollCount"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Host      string `yaml:"host"`
}

var instance *AppConfig

func GetInstance() *AppConfig {
	return instance
}

func Init() {
	pollCount, _ := strconv.Atoi(os.Getenv("POLLCOUNT"))
	if pollCount == 0 {
		pollCount = 10
	}
	appConfig := &AppConfig{
		DbConfig: DbConfig{
			Address:  os.Getenv("DB_HOST"),
			Database: os.Getenv("DB_DATABASE"),
			Username: os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASS"),
		},
		OssConfig: OssConfig{
			AccessKey: os.Getenv("OSS_ACCESS_KEY"),
			Bucket:    os.Getenv("OSS_BUCKET"),
			SecretKey: os.Getenv("OSS_SECRET_KEY"),
			Cdn:       os.Getenv("OSS_CDN"),
			Callback:  os.Getenv("OSS_CALLBACK"),
			Endpoint:  os.Getenv("OSS_ENDPOINT"),
		},
		EmailConfig: EmailConfig{
			Address:   os.Getenv("ADDRESS"),
			Username:  os.Getenv("USERNAME"),
			Password:  os.Getenv("PASSWORD"),
			Host:      os.Getenv("HOST"),
			PollCount: pollCount,
		},
	}
	instance = appConfig

}
