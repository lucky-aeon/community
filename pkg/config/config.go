package config

// 读取配置信息并且给各个配置类赋值

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
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
	Endpoint  string `yaml:"endpoint"`
	Callback  string `json:"callback"`
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
	appConfig := &AppConfig{}

	file, err := os.ReadFile("cmd/community/config.yaml")
	if err != nil {
		panic(err.Error)
	}
	configData := os.ExpandEnv(string(file))
	err = yaml.Unmarshal([]byte(configData), &appConfig)
	if err != nil {
		panic(err.Error())
	}

	instance = appConfig

}
