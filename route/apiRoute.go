package route

import (
	"CatSearch/config"
	"CatSearch/controller"
	"CatSearch/middleware"
	"github.com/gin-gonic/gin"
	"log"
)

func StartApi() {
	r := gin.Default()

	r.Use(middleware.Cors())

	r.POST("/translate", controller.Translate)
	r.GET("/codeAnswer", controller.CodeAnswer)

	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	err := r.Run(config.ServerIp + ":" + config.ServerApiPort)
	if err != nil {
		log.Fatal("Api服务器启动失败！" + err.Error())
	}
}
