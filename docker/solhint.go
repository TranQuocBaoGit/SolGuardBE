package docker

// package docker

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"getContractDeployment/helper"
// 	"getContractDeployment/models"
// 	"os"

// 	"github.com/docker/docker/api/types"
// 	"github.com/docker/docker/client"
// )

// func RunSolhintAnalysis(file string, contractFolder string, remappingJSON bool) (models.SolhintResultDetail, error) { //(string, error)
// 	ctx := context.Background()

// 	// Create a Docker client
// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil{
// 		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) new docker client")
// 	}

// 	result, err := runSolhintContainer(ctx, cli, file, contractFolder, remappingJSON)
// 	if err != nil{
// 		return models.SolhintResultDetail{}, err
// 	}

// 	jsonData, err := json.Marshal(result)
// 	if err != nil {
// 		return models.SolhintResultDetail{}, err
// 	}

// 	helper.WriteJSONToFile(string(jsonData), "solhint.json")
// 	// result = helper.PreprocessJSON(result)
// 	// helper.WriteFile(result, "wtf.txt")
	

// 	// var returnResult models.SlitherResultDetail
// 	// err = json.Unmarshal([]byte(result), &returnResult)
// 	// if err != nil {
// 	// 	helper.WriteFile(result, "wtf.txt")
// 	// 	return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) json unmarshal") 
// 	// }

// 	if !result.Success{
// 		return models.SolhintResultDetail{
// 			Error: helper.MakeError(err, "(slither) failed to analyze contract"),
// 			Results: models.SlitherDetectorDetail{},
// 			Success: false,
// 		}, nil
// 	}

// 	return result, nil
// }

// func runSolhintContainer(ctx context.Context, cli *client.Client, file string, contractFolder string, remappingJSON bool) (models.SolhintResultDetail, error) {

// 	currentDir, err := os.Getwd()
// 	if err != nil {
// 		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) get directory")
// 	}

// 	hostConfig := createHostConfig(currentDir, "/share")

// 	resp, err := createContainer(ctx, cli, "solhint", true, nil, &hostConfig, "")
// 	if err != nil {
// 		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) create container")
// 	}
// 	defer func(){
// 		err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
// 	if err != nil {
// 		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) container start")
// 	}

// 	var setupCmd []string = []string{"sh", "-c", fmt.Sprintf("solhint --init &&  && sed -i 's/\"extends\": \"solhint:default\"/\"extends\": \"solhint:recommended\"/' .solhint.json")}
// 	// fmt.Println(cmd)

// 	_, err = performExec(cli, resp, setupCmd)
// 	if err != nil {
// 		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) perform setup execution")
// 	}

// 	var cmd []string = []string{"sh", "-c", fmt.Sprintf("solhint /share/%s/%s", contractFolder, file)}
// 	result, err := performExec(cli, resp, cmd)
// 	if err != nil {
// 		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) perform analyze execution")
// 	}


// 	var returnResult models.SolhintResultDetail
// 	if err := json.Unmarshal(result, &returnResult); err != nil {
// 		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) json unmarshal")
// 	}

// 	// decoder := json.NewDecoder(execResp.Reader)
// 	// err = decoder.Decode(&returnResult)
// 	// if err != nil {
// 	// 	return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) decode result")
// 	// }

// 	// result = helper.RemoveAfterFirstChar(result,"{")

// 	return returnResult, nil
// }


// func GetSolhintSumUp(detail models.SolhintResultDetail) []models.SumUp{
// 	var sumups []models.SumUp
// 	for _, issue := range detail.Issues{
// 		sumup := models.SumUp{
// 			Name: SlitherVulnaClass[issue.ID],
// 			Description: issue.Description,
// 			Severity: issue.Impact,
// 		}
// 		sumups = append(sumups, sumup)
// 	}
// 	return sumups
// }