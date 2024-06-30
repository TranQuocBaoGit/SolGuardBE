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

func RunSlitherAnalysisWithTimeOut(file string, contractFolder string, address string, remappingJSON bool) (models.SlitherResultDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	resultChan := make(chan struct {
		result models.SlitherResultDetail
		err    error
	}, 1)

	go func() {
		result, err := RunSlitherAnalysis(file, contractFolder, address, remappingJSON)
		resultChan <- struct {
			result models.SlitherResultDetail
			err    error
		}{result, err}
	}()

	select {
	case res := <-resultChan:
		return res.result, res.err
	case <-ctx.Done():
		return models.SlitherResultDetail{}, fmt.Errorf("Slither time out")
	}
}

func RunSlitherAnalysis(file string, contractFolder string, address string, remappingJSON bool) (models.SlitherResultDetail, error) { //(string, error)
	ctx := context.Background()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil{
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) new docker client")
	}

	result, err := runSlitherContainer(ctx, cli, file, contractFolder, address, remappingJSON)
	if err != nil{
		return models.SlitherResultDetail{}, err
	}

	if !result.Success{
		return models.SlitherResultDetail{
			Error: helper.MakeError(err, "(slither) failed to analyze contract"),
			Results: models.SlitherDetectorDetail{},
			Success: false,
		}, nil
	}

	return result, nil
}

func runSlitherContainer(ctx context.Context, cli *client.Client, file string, contractFolder string, address string, remappingJSON bool) (models.SlitherResultDetail, error) {

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

	var cmd []string = []string{"sh", "-c", fmt.Sprintf("cd /share/result/%s && slither %s --json -", contractFolder, file)}
	fmt.Println(cmd)

	result, err := performExec(cli, resp, cmd, false)
	if err != nil {
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) perform execution")
	}

	if string(result) == "" {
		cmd = []string{"sh", "-c", fmt.Sprintf("slither %s --etherscan-apikey 4I6DEE2HEKA8SVDQW59VYZSJCX7HQPQ8JK --json -", address)}
		fmt.Println(cmd)
	
		result, err = performExec(cli, resp, cmd, false)
		if err != nil {
			return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) perform execution")
		}
	}

	helper.WriteJSONToFile(string(result), "slitherCheckBefore.json")

	result = []byte(helper.RemoveAfterFirstChar(string(result),"{"))
	result = helper.CleanupJSON(result)
	// result = helper.CleanupJSON(result)

	helper.WriteJSONToFile(string(result), "slither.json")

	var returnResult models.SlitherResultDetail
	if err := json.Unmarshal(result, &returnResult); err != nil {
		return models.SlitherResultDetail{}, helper.MakeError(err, "(slither) unmarshal json")
	}

	return returnResult, nil
}


func GetSlitherSumUp(detail models.SlitherResultDetail, err error) []models.SumUp{
	var sumups []models.SumUp
	if err != nil {
		sumups = append(sumups, models.SumUp{
			Name: "SLITHER ERROR",
			Description: err.Error(),
			Severity: "",
			Location: models.Location{},
		})
		return sumups
	}
	haveIssue := false
	for _, issue := range detail.Results.Detectors{
		if issue.Impact == "informaltional" || issue.Impact == "Informational" || issue.Impact == "Optimization" || issue.Impact == "optimization"{
			continue
		}
		haveIssue = true
		location := getSlitherLocation(issue.Elements)
		sumup := models.SumUp{
			Name: SlitherVulnaClass[issue.Check],
			Description: issue.Description,
			Severity: issue.Impact,
			Location: location,
		}
		sumups = append(sumups, sumup)
	}
	if !haveIssue{
		return []models.SumUp{ models.SumUp{
			Name: "",
			Description: "No error found",
			Severity: "",
			Location: models.Location{},
		},
	}
	}
	return sumups
}

func getSlitherLocation(elements []models.SlitherElement) models.Location{
	var location models.Location
	for _, element := range elements{
		if element.Type == "contract"{
			location.Contract = element.Name
			location.Line = element.SourceMapping.Lines
		} else if element.Type == "function" {
			location.Function = element.TypeSpecificFields.Signature
			location.Line = element.SourceMapping.Lines
		} else if element.Type == "node"{
			location.Line = element.SourceMapping.Lines
		}
	}
	return location
}