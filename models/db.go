package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type StandardizeResults struct {
	NoError int `json:"no_error"`
	Result []StandardizeResult `json:"result"`
}

type StandardizeResult struct {
	Name string `json:"name"`
	Severity string `json:"severity"`
}

type Contract struct {
	ContractID 		int  				`bson:"_id,omitempty" json:"_id,omitempty"`
	Address    		string 				`bson:"address,omitempty" json:"address,omitempty"`
	ChainID			int					`bson:"chain_id,omitempty" json:"chain_id,omitempty"`
	NoContract		int					`bson:"no_contract,omitempty" json:"no_contract,omitempty"`
	MainContract 	string 				`bson:"main_contract,omitempty" json:"main_contract,omitempty"`
	Content    		[]ContractContent	`bson:"content,omitempty" json:"content,omitempty"`
	FromDB			bool				`bson:"from_db" json:"from_db"`
}

type ContractContent struct{
	ContractName 		string `bson:"contract_name,omitempty" json:"contract_name,omitempty"`
	ContractContent		string `bson:"content,omitempty" json:"content,omitempty"`
}

// AutoIncrementID increments the ID field automatically
func AutoIncrementID(col *mongo.Collection, model *Contract) error {
	count, err := col.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	model.ContractID = int(count + 1)
	return nil
}

type Result struct {
	ContractID       int   			`bson:"contract_id,omitempty" json:"contract_id,omitempty"`
//  NoTool           int          	`bson:"no_tool,omitempty" json:"no_tool,omitempty"`
	ToolsResult      []ToolResult 	`bson:"tools_result,omitempty" json:"tools_result,omitempty"`
//  TotalTimeElapsed float64        `bson:"total_time_elapsed,omitempty" json:"total_time_elapsed,omitempty"`
}

type ToolResult struct {
	ToolName    string      `bson:"name,omitempty" json:"name,omitempty"`
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