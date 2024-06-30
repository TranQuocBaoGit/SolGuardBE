package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"getContractDeployment/configs"
	"getContractDeployment/docker"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// type AnalysisReturn struct{
// 	Result models.Result `json:"result"`
// 	UniqueID string 	`json:"unique_id"`
// }

func PostAnalysisByFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData FileAnalysisFormData
		err := c.ShouldBind(&formData)
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) Analysis by file: ", formData.File.Filename), "log.txt")
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, helper.MakeError(err, "(bind error)").Error(), nil)
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

		contractData := models.Contract{
			ContractID: -1,
			Address: "",
			ChainID: -1,
			NoContract: 1,
			MainContract: file.Filename,
			Content: nil,
		}
		
		dataReturn, err := returnFullResult(contractData, contractFolder, false)
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
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow()," (Analysis API) Analysis contract: ",formData.Address), "log.txt")
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow()," (Analysis API) From address     : ",formData.WalletAddress),"log.txt")
		walletAddress := strings.ToLower(formData.WalletAddress)
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, helper.MakeError(err, "(bind error)").Error(), nil)
			return
		}

		chainid, err := strconv.Atoi(formData.ChainID)
		if err != nil{
			responsesReturn(c, http.StatusInternalServerError, helper.MakeError(err, "(bind error)").Error(), nil)
			return
		}

		contractData, err := getContractSourceCode(chainid, formData.Address)
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

		remap := false
		// walletAddress := formData.WalletAddress
		if walletAddress == ""{
			walletAddress = "0x0"
		}

		analysisResult, historyResult, err := getAnalysisResult(walletAddress, contractData, contractFolder, formData.Dapp, formData.Decision, remap)
		if err != nil{
			helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow()," (Analysis API) Analyze ID: ",analysisResult.AnalyzeID),"log.txt")

		response := make(map[string]interface{})
		data, _ := json.Marshal(analysisResult)
		json.Unmarshal(data, &response)
		response["unique_id"] = historyResult.UniqueID
		helper.DeleteContractsFolder(fmt.Sprintf("./result/%s",contractFolder))
		responsesReturn(c, http.StatusOK, "success", response)
	}
}

func getAnalysisResult(walletAddress string, contractData models.Contract, contractFolder, dapp, decision string, remapping bool)(models.Result, models.AnalyzeHistory, error){

	timePerform := time.Now()

	if contractData.ContractID != -1{
		var result models.Result
		config, err := configs.LoadConfig(".")
		if err != nil {
			return models.Result{}, models.AnalyzeHistory{}, err
		}
		client := configs.ConnectDB(config)
		ctx, _ := context.WithTimeout(context.Background(), 3600*time.Second)
		analysisCol := getCollection(client, config, "analysis")
		userCol := getCollection(client, config, "user")
	
		result, err = getAnalysisViaContractFromDB(ctx, analysisCol, contractData.ContractID)
		if err == nil {
			analyzeHistory := models.AnalyzeHistory{
				UniqueID: strings.ToLower(uuid.New().String()),
				AnalyzeID: result.AnalyzeID,
				ContractAddress: contractData.Address,
				ChainID: contractData.ChainID,
				TimePerform: timePerform,
				Dapp: dapp,
				Decision: decision,
			}
			AddAnalysisToUserHistory(ctx, userCol, walletAddress, analyzeHistory)
			helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) Found analysis in database: ", result.AnalyzeID), "log.txt")
			return result, analyzeHistory, nil
		} else if err != mongo.ErrNoDocuments {
			return models.Result{}, models.AnalyzeHistory{}, err
		}
		
		for _, file := range contractData.Content{
			if file.ContractName == contractData.MainContract{
				content, err := helper.ChangeSolidityVersion(file.ContractContent, "0.8.22")
				if err != nil {
					return models.Result{}, models.AnalyzeHistory{}, err
				}
				helper.WriteFile(content, filepath.Join("result", contractFolder, file.ContractName))
			}
		}

		fmt.Println(contractData.Address)

		fullResult, err := returnFullResult(contractData, contractFolder, remapping)
		if err != nil {
			return models.Result{}, models.AnalyzeHistory{}, err
		}
		err = saveAnalysisToDB(ctx, analysisCol, &fullResult)
		analyzeHistory := models.AnalyzeHistory{
			UniqueID: strings.ToLower(uuid.New().String()),
			AnalyzeID: fullResult.AnalyzeID,
			ContractAddress: contractData.Address,
			TimePerform: timePerform,
			Dapp: dapp,
			Decision: decision,
		}
		AddAnalysisToUserHistory(ctx, userCol, walletAddress, analyzeHistory)
		helper.WriteFileExtra(fmt.Sprint(" (Analysis API) Save analysis success at ID: ", fullResult.AnalyzeID), "log.txt")
		if err != nil {
			return models.Result{}, models.AnalyzeHistory{}, err
		}

		return fullResult, analyzeHistory, nil
	}

	fullResult, err := returnFullResult(contractData, contractFolder, remapping)
	if err != nil {
		return models.Result{}, models.AnalyzeHistory{}, err
	}

	return fullResult, models.AnalyzeHistory{}, nil

}

func returnFullResult(contractData models.Contract, contractFolder string, remapping bool) (models.Result, error){

	var toolsResult []models.ToolResult
	start := time.Now()

	mythrilStart := time.Now()
	mythrilDetail, err := docker.RunMythrilAnalysisWithTimeOut(contractData.MainContract, contractFolder, contractData.Address, remapping)
	if err != nil{
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) (Error) ", err.Error(), ), "log.txt")
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
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) Mythril Done"), "log.txt")



	slitherStart := time.Now()
	slitherDetail, err := docker.RunSlitherAnalysisWithTimeOut(contractData.MainContract, contractFolder, contractData.Address, remapping)
	if err != nil{
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) (Error) ", err.Error(), ), "log.txt")
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
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) Slither Done"), "log.txt")



	solhintStart := time.Now()
	solhintDetail, err := docker.RunSolHintAnalysisWithTimeOut(contractData.MainContract, contractFolder, remapping)
	if err != nil{
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) (Error) ", err.Error(), ), "log.txt")
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
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) Solhint Done"), "log.txt")




	honeybadgerStart := time.Now()
	honeybadgerDetail, err := docker.RunHoneyBadgerAnalysisWithTimeOut(contractData.MainContract, contractFolder, remapping)
	if err != nil{
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) (Error) ", err.Error(), ), "log.txt")
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
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) HoneyBadger Done"), "log.txt")

	running_ids, err := docker.ListRunningContainerIDs()
	if err != nil{
		helper.WriteFileExtra(err.Error(), "log.txt")
		// return models.Result{}, err
	}

	for _, id := range running_ids{
		err = docker.RemoveContainerByID(id)
		if err != nil{
			helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) (Error) ", err.Error(), ), "log.txt")
			// return models.Result{}, err
		}
	}

	standardize, err := StandardizeResult(toolsResult)
	if err != nil{
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis API) (Error) ", err.Error(), ), "log.txt")
		return models.Result{}, err
	}

	return models.Result{
		ContractID: contractData.ContractID,
		ToolsResult: toolsResult,
		StandardizeResult: standardize,
		CreatedAt: start,
	}, nil
}

