package controller

import (
	"getContractDeployment/models"

	"github.com/gin-gonic/gin"
)

func responsesReturn(c *gin.Context, status int, message string, data map[string]interface{}){
	res := models.Responses{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.JSON(status, res)
}