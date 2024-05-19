package routes

import (
	"getContractDeployment/configs"
	"getContractDeployment/controller"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc{
	return func(c *gin.Context){
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		c.Next()
	}
}

func Route(config configs.Config) {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/api/source_code", controller.GetContractSourceCode())
	router.GET("/api/bytecode", controller.GetContractByteCode())
	router.POST("/api/analysis/file", controller.PostAnalysisByFile())
	router.POST("/api/analysis/address", controller.PostAnalysisByAddress())
	

	router.Run(config.SERVER)
}