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
	
		result, err = getAnalysisViaContractFromDB(ctx, col, contractData.ContractID)
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
	start := time.Now()

	mythrilStart := time.Now()
	mythrilDetail, err := docker.RunMythrilAnalysisWithTimeOut(mainFile, contractFolder, remapping)
	if err != nil{
		helper.WriteFileExtra(err.Error(), "log.txt")
		// return models.Result{}, err
	}
	mythrilEnd := time.Since(mythrilStart)
	mythrilSumUp := docker.GetMythrilSumUp(mythrilDetail, err)
	var mythril models.ToolResult
	mythril.ToolName = "mythril"
	mythril.NoError = len(mythrilDetail.Issues)
	mythril.SumUps = mythrilSumUp
	mythril.Detail = mythrilDetail
	mythril.TimeElapsed = mythrilEnd.Seconds()
	toolsResult = append(toolsResult, mythril)
	helper.WriteFileExtra(fmt.Sprint("pass mythril"), "log.txt")



	slitherStart := time.Now()
	slitherDetail, err := docker.RunSlitherAnalysisWithTimeOut(mainFile, contractFolder, remapping)
	if err != nil{
		helper.WriteFileExtra(err.Error(), "log.txt")
		// return models.Result{}, err
	}
	slitherSumUp := docker.GetSlitherSumUp(slitherDetail, err)
	slitherEnd := time.Since(slitherStart)
	var slither models.ToolResult
	slither.ToolName = "slither"
	slither.SumUps = slitherSumUp
	if slitherSumUp[0].Name != "SLITHER ERROR" {
		slither.NoError = len(slitherSumUp)
	} else {
		slither.NoError = 0
	}
	slither.Detail = slitherDetail
	slither.TimeElapsed = slitherEnd.Seconds()
	toolsResult = append(toolsResult, slither)
	helper.WriteFileExtra(fmt.Sprint("pass slither"), "log.txt")



	solhintStart := time.Now()
	solhintDetail, err := docker.RunSolHintAnalysisWithTimeOut(mainFile, contractFolder, remapping)
	if err != nil{
		helper.WriteFileExtra(err.Error(), "log.txt")
		// return models.Result{}, err
	}
	solhintSumUp := docker.GetSolhintSumUp(solhintDetail, err)
	solhintEnd := time.Since(solhintStart)
	var solhint models.ToolResult
	solhint.ToolName = "solhint"
	solhint.SumUps = solhintSumUp
	if solhintSumUp[0].Name != "SOLHINT ERROR" {
		solhint.NoError = len(solhintSumUp)
	} else {
		solhint.NoError = 0
	}
	solhint.Detail = solhintDetail
	solhint.TimeElapsed = solhintEnd.Seconds()
	toolsResult = append(toolsResult, solhint)
	helper.WriteFileExtra(fmt.Sprint("pass solhint"), "log.txt")




	honeybadgerStart := time.Now()
	honeybadgerDetail, err := docker.RunHoneyBadgerAnalysisWithTimeOut(mainFile, contractFolder, remapping)
	if err != nil{
		helper.WriteFileExtra(err.Error(), "log.txt")
		// return models.Result{}, err
	}
	// fmt.Print(detail)
	honeybadgerSumUp := docker.GetHoneyBadgerSumUp(honeybadgerDetail, err)
	honeybadgerEnd := time.Since(honeybadgerStart)
	var honeybadger models.ToolResult
	honeybadger.ToolName = "honeybadger"
	honeybadger.SumUps = honeybadgerSumUp
	if honeybadgerSumUp[0].Name != "HONEYBADGER ERROR" {
		honeybadger.NoError = len(honeybadgerSumUp)
	} else {
		honeybadger.NoError = 0
	}
	honeybadger.Detail = honeybadgerDetail
	honeybadger.TimeElapsed = honeybadgerEnd.Seconds()
	toolsResult = append(toolsResult, honeybadger)
	helper.WriteFileExtra(fmt.Sprint("pass honeybadger"), "log.txt")



	standardize := StandardizeResult(toolsResult)

	return models.Result{
		ContractID: contractID,
		ToolsResult: toolsResult,
		StandardizeResult: standardize,
		CreatedAt: start,
	}, nil
}


func StandardizeResult(toolsResult []models.ToolResult) models.StandardizeResults {

	vulneToServerity := make(map[string]string)
	for _, toolResult := range toolsResult{
		for _, sumUp := range toolResult.SumUps{
			vulne, severity := IdentifyVulnerability(sumUp, toolResult.ToolName)
			if vulne == ""{
				continue
			}
			_, exist := vulneToServerity[vulne]
			if !exist {
				vulneToServerity[vulne] = severity
			} else {
				newSeverity := GetHighestSeverity(vulneToServerity[vulne], severity)
				vulneToServerity[vulne] = newSeverity
			}
		}
	} 

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
		return docker.SlitherStandardize(sumup)
	} else if tool == "solhint"{
		return docker.SolhintStandardize(sumup)
	} else if tool == "honeybadger"{
		return docker.HoneyBadgerStandardize(sumup)
	}
	return  "", ""
}