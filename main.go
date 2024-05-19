package main

import (
	"getContractDeployment/configs"
	"getContractDeployment/helper"
	"getContractDeployment/routes"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)



func main() {
	config, err := configs.LoadConfig(".")
	helper.CheckError(err)

	// // Disable Console Color, you don't need console color when writing the logs to file.
	// gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	routes.Route(config)

	// file := "example.sol"
	// contractFolder := "contracts"
	// remappingJSON := false


	// docker.Test(file, contractFolder, remappingJSON)

}