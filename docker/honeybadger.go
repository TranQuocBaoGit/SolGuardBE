package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"os"
	"reflect"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RunHoneyBadgerAnalysisWithTimeOut(file string, contractFolder string, remappingJSON bool) (models.HoneyBadgerResultDetail, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Channel to receive the result
	resultChan := make(chan struct {
		result models.HoneyBadgerResultDetail
		err    error
	}, 1)

	// Run the fetchMythrilResult function in a separate goroutine
	go func() {
		result, err := RunHoneyBadgerAnalysis(file, contractFolder, remappingJSON)
		resultChan <- struct {
			result models.HoneyBadgerResultDetail
			err    error
		}{result, err}
	}()

	// Use a select statement to wait for either the result or the context timeout
	select {
	case res := <-resultChan:
		return res.result, res.err
	case <-ctx.Done():
		return models.HoneyBadgerResultDetail{}, fmt.Errorf("Honeybadger time out")
	}
}

func RunHoneyBadgerAnalysis(file string, contractFolder string, remappingJSON bool) (models.HoneyBadgerResultDetail, error) { //(string, error)
	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil{
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) new docker client")
	}

	result, err := runHoneyBadgerContainer(ctx, cli, file, contractFolder, remappingJSON)
	if err != nil{
		return models.HoneyBadgerResultDetail{}, err
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return models.HoneyBadgerResultDetail{}, err
	}

	helper.WriteJSONToFile(string(jsonData), "honeybadger.json")

	return result, nil
}

func runHoneyBadgerContainer(ctx context.Context, cli *client.Client, file string, contractFolder string, remappingJSON bool) (models.HoneyBadgerResultDetail, error) {

	currentDir, err := os.Getwd()
	if err != nil {
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) get directory")
	}

	hostConfig := createHostConfig(currentDir, "/share")

	resp, err := createContainer(ctx, cli, "christoftorres/honeybadger", true, nil, &hostConfig, "")
	if err != nil {
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) create container")
	}
	defer func(){
		err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}()

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) container start")
	}


	var setupCmd []string = []string{"bash", "-c", fmt.Sprintf(`mkdir results && touch results/%s`, helper.ChangeFromSolToJson(file))}
	// fmt.Println(cmd)

	_, err = performExec(cli, resp, setupCmd)
	if err != nil {
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) perform setup execution")
	}

	var runAnalysisCmd []string = []string{"bash", "-c", fmt.Sprintf("python honeybadger/honeybadger.py -s /share/result/%s/%s -j", contractFolder, file)}
	_, err = performExec(cli, resp, runAnalysisCmd)
	if err != nil {
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) perform analyze execution")
	}

	var cmd []string = []string{"bash", "-c", fmt.Sprintf("cat results/%s", helper.ChangeFromSolToJson(file))}
	result, err := performExec(cli, resp, cmd)
	if err != nil {
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) get analyze execution")
	}

	resultStr := helper.RemoveAfterXChar(string(result),"{", 2)
	if resultStr == ""{
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) empty result")
	}
	resultStr = resultStr[:len(resultStr)-1]
	// fmt.Print(resultStr)


	var returnResult models.HoneyBadgerResultDetail
	if err := json.Unmarshal([]byte(resultStr), &returnResult); err != nil {
		return models.HoneyBadgerResultDetail{}, helper.MakeError(err, "(honeybadger) json unmarshal")
	}
	// fmt.Print(returnResult)

	return returnResult, nil
}


func GetHoneyBadgerSumUp(detail models.HoneyBadgerResultDetail, err error) []models.SumUp{
	var sumups []models.SumUp
	if err != nil {
		sumups = append(sumups, models.SumUp{
			Name: "HONEYBADGER ERROR",
			Description: "Honeybadger fail to analyze contract",
			Severity: "",
		})
		return sumups
	}
	balanceDisorder := reflect.TypeOf(detail.BalanceDisorder)
	hiddenStateUpdate  := reflect.TypeOf(detail.HiddenStateUpdate)
	hiddenTransfer := reflect.TypeOf(detail.HiddenTransfer)
	skipEmptyStringLiteral  := reflect.TypeOf(detail.SkipEmptyStringLiteral)
	inheritanceDisorder := reflect.TypeOf(detail.InheritanceDisorder)
	uninitialisedStruct := reflect.TypeOf(detail.UninitialisedStruct)
	strawManContract := reflect.TypeOf(detail.StrawManContract)
	typeDeductionOverflow := reflect.TypeOf(detail.TypeDeductionOverflow)

	if balanceDisorder.Kind() == reflect.String{
		sumups = append(sumups, models.SumUp{
			Name: "Balance Disorder",
			Severity: "Medium",
			Description: detail.BalanceDisorder.(string),
		})
	}
	if hiddenStateUpdate.Kind() == reflect.String{
		sumups = append(sumups, models.SumUp{
			Name: "Hidden State Update",
			Severity: "Medium",
			Description: detail.HiddenStateUpdate.(string),
		})
	}
	if hiddenTransfer.Kind() == reflect.String{
		sumups = append(sumups, models.SumUp{
			Name: "Hidden Transfer",
			Severity: "Medium",
			Description: detail.HiddenTransfer.(string),
		})
	}
	if skipEmptyStringLiteral.Kind() == reflect.String{
		sumups = append(sumups, models.SumUp{
			Name: "Skip String Literal",
			Severity: "Medium",
			Description: detail.SkipEmptyStringLiteral.(string),
		})
	}
	if inheritanceDisorder.Kind() == reflect.String{
		sumups = append(sumups, models.SumUp{
			Name: "Inheritance Disorder",
			Severity: "Medium",
			Description: detail.InheritanceDisorder.(string),
		})
	}
	if uninitialisedStruct.Kind() == reflect.String{
		sumups = append(sumups, models.SumUp{
			Name: "Uninitalised Struct",
			Severity: "Medium",
			Description: detail.UninitialisedStruct.(string),
		})
	}
	if strawManContract.Kind() == reflect.String{
		sumups = append(sumups, models.SumUp{
			Name: "Straw Man Hat",
			Severity: "Medium",
			Description: detail.StrawManContract.(string),
		})
	}
	if typeDeductionOverflow.Kind() == reflect.String{
		sumups = append(sumups, models.SumUp{
			Name: "Type Deduction Overflow",
			Severity: "Medium",
			Description: detail.TypeDeductionOverflow.(string),
		})
	}

	return sumups
}