package config

import (
	"fmt"
	"os"

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
	config := AppConfig{}

	file, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic(err.Error)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(config)
}
