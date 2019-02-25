package router

import (
	"net/http"
	"sync"

	"github.com/peersafe/factoring/apiserver/handler"

	"github.com/gin-gonic/gin"
)

// Router 全局路由
var router *gin.Engine
var onceCreateRouter sync.Once

func GetRouter() *gin.Engine {
	onceCreateRouter.Do(func() {
		router = createRouter()
	})

	return router
}

func createRouter() *gin.Engine {
	router := gin.Default()
	// 版本控制
	// v1 := Router.Group("/v1")
	// {
	// factor
	factor := router.Group("/trade")
	{
		factor.POST("/savedata", handler.SaveData)
		factor.GET("/tradeinfo", handler.QueryByInfo)

		factor.GET("/blockheight", handler.BlockHeight)
		factor.HEAD("/keepalive", handler.KeepaliveQuery)
		factor.GET("/keepalive", handler.Keepalive)
		factor.GET("/kafkaNumber", handler.KafkaNumber)
		factor.GET("/version", handler.Version)

		//log interface
		factor.POST("/setLogLevel", handler.SetLogLevel)
		factor.GET("/getLogLevel", handler.GetLogLevel)

		//BlockQuery and BlockQueryEx is dropped, which are used in couchDB environment.
		//factor.GET("/block/:id", handler.BlockQuery)
		//factor.GET("/blockQuery/:id", handler.BlockQueryEx)
	}
	//upload schema json file
	router.StaticFS("/schema", http.Dir("./schema"))
	return router
}
