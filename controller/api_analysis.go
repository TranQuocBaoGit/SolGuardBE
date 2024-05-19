package controller

import (
	"context"
	"fmt"
	"getContractDeployment/configs"
	"getContractDeployment/docker"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func PostAnalysisByFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData FileAnalysisFormData
		err := c.ShouldBind(&formData)
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		contractFolder, err := helper.CreateNewContractsFolder()
		if err != nil {
			helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		file := formData.File
		if err != http.ErrMissingFile{
			// Save smart contract file
			saveFilePath := fmt.Sprintf("./result/%s/%s",contractFolder, file.Filename) // "./result/contracts/"
			err =  c.SaveUploadedFile(file, saveFilePath)
			if err != nil && err != http.ErrMissingFile {
				fmt.Print("missing file")
				helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
				responsesReturn(c, http.StatusInternalServerError, "Error when processing file", nil)
				return
			}
		}
		
		dataReturn, err := returnFullResult(file.Filename, contractFolder, -1, false)
		if err != nil{
			fmt.Print("error analysis")
			helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
		responsesReturn(c, http.StatusOK, "success", dataReturn)
	}
}

func PostAnalysisByAddress() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData AddressAnalysisFormData
		err := c.ShouldBind(&formData)
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		contractData, err := getContractSourceCode(formData.ChainID, formData.Address)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		contractFolder, err := helper.CreateNewContractsFolder()
		if err != nil {
			helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		helper.WriteFile(fmt.Sprint("create ",contractFolder), "log.txt")

		remap := false
		// if len(contractData.Content) >= 2{
		// 	err = docker.CreateMythrilMappingJson(contractData, contractFolder)
		// 	if err != nil{
		// 		helper.DeleteContractsFolder(contractFolder)
		// 		responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
		// 		return
		// 	}
		// 	remap = true
		// }

		dataReturn, err := getAnalysisResult(contractData, contractFolder, remap)
		if err != nil{
			helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
		responsesReturn(c, http.StatusOK, "success", dataReturn)
	}
}

func getAnalysisResult(contractData models.Contract, contractFolder string, remapping bool)(models.Result, error){

	helper.WriteFileExtra(fmt.Sprint("begin analysis"), "log.txt")

	if contractData.ContractID != -1{
		var result models.Result
		config, err := configs.LoadConfig(".")
		if err != nil {
			return models.Result{}, err
		}
		client := configs.ConnectDB(config)
		ctx, _ := context.WithTimeout(context.Background(), 3600*time.Second)
		col := getCollection(client, config, "analysis")
	
		result, err = getAnalysisFromDB(ctx, col, contractData.ContractID)
		if err == nil {
			return result, nil
		} else if err != mongo.ErrNoDocuments {
			return models.Result{}, err
		}
		helper.WriteFileExtra(fmt.Sprint("pass get analysis from DB"), "log.txt")
		
		for _, file := range contractData.Content{
			if file.ContractName == contractData.MainContract{
				helper.WriteFile(file.ContractContent, filepath.Join("result", contractFolder, file.ContractName))
			}
		}

		fullResult, err := returnFullResult(contractData.MainContract, contractFolder, contractData.ContractID, remapping)
		if err != nil {
			return models.Result{}, err
		}
		err = saveAnalysisToDB(ctx, col, &fullResult)
		helper.WriteFileExtra(fmt.Sprint("save analysis success"), "log.txt")
		if err != nil {
			return models.Result{}, err
		}

		return fullResult, nil
	}

	return returnFullResult(contractData.MainContract, contractFolder, contractData.ContractID, remapping)

}

func returnFullResult(mainFile string, contractFolder string, contractID int, remapping bool) (models.Result, error){

	var toolsResult []models.ToolResult

	// Start analysis timer
	mythrilStart := time.Now()
	
	// Run container
	mythrilDetail, err := docker.RunMythrilAnalysis(mainFile, contractFolder, remapping)

	// End analysis timer
	mythrilEnd := time.Since(mythrilStart)

	// Get sum up result
	mythrilSumUp := docker.GetMythrilSumUp(mythrilDetail)


	var mythril models.ToolResult
	mythril.ToolName = "mythril"
	mythril.NoError = len(mythrilDetail.Issues)
	mythril.SumUps = mythrilSumUp
	mythril.Detail = mythrilDetail
	mythril.TimeElapsed = mythrilEnd.Seconds()
	toolsResult = append(toolsResult, mythril)

	// var wg sync.WaitGroup
	// wg.Add(2)

	// outputChan := make(chan models.ToolResult)

	// go func(){
	// 	defer wg.Done()



	// if err != nil {
	// 	helper.WriteFileExtra(err.Error(), "log.txt")
	// } else {

	// }


	// toolsResult = append(toolsResult, mythril)
	// outputChan <- mythril
	// totalTimeElapsed += mythrilEnd.Seconds()
	helper.WriteFileExtra(fmt.Sprint("pass mythril"), "log.txt")
	// }()
	// mythrilStart := time.Now()
	// mythrilDetail, err := docker.RunMythrilAnalysis(mainFile, contractFolder, remapping)
	// if err != nil{
	// 	return models.Result{}, err
	// }
	// mythrilSumUp := docker.GetMythrilSumUp(mythrilDetail)
	// mythrilEnd := time.Since(mythrilStart)

	// var mythril models.ToolResult
	// mythril.ToolName = "mythril"
	// mythril.NoError = len(mythrilDetail.Issues)
	// mythril.SumUps = mythrilSumUp
	// mythril.Detail = mythrilDetail
	// mythril.TimeElapsed = mythrilEnd.Seconds()
	// toolsResult = append(toolsResult, mythril)
	// totalTimeElapsed += mythrilEnd.Seconds()
	// helper.WriteFileExtra(fmt.Sprint("pass mythril"), "log.txt")

	slitherStart := time.Now()
	slitherDetail, err := docker.RunSlitherAnalysis(mainFile, contractFolder, remapping)
	if err != nil{
		helper.WriteFileExtra(err.Error(), "log.txt")
		return models.Result{}, err
	}
	slitherSumUp := docker.GetSlitherSumUp(slitherDetail)
	slitherEnd := time.Since(slitherStart)

	var slither models.ToolResult
	slither.ToolName = "slither"
	slither.SumUps = slitherSumUp
	slither.NoError = len(slitherSumUp)
	slither.Detail = slitherDetail
	slither.TimeElapsed = slitherEnd.Seconds()
	toolsResult = append(toolsResult, slither)
	// toolsResult = append(toolsResult, slither)
	// outputChan <- slither
	// totalTimeElapsed += slitherEnd.Seconds()
	helper.WriteFileExtra(fmt.Sprint("pass slither"), "log.txt")

	// Close the output channel once both Goroutines are done.
	// go func() {
	// 	wg.Wait()
	// 	close(outputChan)
	// }()
	
	// for output := range outputChan {
	// 	toolsResult = append(toolsResult, output)
	// }

	return models.Result{
		ContractID: contractID,
		ToolsResult: toolsResult,
	}, nil
}


func StandardizeResult(toolsResult []models.ToolResult) models.StandardizeResults {

	// Map vulnerability to its highest severity
	vulneToServerity := make(map[string]string)
	for _, toolResult := range toolsResult{
		for _, sumUp := range toolResult.SumUps{
			vulne, severity := IdentifyVulnerability(sumUp, toolResult.ToolName)
			_, exist := vulneToServerity[vulne]
			if !exist {
				vulneToServerity[vulne] = severity
			} else {
				newSeverity := GetHighestSeverity(vulneToServerity[vulne], severity)
				vulneToServerity[vulne] = newSeverity
			}
		}
	} 

	// Get standardize result
	standardizeResult := models.StandardizeResults{
		NoError: 0,
		Result: []models.StandardizeResult{},
	}
	for vulnerability, severity := range vulneToServerity {
		standardizeResult.NoError++
		standardizeResult.Result = append(standardizeResult.Result, models.StandardizeResult{
			Name: vulnerability,
			Severity: severity,
		})
	}

	return standardizeResult
}

func GetHighestSeverity(sever1, sever2 string) string {
	if sever1 == "High" || sever2 == "High" {return "High"}
	if sever1 == "Medium" || sever2 == "Medium" {return "Medium"}
	return "Low"
}


func IdentifyVulnerability(sumup models.SumUp, tool string) (string, string) {
	if tool == "mythril"{
		return docker.MythrilStandardize(sumup)
	} else if tool == "slither"{
		return "", ""
	}
	return  "", ""
}