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
		c.Writer.Header().Set("Content-Type", "text/plain; multipart/form-data; application/json; charset=utf-8")
		// c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, ngrok-skip-browser-warning, application/json")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, ngrok-skip-browser-warning")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

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
	router.POST("/api/user/login", controller.Login())
	router.GET("/api/user/history", controller.GetUserHistory())
	router.POST("/api/user/history/delete", controller.DeleteHistory())
	

	router.Run(config.SERVER)
}


	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{"*"} // Allow requests from your frontend domain
	// corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	// corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "ngrok-skip-browser-warning"}
	// corsConfig.ExposeHeaders = []string{"Content-Length"}
	// corsConfig.AllowCredentials = true
	// corsConfig.MaxAge = 12 * time.Hour

	// // if c.Request.Method == "OPTIONS" {
	// // 	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	// // 	c.Writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	// // 	c.AbortWithStatus(204)
	// // 	return
	// // }
	
	// router.Use(cors.New(corsConfig))