package config

// 读取配置信息并且给各个配置类赋值

import (
	"os"

	"gopkg.in/yaml.v3"
	"xhyovo.cn/community/server/dao"
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
	appConfig := &AppConfig{}

	file, err := os.ReadFile("cmd/community/config.yaml")
	if err != nil {
		panic(err.Error)
	}
	err = yaml.Unmarshal(file, &appConfig)
	if err != nil {
		panic(err.Error())
	}
	db := appConfig.DbConfig
	dao.Init(db.Username, db.Password, db.Address, db.Database)

}
