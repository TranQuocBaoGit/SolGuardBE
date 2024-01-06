package controller

import (
	"encoding/json"
	"getContractDeployment/configs"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ContractAddress struct {
	Address string `form:"address"`
	ChainID int    `form:"chainid"`
}

func GetContractSourceCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		var contractAddr ContractAddress
		err := c.ShouldBind(&contractAddr)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		
		dataReturn, err := getContractSourceCode(contractAddr)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		
		responsesReturn(c, http.StatusOK, "successful", dataReturn)
	}
}

func getContractSourceCode(contractAddr ContractAddress)(map[string]interface{}, error){
	config, err := configs.LoadConfig(".")
	if err != nil{
		return nil, err
	}

	data := getContractSourceCodeEthscan(contractAddr, config)
	return contractCodeDataHandler(data, false)
}

// stupid address = 0xDeFB0B264032e4e128b00D02b3FD0aA00331237b
// not so stupid address = 0xdAC17F958D2ee523a2206206994597C13D831ec7

func contractCodeDataHandler(data interface{}, write bool) (map[string]interface{}, error){
	var response models.EtherscanSourceResponses

	dataInByte, err := json.Marshal(data)
	if err != nil{
		return nil, err
	}
	json.Unmarshal(dataInByte, &response)

	result := response.Result[0]
	source := result.SourceCode

	finalReturn := make(map[string]interface{})
	if (source[0] == '{'){
		var allFileMap map[string]interface{}
		source = source[1:len(source)-1]

		err := json.Unmarshal([]byte(source), &allFileMap)
		if err != nil{
			return nil, err
		}

		allFile := allFileMap["sources"].(map[string]interface{})

		for contract, contentInterface := range allFile {
			content := contentInterface.(map[string]interface{})["content"].(string)

			if strings.HasPrefix(contract, "@openzeppelin"){
				contract = "openzeppelin/" + helper.GetLastFilePath(contract)
				content = helper.ReplacePathWithFilename(content)
			}else {
				contract = helper.GetLastFilePath(contract)
			}

			if write{
				writePath := "./result/contracts/" + contract
				helper.WriteFile(content, writePath)
			}

			finalReturn[contract] = content
		}
	} else {
		contract :=  result.ContractName + ".sol"
		finalReturn[contract] = source
	}

	finalReturn["Main Contract"] = result.ContractName + ".sol"

	return finalReturn, nil
}