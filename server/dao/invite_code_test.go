package dao

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
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

func initConfig() {
	appConfig := &AppConfig{}

	file, err := os.ReadFile("D:\\go_project\\community\\config.yaml")
	if err != nil {

		panic(err.Error)

	}
	err = yaml.Unmarshal(file, &appConfig)
	if err != nil {
		panic(err.Error())
	}
	db := appConfig.DbConfig
	Init(db.Username, db.Password, db.Address, db.Database)

}

func TestMain(m *testing.M) {

	initConfig()

	m.Run()
}

func TestExist(t *testing.T) {
	var i inviteCode
	fmt.Println(i.Exist(1))
}
