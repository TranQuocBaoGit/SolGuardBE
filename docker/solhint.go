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

func RunSolHintAnalysisWithTimeOut(file string, contractFolder string, remappingJSON bool) (models.SolhintResultDetail, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Channel to receive the result
	resultChan := make(chan struct {
		result models.SolhintResultDetail
		err    error
	}, 1)

	// Run the fetchMythrilResult function in a separate goroutine
	go func() {
		result, err := RunSolhintAnalysis(file, contractFolder, remappingJSON)
		resultChan <- struct {
			result models.SolhintResultDetail
			err    error
		}{result, err}
	}()

	// Use a select statement to wait for either the result or the context timeout
	select {
	case res := <-resultChan:
		return res.result, res.err
	case <-ctx.Done():
		return models.SolhintResultDetail{}, fmt.Errorf("Solhint time out")
	}
}

func RunSolhintAnalysis(file string, contractFolder string, remappingJSON bool) (models.SolhintResultDetail, error) { //(string, error)
	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil{
		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) new docker client")
	}

	result, err := runSolhintContainer(ctx, cli, file, contractFolder, remappingJSON)
	if err != nil{
		return models.SolhintResultDetail{}, err
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return models.SolhintResultDetail{}, err
	}

	helper.WriteJSONToFile(string(jsonData), "solhint.json")

	return result, nil
}

func runSolhintContainer(ctx context.Context, cli *client.Client, file string, contractFolder string, remappingJSON bool) (models.SolhintResultDetail, error) {

	currentDir, err := os.Getwd()
	if err != nil {
		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) get directory")
	}

	hostConfig := createHostConfig(currentDir, "/share")

	resp, err := createContainer(ctx, cli, "solhint", true, nil, &hostConfig, "")
	if err != nil {
		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) create container")
	}
	defer func(){
		err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}()

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) container start")
	}

	newConfigFile := `"extends": "solhint:recommended","plugins": [],"rules": {"avoid-call-value": "warn","avoid-low-level-calls": "warn","avoid-throw": "warn","avoid-tx-origin": "error","check-send-result": "warn","func-visibility": "error","multiple-sends": "off","no-complex-fallback": "off","no-inline-assembly": "warn","not-rely-on-block-hash": "error","not-rely-on-time": "error","reentrancy": "error","state-visibility": "warn","avoid-suicide": "error","avoid-sha3": "warn"}`

	var setupCmd []string = []string{"sh", "-c", fmt.Sprintf(`solhint --init && sed -i 's/"extends": "solhint:default"/%s/' .solhint.json && cat .solhint.json`, newConfigFile)}
	// fmt.Println(cmd)

	_, err = performExec(cli, resp, setupCmd)
	if err != nil {
		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) perform setup execution")
	}
	// fmt.Print("Solhint haiz is: ", string(haiz))

	var cmd []string = []string{"sh", "-c", fmt.Sprintf("solhint /share/result/%s/%s -f json", contractFolder, file)}
	result, err := performExec(cli, resp, cmd)
	if err != nil {
		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) perform analyze execution")
	}

	result = []byte(helper.RemoveAfterFirstChar(string(result),"["))
	// fmt.Print("Solhint result is: ", string(result))

	// if rune(string(result)[0]) != rune('[') {
	// 	return models.SolhintResultDetail{}, err
	// }


	var returnResult models.SolhintResultDetail
	if err := json.Unmarshal(result, &returnResult); err != nil {
		return models.SolhintResultDetail{}, helper.MakeError(err, "(solhint) json unmarshal")
	}

	return returnResult, nil
}


func GetSolhintSumUp(detail models.SolhintResultDetail, err error) []models.SumUp{
	var sumups []models.SumUp
	if err != nil {
		sumups = append(sumups, models.SumUp{
			Name: "SOLHINT ERROR",
			Description: "Solhint fail to analyze contract",
			Severity: "",
		})
		return sumups
	}
	for _, issue := range detail{
		if issue.Issues != nil{
			seveDefine := ""
			switch issue.Issues.Severity{
			case "off":
				seveDefine = "Low"
				break
			case "warn":
				seveDefine = "Medium"
				break
			case "error":
				seveDefine = "High"
				break
			default:
				seveDefine = "Low"
				break
			} 
			sumup := models.SumUp{
				Name: issue.Issues.RuleID,
				Description: issue.Issues.Message,
				Severity: seveDefine,
			}
			sumups = append(sumups, sumup)
		}
	}
	return sumups
}