package config

// 读取配置信息并且给各个配置类赋值

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	DbConfig   DbConfig   `yaml:"db"`
	KodoConfig KodoConfig `yaml:"kodo"`
}

type DbConfig struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type KodoConfig struct {
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
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
