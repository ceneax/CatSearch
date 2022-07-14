package route

import (
	"CatSearch/config"
	"CatSearch/model"
	"CatSearch/search"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func StartSearcher() {
	r := gin.Default()
	r.Static("/static", "static")
	r.LoadHTMLGlob("template/*")

	r.GET("/", func(context *gin.Context) {
		context.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/search", func(context *gin.Context) {
		p := model.Params{
			Query: context.Query("q"),
			Page:  context.Query("page"),
		}

		if len(p.Query) <= 0 {
			context.Redirect(http.StatusFound, "/")
			return
		}

		search.Search(context, p)
	})

	if config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	err := r.Run(config.ServerIp + ":" + config.ServerPort)
	if err != nil {
		log.Fatal("搜索引擎服务器启动失败！" + err.Error())
	}
}
