package main

import (
	"xhyovo.cn/community/server/config"
	"xhyovo.cn/community/server/routers"
)

func main() {

	config.InitConfig()

	routers.InitRouter()

}
