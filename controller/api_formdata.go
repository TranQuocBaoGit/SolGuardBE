package controller

import (
	"getContractDeployment/models"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type FileAnalysisFormData struct {
	File  *multipart.FileHeader `form:"file"`
	// Tools []string 				`form:"tools"`
}

type User struct {
	WalletAddress string   `form:"wallet_address" json:"wallet_address"`
}

type AddressAnalysisFormData struct {
	WalletAddress string `form:"wallet_address" json:"wallet_address"`
	Address string   `form:"address" json:"address"`
	ChainID string   `form:"chainid" json:"chainid"`
	Dapp 	string 	 `form:"dApp,omitempty" json:"dApp,omitempty"`
	Decision string  `form:"decision,omitempty" json:"decision,omitempty"`
}

type AddressFormData struct {
	Address string `form:"address" json:"address"`
	ChainID int    `form:"chainid" json:"chainid"`
}

type DeleteHistoryFormData struct {
	WalletAddress 		string `form:"wallet_address" json:"wallet_address"`
	HistoryAnalyzeID 	string `form:"history_analyze_id" json:"history_analyze_id"`
}

func responsesReturn(c *gin.Context, status int, message string, data interface{}){
	res := models.Responses{
		Status:  status,
		Message: message,
		Data:    data,
	}
	c.JSON(status, res)
}