package controller

import (
	"encoding/json"
	"fmt"
	"getContractDeployment/configs"
	"getContractDeployment/helper"
	"sync"
)


func getContractSourceCodeEthscan(contractAddr ContractAddress, config configs.Config) interface{}{
	var baseUrl string
	switch contractAddr.ChainID{
	case 1:
		baseUrl = "https://api.etherscan.io/"
	case 5:
		baseUrl = "https://api-goerli.etherscan.io/"
	case 11155111:
		baseUrl = "https://api-sepolia.etherscan.io/"
	default:
		baseUrl = "https://api.etherscan.io/"
	}

	url  := fmt.Sprintf("%sapi?module=%s&action=%s&address=%s&apikey=%s", baseUrl, "contract", "getsourcecode", contractAddr.Address, config.ETHER_SCAN_API)
	
	bodyChan := make(chan []byte)
	var wg sync.WaitGroup
	wg.Add(1)
	go fetchAPIData(&wg, url ,bodyChan)
	body := <- bodyChan
	wg.Wait()

	var data interface{}
	err := json.Unmarshal(body, &data)
	helper.CheckError(err)

	return data
}