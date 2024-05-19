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
		var formData AddressFormData
		err := c.ShouldBind(&formData)
		helper.CheckError(err)

		baseUrl := ""
		switch formData.ChainID{
		case 1:
			baseUrl = "https://eth.rpc.blxrbdn.com"
		case 5:
			baseUrl = "https://goerli.gateway.tenderly.co"
		case 11155111:
			baseUrl = "https://sepolia.gateway.tenderly.co"
		default:
			baseUrl = "https://eth.rpc.blxrbdn.com"
		}

		client, err := ethclient.Dial(baseUrl)
		helper.CheckError(err)

		data := getSmartContractByteCode(client, formData.Address)

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
