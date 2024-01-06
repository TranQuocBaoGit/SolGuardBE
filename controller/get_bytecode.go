package controller

import (
	"getContractDeployment/helper"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

// type contractAddress struct {
// 	Address string `form:"address"`
// }


func GetContractByteCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		var smcAddress ContractAddress
		err := c.ShouldBind(&smcAddress)
		helper.CheckError(err)

		client, err := ethclient.Dial("https://sepolia.gateway.tenderly.co")
		helper.CheckError(err)

		data := getSmartContractByteCode(client, smcAddress.Address)

		dataReturn := make(map[string]interface{})
		dataReturn["bytecode"] = data
		
		responsesReturn(c, http.StatusOK, "successful", dataReturn)
	}
}

type cache struct {
	mu sync.Mutex
	counter int
	
}

var counter int
var mu sync.Mutex

func GetTest() gin.HandlerFunc {
	return func(c *gin.Context) {

		mu.Lock()
		defer mu.Unlock()

		counter++
		dataReturn := make(map[string]interface{})
		dataReturn["result"] = counter

		responsesReturn(c, http.StatusOK, "success", dataReturn)
	}
}
