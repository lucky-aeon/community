package routers

import (
	"github.com/gin-gonic/gin"
	"log"
)

// init router
func InitRouter() {
	r := gin.Default()

	InitLoginRegisterRouter(r)
	InitFileRouter(r)
	err := r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
