package main

import (
	"xhyovo.cn/community/cmd/community/routers"
	"xhyovo.cn/community/server/config"
)

func main() {

	config.InitConfig()

	routers.InitRouter()

}
