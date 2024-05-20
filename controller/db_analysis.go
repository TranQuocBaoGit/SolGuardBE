package controller

import (
	"context"
	"getContractDeployment/helper"
	"getContractDeployment/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getAnalysisViaContractFromDB(ctx context.Context, collection *mongo.Collection, contractID int) (models.Result, error) {
	filter := bson.D{{"contractID", contractID}}

	var result models.Result
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments{
		return models.Result{}, err
	} else if err != nil {
		return models.Result{}, helper.MakeError(err, "(db_analysis) get analysis from contract id")
	}
	
	return result, nil
}

func getAnalysisViaUserFromDB(ctx context.Context, collection *mongo.Collection, analyzeID int) (models.Result, error) {
	filter := bson.M{"_id": analyzeID}

	var result models.Result
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return models.Result{}, err
		}
		return models.Result{}, helper.MakeError(err, "(db_analysis) get analysis from user")
	}

	return result, nil
}

func saveAnalysisToDB(ctx context.Context, collection *mongo.Collection, data *models.Result) error {

	err := models.AutoIncrementIDAnalysis(collection, data)
	if err != nil {
		return helper.MakeError(err, "(db_analysis) auto increment id")
	}

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return helper.MakeError(err, "(db_analysis) save analysis")
	}

	return nil
}

// TO DO