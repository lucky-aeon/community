package dao_test

import (
	"gopkg.in/yaml.v3"
	"os"
	"testing"
	"xhyovo.cn/community/server/config"
	"xhyovo.cn/community/server/dao"
)

func initConfig() {
	appConfig := &config.AppConfig{}

	file, err := os.ReadFile("D:\\go_project\\community\\config.yaml")
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

func TestMain(m *testing.M) {

	initConfig()

	m.Run()
}

func TestQueryUserByAccountAndPswd(t *testing.T) {

}
