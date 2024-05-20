package controller

import (
	"context"
	"getContractDeployment/helper"
	"getContractDeployment/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddNewUser(ctx context.Context, collection *mongo.Collection, walletAddress string) error {
	var user models.User = models.User{
		WalletAddress: walletAddress,
		AnalyzeHistory: []int{},
	}

	_, err := collection.InsertOne(ctx, &user)
	if err != nil {
		return helper.MakeError(err, "(db_user) add new user")
	}
	return nil
}

func AddAnalysisToUserHistory(ctx context.Context, collection *mongo.Collection, walletAddress string, analysisCode int) error {
	filter := bson.M{"wallet_address": walletAddress}
	update := bson.M{"$push": bson.M{"analyze_history": analysisCode}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return helper.MakeError(err, "(db_user) add analysis history to user")
	}

	return nil
}

func GetUserHistory(ctx context.Context, collection *mongo.Collection, walletAddress string) ([]int, error) {

	filter := bson.M{"wallet_address": walletAddress}

	var user models.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, helper.MakeError(err, "(db_user) no user found for " + walletAddress)
		}
		return nil, helper.MakeError(err, "(db_user) fail to get user " + walletAddress)
	}

	return user.AnalyzeHistory, nil
}