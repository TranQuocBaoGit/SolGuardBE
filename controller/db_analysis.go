package controller

import (
	"context"
	"fmt"
	"getContractDeployment/helper"
	"getContractDeployment/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getAnalysisViaContractFromDB(ctx context.Context, analysisCollection *mongo.Collection, contractID int) (models.Result, error) {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis DB) Get analysis of contract ID: ", contractID), "log.txt")

	filter := bson.D{{"contract_id", contractID}}
	var result models.Result
	err := analysisCollection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments{
		return models.Result{}, err
	} else if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis DB) (Error) ", err.Error()), "log.txt")
		return models.Result{}, helper.MakeError(err, "(db_analysis) get analysis from contract id")
	}
	
	return result, nil
}

func getAnalysisViaUserFromDB(ctx context.Context, analysisCollection *mongo.Collection, analyzeID int) (models.Result, error) {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis DB) Get analysis ID: ", analyzeID), "log.txt")

	filter := bson.M{"_id": analyzeID}
	var result models.Result
	err := analysisCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Result{}, nil
		}
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis DB) (Error) ", err.Error()), "log.txt")
		return models.Result{}, helper.MakeError(err, "(db_analysis) get analysis from user")
	}

	return result, nil
}

func saveAnalysisToDB(ctx context.Context, analysisCollection *mongo.Collection, data *models.Result) error {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis DB) Save analysis ID ", data.AnalyzeID, " to database"), "log.txt")

	err := models.AutoIncrementIDAnalysis(analysisCollection, data)
	if err != nil {	
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis DB) (Error) ", err.Error()), "log.txt")
		return helper.MakeError(err, "(db_analysis) auto increment id")
	}

	_, err = analysisCollection.InsertOne(ctx, data)
	if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Analysis DB) (Error) ", err.Error()), "log.txt")
		return helper.MakeError(err, "(db_analysis) save analysis")
	}

	return nil
}

// TO DO