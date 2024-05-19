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

type RemappingJSON struct {
	Remappings []string `json:"remappings"`
}

func RunMythrilAnalysis(file string, contractFolder string, remappingJSON bool) (models.MythrilResultDetail, error) {
	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil{
		return models.MythrilResultDetail{}, helper.MakeError(err, "(mythril) new docker client")
	}

	// Run Mythril Docker container to analyze the smart contract
	result, err := runMythrilContainer(ctx, cli, contractFolder, file, remappingJSON)
	if err != nil{
		return models.MythrilResultDetail{}, err
	}

	result = []byte(helper.RemoveAfterFirstChar(string(result),"{"))
	result = helper.CleanupJSON(result)

	// Convert result to models
	var returnResult models.MythrilResultDetail
	err = json.Unmarshal([]byte(result), &returnResult)
	if err != nil {
		return models.MythrilResultDetail{}, helper.MakeError(err, "(mythril) json unmarshal")
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return models.MythrilResultDetail{}, err
	}

	helper.WriteJSONToFile(string(jsonData), "mythril.json")

	if !returnResult.Success{
		return models.MythrilResultDetail{
			Error: helper.MakeError(err, "(mythril) failed to analyze contract"),
			Issues: nil,
			Success: false,
		}, nil
	}

	return returnResult, nil
}

func runMythrilContainer(ctx context.Context, cli *client.Client, contractFolder string, file string, remappingJSON bool) ([]byte, error) {

	currentDir, err := os.Getwd()
	if err != nil {
		return nil, helper.MakeError(err, "(mythril) get directory")
	}

	hostConfig :=  createHostConfig(currentDir, "/mnt")

	path := fmt.Sprintf("/mnt/result/%s/%s", contractFolder, file)
	// Create command
	var cmd []string = []string{"analyze"}
	cmd = append(cmd, path, "-o", "json") //"-o", "jsonv2"
	if remappingJSON {
		cmd = append(cmd, "--solc-json", "/mnt/result/remapping/mythril_remappings.json")
	}
	
	resp, err := createContainer(ctx, cli, "mythril/myth", false, cmd, &hostConfig, "")
	if err != nil {
		return nil, helper.MakeError(err, "(mythril) create container")
	}
	defer func(){
		err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
		if err != nil {
			panic(err)
		}
	}()

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, helper.MakeError(err, "(mythril) container start")
	}

	// Wait for the container to finish
	err = waitContainer(ctx, cli, resp)
	if err != nil {
		return nil, helper.MakeError(err, "(mythril) wait container finish")
	}

	// Retrieve container logs
	result, err := retrieveContainerLogs(ctx, cli, resp)
	if err != nil {
		return nil, helper.MakeError(err, "(mythril) retrieve container logs")
	}
	// result = helper.RemoveAfterFirstChar(result, "{")
	
	return result, nil
}

// func CreateMythrilMappingJson(data []models.ContractContent, contractFolder string) error {
// 	results := make(map[string]map[string]string)

// 	for _, content := range data {
// 		if content.ContractName == "Main Contract"{
// 			continue
// 		}
//         sourcePath := fmt.Sprintf("./result/%s/%s",contractFolder, contract)
//         helper.WriteFile(content.(string), sourcePath)

//         allImportPath := helper.FindImportPath(content.(string))
//         importReplacement := make(map[string]string)
//         for _, eachImportPath := range allImportPath{
//             if strings.HasPrefix(eachImportPath, "@openzeppelin"){
//                 replacePath := fmt.Sprintf("/mnt/%s/openzeppelin/",contractFolder) + helper.GetLastFilePath(eachImportPath) // "/mnt/contracts/openzeppelin/"
//                 importReplacement[eachImportPath] = replacePath
//             } else if strings.HasPrefix(contract, "openzeppelin"){
// 				replacePath := helper.GetLastFilePath(eachImportPath)
//                 importReplacement[eachImportPath] = replacePath
// 			} else {
//                 replacePath := fmt.Sprintf("/mnt/%s",contractFolder) + helper.GetLastFilePath(eachImportPath)
//                 importReplacement[eachImportPath] = replacePath
//             }
//         }
//         results[contract] = importReplacement
//     }

// 	var remappingJSON RemappingJSON

//     for _, importReplacement := range results{
//         for importPath, replacement := range importReplacement {
//             result := fmt.Sprint(importPath, "=", replacement)
// 			remappingJSON.Remappings = append(remappingJSON.Remappings, result)
//         }
//     }
// 	jsonData, err := json.Marshal(remappingJSON)
// 	if err != nil {
// 		return err
// 	}

// 	helper.WriteFile(string(jsonData), "./result/remapping/mythril_remappings.json")
// 	return nil
// }

// func mythrilDataHandler(){

// }

func GetMythrilSumUp(detail models.MythrilResultDetail) []models.SumUp{
	var sumups []models.SumUp
	for _, issue := range detail.Issues{
		sumup := models.SumUp{
			Name: MythrilVulnaClass[issue.SwcID],
			Description: issue.Description,
			Severity: issue.Severity,
		}
		sumups = append(sumups, sumup)
	}
	return sumups
}

