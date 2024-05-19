package controller

import (
	"context"
	"getContractDeployment/configs"
	"getContractDeployment/helper"

	"go.mongodb.org/mongo-driver/mongo"
)


func checkRecordExistsByFieldValue(ctx context.Context, col *mongo.Collection, fieldName string, fieldValue interface{}) (bool, error) {
	filter := map[string]interface{}{
		fieldName: fieldValue,
	}

	err := col.FindOne(ctx, filter).Err()
	if err == mongo.ErrNoDocuments {
		return false, nil
	} else if err != nil {
		return false, helper.MakeError(err, "(db_overall) check record")
	}

	return true, nil
}

func getCollection(client *mongo.Client, config configs.Config, colName string) *mongo.Collection{
	col := client.Database(config.DATABASE_NAME).Collection(colName)
	return col
}