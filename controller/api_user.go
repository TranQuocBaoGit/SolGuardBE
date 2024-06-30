package controller

import (
	"context"
	"fmt"
	"getContractDeployment/configs"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)


func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData User
		err := c.ShouldBind(&formData)
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow()," (User API) Login: ",formData.WalletAddress), "log.txt")
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, helper.MakeError(err, "(bind error)").Error(), nil)
			return
		}

		config, err := configs.LoadConfig(".")
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		client := configs.ConnectDB(config)
		ctx, _ := context.WithTimeout(context.Background(), 3600*time.Second)
		col := getCollection(client, config, "user")
		walletAddress := strings.ToLower(formData.WalletAddress)

		check, err := checkRecordExistsByFieldValue(ctx, col, "wallet_address", walletAddress)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		if !check{
			err = AddNewUser(ctx, col,  walletAddress)
			if err != nil {
				responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
				return
			}
		}

		responsesReturn(c, http.StatusOK, "successful", formData.WalletAddress)
	}
}

func GetUserHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		walletAddress := c.Query("wallet_address")
		// analyzeIDStr := c.Query("analyzeID")
		uniqueIDStr := c.Query("id")
		walletAddress = strings.ToLower(walletAddress)
		if uniqueIDStr != ""{
			helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow()," (User API) Get analyze history of unique ID ", uniqueIDStr, " from user: ",walletAddress), "log.txt")
			oneUserHistory, err := getOneUserHistory(walletAddress, uniqueIDStr)
			if err != nil {
				responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
				return
			}
	
			responsesReturn(c, http.StatusOK, "successful", oneUserHistory)
			return
		}
		
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow()," (User API) Get all user history: ",walletAddress), "log.txt")

		userOverallHistory, err := getAllUserOverallHistory(walletAddress)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, err.Error(), nil)
			return
		}

		responsesReturn(c, http.StatusOK, "successful", userOverallHistory)
	}
}

func DeleteHistory() gin.HandlerFunc {
	return func(c *gin.Context) {
		var formData DeleteHistoryFormData
		if err := c.ShouldBind(&formData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		config, err := configs.LoadConfig(".")
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, helper.MakeError(err, "(load config)").Error(), nil)
			return
		}

		client := configs.ConnectDB(config)
		ctx, _ := context.WithTimeout(context.Background(), 3600*time.Second)
		userCol := getCollection(client, config, "user")
	
		walletAddress := strings.ToLower(formData.WalletAddress)
		err = DeleteAnalysisFromUserHistory(ctx, userCol, walletAddress, formData.HistoryAnalyzeID)
		if err != nil {
			responsesReturn(c, http.StatusInternalServerError, helper.MakeError(err, "(load config)").Error(), nil)
			return
		}
	
		responsesReturn(c, http.StatusOK, "successful", "Successfully delete history analysis")
	}
}

type UserHistoryOverall struct{
	UniqueID string `json:"unique_id"`
	AnalyzeID int `json:"analyze_id"`
	ContractAddress string `json:"contract_address"`
	ChainID int `json:"chainid"`
	TimeCreated int `json:"time_created"`
	StandardizeResults models.StandardizeResults `json:"standardize_result"`
	Dapp string `json:"dApp"`
	Decision string `json:"decision"`
}

type UserHistory struct{
	UniqueID string `json:"unique_id"`
	AnalyzeID int `json:"analyze_id"`
	ContractAddress string `json:"contract_address"`
	ChainID int `json:"chainid"`
	TimeCreated int `json:"time_created"`
	StandardizeResults models.StandardizeResults `json:"standardize_result"`
	ToolsResult []models.ToolResult `json:"tools_result"`
	Dapp string `json:"dApp"`
	Decision string `json:"decision"`
}	

func getAllUserOverallHistory(walletAddress string) ([]UserHistoryOverall, error) {
	config, err := configs.LoadConfig(".")
	if err != nil {
		return nil, err
	}
	client := configs.ConnectDB(config)
	ctx, _ := context.WithTimeout(context.Background(), 3600*time.Second)
	userCol := getCollection(client, config, "user")
	analysisCol := getCollection(client, config, "analysis")

	analyzeHistories, err := GetUserHistoryFromDB(ctx, userCol, walletAddress)
	if err != nil {
		return nil, err
	}
	
	var userOverallHistory []UserHistoryOverall
	for _, analyzeHistory := range  analyzeHistories{
		analysis, err := getAnalysisViaUserFromDB(ctx, analysisCol, analyzeHistory.AnalyzeID)
		if err != nil {
			return nil, err
		}
		temp := UserHistoryOverall{
			UniqueID: analyzeHistory.UniqueID,
			AnalyzeID: analyzeHistory.AnalyzeID,
			ContractAddress: analyzeHistory.ContractAddress,
			ChainID: analyzeHistory.ChainID,
			TimeCreated: int(analyzeHistory.TimePerform.Unix()),
			StandardizeResults: analysis.StandardizeResult,
			Dapp: analyzeHistory.Dapp,
			Decision: analyzeHistory.Decision,
		}
		userOverallHistory = append(userOverallHistory, temp)
	}

	return userOverallHistory, nil
}


func getOneUserHistory(walletAddress, uniqueID string)(UserHistory, error){
	config, err := configs.LoadConfig(".")
	if err != nil {
		return UserHistory{}, err
	}
	client := configs.ConnectDB(config)
	ctx, _ := context.WithTimeout(context.Background(), 3600*time.Second)
	// userCol := getCollection(client, config, "user")
	analysisCol := getCollection(client, config, "analysis")

	// id, err :=  strconv.Atoi(uniqueID)
	// if err != nil {
	// 	return UserHistory{}, err
	// }

	// analyzeHistory, err := GetOneUserHistoryFromDB(ctx, userCol, walletAddress, id)
	// if err != nil {
	// 	return UserHistory{}, err
	// }

	// analysis, err := getAnalysisViaUserFromDB(ctx, analysisCol, id)
	// if err != nil {
	// 	return UserHistory{}, err
	// }

	allOverallHistories, err := getAllUserOverallHistory(walletAddress)
	if err != nil {
		return UserHistory{}, err
	}

	for _, oneOverallHistories := range allOverallHistories {
		if oneOverallHistories.UniqueID == uniqueID{
			analysis, err := getAnalysisViaUserFromDB(ctx, analysisCol, oneOverallHistories.AnalyzeID)
			if err != nil {
				return UserHistory{}, err
			}
			return UserHistory{
				UniqueID: oneOverallHistories.UniqueID,
				AnalyzeID: oneOverallHistories.AnalyzeID,
				ContractAddress: oneOverallHistories.ContractAddress,
				ChainID: oneOverallHistories.ChainID,
				StandardizeResults: analysis.StandardizeResult,
				ToolsResult: analysis.ToolsResult,
				TimeCreated: int(analysis.CreatedAt.Unix()),
				Dapp: oneOverallHistories.Dapp,
				Decision: oneOverallHistories.Decision,
			}, nil
		}
	}


	return UserHistory{}, nil
}
