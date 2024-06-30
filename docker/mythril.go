package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type RemappingJSON struct {
	Remappings []string `json:"remappings"`
}

func RunMythrilAnalysisWithTimeOut(file string, contractFolder string, address string, remappingJSON bool) (models.MythrilResultDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	resultChan := make(chan struct {
		result models.MythrilResultDetail
		err    error
	}, 1)

	go func() {
		result, err := RunMythrilAnalysis(file, contractFolder, address, remappingJSON)
		resultChan <- struct {
			result models.MythrilResultDetail
			err    error
		}{result, err}
	}()

	select {
	case res := <-resultChan:
		return res.result, res.err
	case <-ctx.Done():
		return models.MythrilResultDetail{}, fmt.Errorf("Mythril time out")
	}
}


func RunMythrilAnalysis(file string, contractFolder string, address string, remappingJSON bool) (models.MythrilResultDetail, error) {
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil{
		return models.MythrilResultDetail{}, helper.MakeError(err, "(mythril) new docker client")
	}

	result, resID, err := runMythrilContainer(ctx, cli, contractFolder, file, address, remappingJSON)
	if err != nil{
		return models.MythrilResultDetail{}, err
	}

	defer func(){
		err := cli.ContainerRemove(context.Background(), resID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}()


	result = []byte(helper.RemoveAfterFirstChar(string(result),"{"))
	result = helper.CleanupJSON(result)

	helper.WriteJSONToFile(string(result), "mythril.json")

	var returnResult models.MythrilResultDetail
	err = json.Unmarshal([]byte(result), &returnResult)
	if err != nil {
		return models.MythrilResultDetail{}, err
	}

	if !returnResult.Success{
		return models.MythrilResultDetail{
			Error: helper.MakeError(err, "(mythril) failed to analyze contract"),
			Issues: nil,
			Success: false,
		}, nil
	}

	return returnResult, nil
}

func runMythrilContainer(ctx context.Context, cli *client.Client, contractFolder string, file string, address string, remappingJSON bool) ([]byte, string, error) {

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, "", helper.MakeError(err, "(mythril) get directory")
	}

	hostConfig :=  createHostConfig(currentDir, "/mnt")

	var cmd []string = []string{"analyze"}
	if address == "" {
		path := fmt.Sprintf("/mnt/result/%s/%s", contractFolder, file)
		cmd = append(cmd, path,"--execution-timeout", "120", "-o", "json")
	} else {
		cmd = append(cmd, "-a", address, "--rpc", "infura-mainnet", "--infura-id", "9115f58b5ee548ce91fe69cd2d2f8dec", "--execution-timeout", "120", "-o", "json")
	}
	fmt.Println(cmd)
	if remappingJSON {
		cmd = append(cmd, "--solc-json", "/mnt/result/remapping/mythril_remappings.json")
	}
	
	resp, err := createContainer(ctx, cli, "mythril/myth", false, cmd, &hostConfig, "")
	if err != nil {
		return nil, resp.ID, helper.MakeError(err, "(mythril) create container")
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, resp.ID, helper.MakeError(err, "(mythril) container start")
	}

	err = waitContainer(ctx, cli, resp)
	if err != nil {
		return nil, resp.ID, helper.MakeError(err, "(mythril) wait container finish")
	}

	result, err := retrieveContainerLogs(ctx, cli, resp)
	if err != nil {
		return nil, resp.ID, helper.MakeError(err, "(mythril) retrieve container logs")
	}
	// result = helper.RemoveAfterFirstChar(result, "{")
	
	return result, resp.ID, nil
}

func GetMythrilSumUp(detail models.MythrilResultDetail, err error) []models.SumUp{
	var sumups []models.SumUp
	if err != nil{
		sumups = append(sumups, models.SumUp{
			Name: "MYTHRIL ERROR",
			Description: err.Error(),
			Severity: "",
		})
		return sumups
	}

	for _, issue := range detail.Issues{
		sumup := models.SumUp{
			Name: MythrilVulnaClass[issue.SwcID],
			Description: issue.Description,
			Severity: issue.Severity,
			Location: models.Location{
				Contract: issue.Contract,
				Function: issue.Function,
				Line: []int{issue.LineNo},
			},
		}
		sumups = append(sumups, sumup)
	}
	return sumups
}
