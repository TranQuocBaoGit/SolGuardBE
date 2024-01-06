package routes

import (
	"getContractDeployment/configs"
	"getContractDeployment/controller"

	"github.com/gin-gonic/gin"
)

func Route(config configs.Config) {
	server := gin.Default()

	// server.POST("/todo_list", controllers.AddTodo())

	server.GET("/api/source_code", controller.GetContractSourceCode())
	server.GET("/api/bytecode", controller.GetContractByteCode())
	server.POST("/api/analysis/file", controller.PostAnalysisByFile())
	server.POST("/api/analysis/address", controller.PostAnalysisByAddress())

	server.GET("/api/test", controller.GetTest())
	

	server.Run(config.SERVER)
}