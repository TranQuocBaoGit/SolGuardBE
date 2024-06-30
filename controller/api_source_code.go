package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"getContractDeployment/configs"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetContractSourceCode() gin.HandlerFunc {
	return func(c *gin.Context) {
		// var formData AddressFormData
		// err := c.ShouldBind(&formData)
		address := c.Query("address")
		chainidStr := c.Query("chainid")
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code API) Chain ID: ", chainidStr), "log.txt")
		chainid, err :=  strconv.Atoi(chainidStr)
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code API) Get source from: ", address), "log.txt")
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		dataReturn, err := getContractSourceCode(chainid, address)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		responsesReturn(c, http.StatusOK, "successful", dataReturn)
	}
}

func getContractSourceCode(chainID int, address string) (models.Contract, error) {

	var result models.Contract

	config, err := configs.LoadConfig(".")
	if err != nil {
		return models.Contract{}, err
	}
	client := configs.ConnectDB(config)
	ctx, _ := context.WithTimeout(context.Background(), 3600*time.Second)
	col := getCollection(client, config, "smc")

	result, err = getSourceCodeViaAddressFromDB(ctx, col, address)
	if err == nil {
		return result, nil
	} else if err != mongo.ErrNoDocuments {
		return models.Contract{}, err
	}

	data := getContractSourceCodeEthscan(chainID, address)
	result, err = contractCodeDataHandler(data, address, chainID)
	if err != nil {
		return models.Contract{}, err
	}

	err = saveSourceCodeToDB(ctx, col, &result)
	if err != nil {
		return models.Contract{}, err
	}

	return result, err
}

// stupid address = 0xDeFB0B264032e4e128b00D02b3FD0aA00331237b
// not so stupid address = 0xdAC17F958D2ee523a2206206994597C13D831ec7


func contractCodeDataHandler(data interface{}, address string, chainID int) (models.Contract, error) {
	var response models.EtherscanSourceResponses

	dataInByte, err := json.Marshal(data)
	if err != nil {
		return models.Contract{}, err
	}
	json.Unmarshal(dataInByte, &response)

	result := response.Result[0]
	source := result.SourceCode

	var finalReturn models.Contract

	mainContract := result.ContractName + ".sol"
	var allContent []models.ContractContent

	if source[0] == '{' {
		var allFileMap map[string]interface{}
		source = source[1 : len(source)-1]

		err := json.Unmarshal([]byte(source), &allFileMap)
		if err != nil {
			return models.Contract{}, err
		}

		allFile := allFileMap["sources"].(map[string]interface{})

		for contract, contentInterface := range allFile {
			content := contentInterface.(map[string]interface{})["content"].(string)

			if strings.HasPrefix(contract, "@openzeppelin") {
				contract = "openzeppelin/" + helper.GetLastFilePath(contract)
				content = helper.ReplacePathWithFilename(content)
			} else {
				contract = helper.GetLastFilePath(contract)
			}
			var oneContent models.ContractContent = models.ContractContent{
				ContractName: contract,
				ContractContent: content,
			}
			allContent = append(allContent, oneContent)
		}
	} else {
		contract := result.ContractName + ".sol"
		var oneContent models.ContractContent = models.ContractContent{
			ContractName: contract,
			ContractContent: source,
		}
		allContent = append(allContent, oneContent)
	}

	finalReturn = models.Contract{
		Address:      address,
		ChainID:      chainID,
		NoContract:   len(allContent),
		MainContract: mainContract,
		Content:      allContent,
	}

	return finalReturn, nil
}

// func contractCodeDataHandler(data interface{}, write bool) (map[string]interface{}, error) {
// 	var response models.EtherscanSourceResponses

// 	dataInByte, err := json.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	json.Unmarshal(dataInByte, &response)

// 	result := response.Result[0]
// 	source := result.SourceCode

// 	finalReturn := make(map[string]interface{})
// 	finalReturn["main_contract"] = result.ContractName + ".sol"
// 	finalReturn["content"] = make(map[string]interface{})
// 	if source[0] == '{' {
// 		var allFileMap map[string]interface{}
// 		source = source[1 : len(source)-1]

// 		err := json.Unmarshal([]byte(source), &allFileMap)
// 		if err != nil {
// 			return nil, err
// 		}

// 		allFile := allFileMap["sources"].(map[string]interface{})

// 		for contract, contentInterface := range allFile {
// 			content := contentInterface.(map[string]interface{})["content"].(string)

// 			if strings.HasPrefix(contract, "@openzeppelin") {
// 				contract = "openzeppelin/" + helper.GetLastFilePath(contract)
// 				content = helper.ReplacePathWithFilename(content)
// 			} else {
// 				contract = helper.GetLastFilePath(contract)
// 			}

// 			if write {
// 				writePath := "./result/contracts/" + contract
// 				helper.WriteFile(content, writePath)
// 			}

// 			finalReturn["content"].(map[string]interface{})[contract] = content
// 		}
// 	} else {
// 		contract := result.ContractName + ".sol"
// 		finalReturn["content"].(map[string]interface{})[contract] = source
// 	}

// 	return finalReturn, nil
// }