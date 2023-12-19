package main

import (
	"github.com/afyi/sytalk/router"
	"github.com/gin-gonic/gin"
)

func main() {

	// 调试模式
	gin.SetMode(gin.DebugMode)

	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default() // create router using gin

	router.Init(r)

	r.Run(":8000") // register router to port 8000

}
