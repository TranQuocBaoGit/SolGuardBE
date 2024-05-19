package controller

import (
	"context"
	"getContractDeployment/helper"
	"getContractDeployment/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


func getSourceCodeFromDB(ctx context.Context, collection *mongo.Collection, addressValue string) (models.Contract, error) {	
	filter := bson.D{{"address", addressValue}}

	var result models.Contract
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err ==  mongo.ErrNoDocuments{
		return models.Contract{}, err
	} else if err != nil {
		return models.Contract{}, helper.MakeError(err, "(db_source_code) get source code")
	}
	
	return result, nil
}


func saveSourceCodeToDB(ctx context.Context, collection *mongo.Collection, data *models.Contract) error {

	err := models.AutoIncrementID(collection, data)
	if err != nil {
		return helper.MakeError(err, "(db_source_code) auto increment id")
	}

	_, err = collection.InsertOne(ctx, data)
	if err != nil {
		return helper.MakeError(err, "(db_source_code) save source code")
	}

	return nil
}