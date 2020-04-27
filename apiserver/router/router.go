package router

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zhj0811/fabric/apiserver/handler"
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

	factory := router.Group("/factocy/v1")
	{
		factory.POST("/data", handler.SaveData)
		factory.GET("data", handler.QueryData)
	}

	block := router.Group("/block") // 区块链网络操作
	{
		block.HEAD("/keepalive", handler.KeepaliveQuery) // 探活查询
		block.GET("/blockheight", handler.BlockHeight)   // 区块高度
		block.GET("/kafkaNumber", handler.KafkaNumber)   // kafka数量

		//block log interface
		block.PUT("/loglevel", handler.SetLogLevel)
		//block.GET("/loglevel", handler.GetLogLevel)
	}

	return router
}
