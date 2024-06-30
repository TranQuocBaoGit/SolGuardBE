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


func getSourceCodeViaAddressFromDB(ctx context.Context, contractCollection *mongo.Collection, addressValue string) (models.Contract, error) {	
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code DB) Get source code of: ", addressValue), "log.txt")

	filter := bson.M{"address": bson.M{"$regex": "^" + regexp.QuoteMeta(addressValue) + "$", "$options": "i"}}
	var result models.Contract
	err := contractCollection.FindOne(ctx, filter).Decode(&result)
	if err ==  mongo.ErrNoDocuments{
		return models.Contract{}, err
	} else if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code DB) (Error) ", err.Error()), "log.txt")
		return models.Contract{}, helper.MakeError(err, "(db_source_code) get source code")
	}
	
	return result, nil
}

func getSourceCodeViaIDFromDB(ctx context.Context, contractCollection *mongo.Collection, id int) (models.Contract, error) {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code DB) Get source code ID: ", id), "log.txt")

	filter := bson.D{{"_id", id}}
	var result models.Contract
	err := contractCollection.FindOne(ctx, filter).Decode(&result)
	if err ==  mongo.ErrNoDocuments{
		return models.Contract{}, err
	} else if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code DB) (Error) ", err.Error()), "log.txt")
		return models.Contract{}, helper.MakeError(err, "(db_source_code) get source code")
	}
	
	return result, nil
}


func saveSourceCodeToDB(ctx context.Context, contractCollection *mongo.Collection, data *models.Contract) error {
	helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code DB) Save source code of : ", data.Address), "log.txt")

	err := models.AutoIncrementIDContract(contractCollection, data)
	if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code DB) (Error) ", err.Error()), "log.txt")
		return helper.MakeError(err, "(db_source_code) auto increment id")
	}

	_, err = contractCollection.InsertOne(ctx, data)
	if err != nil {
		helper.WriteFileExtra(fmt.Sprint(helper.GetTimeNow(), " (Source Code DB) (Error) ", err.Error()), "log.txt")
		return helper.MakeError(err, "(db_source_code) save source code")
	}

	return nil
}