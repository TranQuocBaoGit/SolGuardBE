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

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	helper.WriteFile("BEGIN SERVER", "log.txt")

	routes.Route(config)

}