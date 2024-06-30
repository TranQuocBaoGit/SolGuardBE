package controller

import (
	"context"
	"fmt"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddNewUser(ctx context.Context, userCollection *mongo.Collection, walletAddress string) error {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) New user : ", walletAddress), "log.txt")
	
	var user models.User = models.User{
		WalletAddress: walletAddress,
		AnalyzeHistory: []models.AnalyzeHistory{},
	}

	_, err := userCollection.InsertOne(ctx, &user)
	if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) (Error) ", err.Error()), "log.txt")
		return helper.MakeError(err, "(db_user) add new user")
	}
	return nil
}

func AddAnalysisToUserHistory(ctx context.Context, userCollection *mongo.Collection, walletAddress string, analysisHistory models.AnalyzeHistory) error {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) Add analysis ", analysisHistory.AnalyzeID, " to user: ", walletAddress), "log.txt")

	// historyExist, err := GetOneUserHistoryFromDB(ctx, userCollection, walletAddress, analysisHistory.AnalyzeID)
	// if err != nil {
	// 	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) (Error) ", err.Error()), "log.txt")
	// 	return helper.MakeError(err, "(db_user) check analysis history error")
	// }

	// update := bson.M{}
	// if historyExist.ContractAddress != "" {
	// 	update = bson.M{
	// 		"$set": bson.M{
	// 			"analyze_history.$.time_perform":  analysisHistory.TimePerform,
	// 		},
	// 	}
	// } else {
	update := bson.M{"$push": bson.M{"analyze_history": analysisHistory}}
	// }
	filter := bson.D{{"wallet_address", walletAddress}}

	_, err := userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) (Error) ", err.Error()), "log.txt")
		return helper.MakeError(err, "(db_user) add analysis history to user")
	}

	return nil
}


func DeleteAnalysisFromUserHistory(ctx context.Context, userCollection *mongo.Collection, walletAddress, analyzeHistoryID string) error {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) Delete history ", analyzeHistoryID, " of user: ", walletAddress), "log.txt")
	var user models.User
	filter := bson.M{"wallet_address": walletAddress}
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) (Error) ", err.Error()), "log.txt")
		return helper.MakeError(err, "(db_user) find wallet to delete")
	}

	// Find and delete the history with the given AnalyzeID
	var found bool
	// var newUser models.User
	for i, history := range user.AnalyzeHistory {
		if history.UniqueID == analyzeHistoryID {
			user.AnalyzeHistory = append(user.AnalyzeHistory[:i], user.AnalyzeHistory[i+1:]...)
			found = true
			break
		}
	}

	if found{
		// Update user in MongoDB
		update := bson.M{"$set": bson.M{"analyze_history": user.AnalyzeHistory}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) (Error) ", err.Error()), "log.txt")
			return helper.MakeError(err, "(db_user) Not found user to delete")
		}
	}


	return nil
}

func GetUserHistoryFromDB(ctx context.Context, userCollection *mongo.Collection, walletAddress string) ([]models.AnalyzeHistory, error) {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) Get overall history of user: ", walletAddress), "log.txt")

	filter := bson.M{"wallet_address": bson.M{"$regex": "^" + regexp.QuoteMeta(walletAddress) + "$", "$options": "i"}}
	var user models.User
	err := userCollection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, helper.MakeError(err, "(db_user) no user found for " + walletAddress)
		}
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) (Error) ", err.Error()), "log.txt")
		return nil, helper.MakeError(err, "(db_user) fail to get user " + walletAddress)
	}


	return user.AnalyzeHistory, nil
}

func GetOneUserHistoryFromDB(ctx context.Context, userCollection *mongo.Collection, walletAddress string, analyzeID int) (models.AnalyzeHistory, error) {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (User DB) Get history ID ", analyzeID, " of user: ", walletAddress), "log.txt")

	analyzeHistories, err := GetUserHistoryFromDB(ctx, userCollection, walletAddress)
	if err != nil {
		return models.AnalyzeHistory{}, err
	}

	var analyzeHistory models.AnalyzeHistory
	for _, oneHistory := range analyzeHistories{
		if oneHistory.AnalyzeID == analyzeID {
			analyzeHistory = oneHistory
		}
	}

	if analyzeHistory.ContractAddress == "" {
		return models.AnalyzeHistory{}, nil
	}

	return analyzeHistory, nil
}

// func GetUserHistoryFromDB(ctx context.Context, collection *mongo.Collection, ids []int) ([]models.Result, error) {
// 	helper.WriteFileExtra(fmt.Sprint(time.Now(), " (User DB) Get all history of user: ", walletAddress), "log.txt")

// 	var results []models.Result
// 	for _, id := range ids{
// 		result, err := getAnalysisViaUserFromDB(ctx, collection, id)
// 		if err != nil{
// 			return nil, err
// 		}
// 		results = append(results, result)
// 	}
// 	return results, nil
// }