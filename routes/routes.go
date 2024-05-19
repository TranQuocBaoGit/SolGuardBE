package routes

import (
	"getContractDeployment/configs"
	"getContractDeployment/controller"

	"github.com/gin-gonic/gin"
)

func Route(config configs.Config) {
	router := gin.Default()

	router.GET("/api/source_code", controller.GetContractSourceCode())
	router.GET("/api/bytecode", controller.GetContractByteCode())
	router.POST("/api/analysis/file", controller.PostAnalysisByFile())
	router.POST("/api/analysis/address", controller.PostAnalysisByAddress())
	

	router.Run(config.SERVER)
}