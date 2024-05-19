package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func RunSlitherAnalysis(file string, contractFolder string, remappingJSON bool) (models.SlitherResultDetail, error) { //(string, error)
	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil{
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) new docker client")
	}

	result, err := runSlitherContainer(ctx, cli, file, contractFolder, remappingJSON)
	if err != nil{
		return models.SlitherResultDetail{}, err
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return models.SlitherResultDetail{}, err
	}

	helper.WriteJSONToFile(string(jsonData), "slither.json")
	// result = helper.PreprocessJSON(result)
	// helper.WriteFile(result, "wtf.txt")
	

	// var returnResult models.SlitherResultDetail
	// err = json.Unmarshal([]byte(result), &returnResult)
	// if err != nil {
	// 	helper.WriteFile(result, "wtf.txt")
	// 	return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) json unmarshal") 
	// }

	if !result.Success{
		return models.SlitherResultDetail{
			Error: helper.MakeError(err, "(slither) failed to analyze contract"),
			Results: models.SlitherDetectorDetail{},
			Success: false,
		}, nil
	}

	return result, nil
}

func runSlitherContainer(ctx context.Context, cli *client.Client, file string, contractFolder string, remappingJSON bool) (models.SlitherResultDetail, error) {

	currentDir, err := os.Getwd()
	if err != nil {
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) get directory")
	}

	hostConfig := createHostConfig(currentDir, "/share")

	resp, err := createContainer(ctx, cli, "trailofbits/eth-security-toolbox", true, nil, &hostConfig, "")
	if err != nil {
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) create container")
	}
	defer func(){
		err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}()

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) container start")
	}

	var cmd []string = []string{"sh", "-c", fmt.Sprintf("ls && cd /share/result/%s && slither %s --json -", contractFolder, file)}
	// fmt.Println(cmd)

	result, err := performExec(cli, resp, cmd)
	if err != nil {
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) perform execution")
	}

	result = []byte(helper.RemoveAfterFirstChar(string(result),"{"))

	var returnResult models.SlitherResultDetail
	if err := json.Unmarshal(result, &returnResult); err != nil {
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) json unmarshal")
	}

	// decoder := json.NewDecoder(execResp.Reader)
	// err = decoder.Decode(&returnResult)
	// if err != nil {
	// 	return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) decode result")
	// }

	// result = helper.RemoveAfterFirstChar(result,"{")

	return returnResult, nil
}


func GetSlitherSumUp(detail models.SlitherResultDetail) []models.SumUp{
	var sumups []models.SumUp
	for _, issue := range detail.Results.Detectors{
		if issue.Impact == "informaltional" || issue.Impact == "Informational" || issue.Impact == "Optimization" || issue.Impact == "optimization"{
			continue
		}
		sumup := models.SumUp{
			Name: SlitherVulnaClass[issue.ID],
			Description: issue.Description,
			Severity: issue.Impact,
		}
		sumups = append(sumups, sumup)
	}
	return sumups
}