package controller

import (
	"context"
	"getContractDeployment/helper"
	"getContractDeployment/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getAnalysisFromDB(ctx context.Context, collection *mongo.Collection, contractID int) (models.Result, error) {
	filter := bson.D{{"contractID", contractID}}

	var result models.Result
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments{
		return models.Result{}, err
	} else if err != nil {
		return models.Result{}, helper.MakeError(err, "(db_analysis) get analysis")
	}
	
	return result, nil
}


func saveAnalysisToDB(ctx context.Context, collection *mongo.Collection, data *models.Result) error {

	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		return helper.MakeError(err, "(db_analysis) save analysis")
	}

	return nil
}

// TO DO