package configs

import (
	"context"
	"getContractDeployment/helper"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(config Config) *mongo.Client{
	// connect to DB
	ctx, _ := context.WithTimeout(context.Background(),3600*time.Second)

	// get client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MONGO_URI))
	helper.CheckError(err)

	// Ping test DB
	err = client.Ping(ctx, nil)
	helper.CheckError(err)
	
	return client
}

// func GetDatabase(client *mongo.Client, databaseName string) *mongo.Database{
// 	database := client.Database(databaseName)
// 	return database
// }

// func GetCollection(client *mongo.Client, databaseName string, collectionName string) *mongo.Collection{
// 	database := GetDatabase(client, databaseName)
// 	collection := database.Collection(collectionName)
// 	return collection
// }

func DisconnectDB(client *mongo.Client){
	ctx, _ := context.WithTimeout(context.Background(),3600*time.Second)
	err := client.Disconnect(ctx)
	helper.CheckError(err)
}