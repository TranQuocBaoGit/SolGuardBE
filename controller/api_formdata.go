package controller

import (
	"getContractDeployment/models"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type FileAnalysisFormData struct {
	File  *multipart.FileHeader `form:"file"`
	Tools []string 				`form:"tools"`
}

type AddressAnalysisFormData struct {
	Address string   `form:"address"`
	ChainID int      `form:"chainid"`
	SumUp 	bool	 `form:"sumup"`
	Tools   []string `form:"tools"`
}

type AddressFormData struct {
	Address string `form:"address"`
	ChainID int    `form:"chainid"`
}

func responsesReturn(c *gin.Context, status int, message string, data interface{}){
	res := models.Responses{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.JSON(status, res)
}