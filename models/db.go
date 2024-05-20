package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	WalletAddress string 	`bson:"wallet_address" json:"wallet_address"`
	AnalyzeHistory []int 	`bson:"analyze_history" json:"analyze_history"`
}

type Contract struct {
	ContractID 		int  				`bson:"_id,omitempty" json:"_id,omitempty"`
	Address    		string 				`bson:"address,omitempty" json:"address,omitempty"`
	ChainID			int					`bson:"chain_id,omitempty" json:"chain_id,omitempty"`
	NoContract		int					`bson:"no_contract,omitempty" json:"no_contract,omitempty"`
	MainContract 	string 				`bson:"main_contract,omitempty" json:"main_contract,omitempty"`
	Content    		[]ContractContent	`bson:"content,omitempty" json:"content,omitempty"`
}

type ContractContent struct{
	ContractName 		string `bson:"contract_name,omitempty" json:"contract_name,omitempty"`
	ContractContent		string `bson:"content,omitempty" json:"content,omitempty"`
}

func AutoIncrementIDContract(col *mongo.Collection, model *Contract) error {
	count, err := col.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	model.ContractID = int(count + 1)
	return nil
}

type Result struct {
	AnalyzeID        int 			`bson:"_id,omitempty" json:"analyze_id,omitempty"`
	ContractID       int   			`bson:"contract_id,omitempty" json:"contract_id,omitempty"`
	StandardizeResult StandardizeResults `bson:"standardize_result,omitempty" json:"standardize_result,omitempty"`
	ToolsResult      []ToolResult 	`bson:"tools_result,omitempty" json:"tools_result,omitempty"`
	CreatedAt 		 time.Time		`bson:"created_at,omitempty" json:"created_at,omitempty"`
}

func AutoIncrementIDAnalysis(col *mongo.Collection, model *Result) error {
	count, err := col.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	model.AnalyzeID = int(count + 1)
	return nil
}

type StandardizeResults struct {
	NoError int 				`bson:"no_error" json:"no_error"`
	Result 	[]StandardizeResult `bson:"result" json:"result"`
}

type StandardizeResult struct {
	Name 		string 	`bson:"name" json:"name"`
	Severity 	string 	`bson:"severity" json:"severity"`
}

type ToolResult struct {
	ToolName    string      `bson:"tool_name,omitempty" json:"tool_name,omitempty"`
	NoError     int         `bson:"no_error" json:"no_error"`
	SumUps      []SumUp     `bson:"sum_up" json:"sum_up"`
	Detail      interface{} `bson:"detail" json:"detail"`
	TimeElapsed float64     `bson:"time_elapsed,omitempty" json:"time_elapsed,omitempty"`
}

type SumUp struct {
	Name        string `bson:"name" json:"name"`
	Description string `bson:"description" json:"description"`
	Severity    string `bson:"severity" json:"severity"`
}