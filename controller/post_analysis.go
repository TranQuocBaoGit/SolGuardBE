package controller

import (
	"getContractDeployment/configs"
	"getContractDeployment/docker"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PostAnalysisByFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("smc")
		if err != nil && err != http.ErrMissingFile{
			responsesReturn(c, http.StatusInternalServerError, "Error when upload file", nil)
			return
		}

		if err != http.ErrMissingFile{
			// Save smart contract file
			saveFilePath := "./result/contracts/" +  file.Filename
			err =  c.SaveUploadedFile(file, saveFilePath)
			if err != nil && err != http.ErrMissingFile {
				responsesReturn(c, http.StatusInternalServerError, "Error when processing file", nil)
				return
			}
		}
		
		dataReturn, err := returnAnalysisResult(file.Filename, false)
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		responsesReturn(c, http.StatusOK, "success", dataReturn)
	}
}

func PostAnalysisByAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var contractAddr ContractAddress
		err := c.ShouldBind(&contractAddr)
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		contractData, err := getContractSourceCode(contractAddr)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		err = docker.CreateMythrilMappingJson(contractData)
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		dataReturn, err := returnAnalysisResult(contractData["Main Contract"].(string), true)
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		responsesReturn(c, http.StatusOK, "success", dataReturn)
	}
}

func returnAnalysisResult(mainFile string, remapping bool)(map[string]interface{}, error){

	config, err := configs.LoadConfig(".")
	if err != nil {
		return nil, err
	}
	contractPath := config.CONTRACT_PATH

	analysisResult, err := docker.RunMythrilAnalysis(contractPath, mainFile, remapping)
	if err != nil {
		return nil, err
	}

	var dataReturn map[string]interface{} = map[string]interface{}{"result": analysisResult}
	return dataReturn, nil
}