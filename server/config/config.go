package config

import (
	"os"
	"xhyovo.cn/community/server/dao"

	"gopkg.in/yaml.v3"
)

type AppConfig struct {
	DbConfig DbConfig `yaml:"db"`
}

type DbConfig struct {
	Address  string `yaml:"address"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func InitConfig() {
	config := &AppConfig{}

	file, err := os.ReadFile("D:\\go_project\\community\\config.yaml")
	if err != nil {

		panic(err.Error)

	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err.Error())
	}
	dao.InitDb(&config.DbConfig)

}
