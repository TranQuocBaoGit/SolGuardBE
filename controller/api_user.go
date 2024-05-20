package controller

import (
	"context"
	"getContractDeployment/configs"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData User
		err := c.ShouldBind(&formData)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		config, err := configs.LoadConfig(".")
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		client := configs.ConnectDB(config)
		ctx := context.Background()
		col := getCollection(client, config, "user")

		err = AddNewUser(ctx, col,  formData.WalletAddress)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		responsesReturn(c, http.StatusOK, "successful", formData.WalletAddress)
	}
}