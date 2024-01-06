package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"getContractDeployment/helper"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type RemappingJSON struct {
	Remappings []string `json:"remappings"`
}

func RunMythrilAnalysis(contractPath string, file string, remappingJSON bool) (string, error) {
	ctx := context.Background()

	// Create a Docker client
	cli, err := client.NewClientWithOpts(client.FromEnv)
	helper.CheckError(err)

	// Run Mythril Docker container to analyze the smart contract
	result, err := runMythrilContainer(ctx, cli, contractPath, file, remappingJSON)

	return result, err
}

func runMythrilContainer(ctx context.Context, cli *client.Client, contractPath string, file string, remappingJSON bool) (string, error) {
	// Convert Windows path to Linux format
	contractPath = strings.ReplaceAll(contractPath, "\\", "/")

	// Create command
	var cmd []string = []string{"analyze"}
	cmd = append(cmd, "/mnt/contracts/" + file, "-o", "jsonv2")
	if remappingJSON {
		cmd = append(cmd, "--solc-json", "/mnt/remapping/mythril_remappings.json")
	}
	fmt.Println(cmd)
	hostConfig := container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: contractPath,
				Target: "/mnt",
			},
		},
	}
	resp, err := createContainer(ctx, cli, "mythril/myth", cmd, &hostConfig)
	if err != nil {
		return "", err
	}

	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", err
	}

	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return "", err
		}
	case <-statusCh:
	}

	// Retrieve container logs
	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return "", err
	}

	defer out.Close()

	// Display container logs
	logBytes, err := ioutil.ReadAll(out)
	if err != nil {
		return "", err
	}

	result := string(logBytes)
	fmt.Println("Mythril Analysis Logs:")
	fmt.Println(result)

	return result, nil
}

func CreateMythrilMappingJson(data map[string]interface{}) error {
	results := make(map[string]map[string]string)

	for contract, content := range data {
		if contract == "Main Contract"{
			continue
		}
        sourcePath := fmt.Sprintf("./result/contracts/%s", contract)
        helper.WriteFile(content.(string), sourcePath)

        allImportPath := helper.FindImportPath(content.(string))
        importReplacement := make(map[string]string)
        for _, eachImportPath := range allImportPath{
            if strings.HasPrefix(eachImportPath, "@openzeppelin"){
                replacePath := "/mnt/contracts/openzeppelin/" + helper.GetLastFilePath(eachImportPath)
                importReplacement[eachImportPath] = replacePath
            } else if strings.HasPrefix(contract, "openzeppelin"){
				replacePath := helper.GetLastFilePath(eachImportPath)
                importReplacement[eachImportPath] = replacePath
			} else {
                replacePath := "/mnt/contracts/" + helper.GetLastFilePath(eachImportPath)
                importReplacement[eachImportPath] = replacePath
            }
        }
        results[contract] = importReplacement
    }

	var remappingJSON RemappingJSON

    for _, importReplacement := range results{
        for importPath, replacement := range importReplacement {
            result := fmt.Sprint(importPath, "=", replacement)
			remappingJSON.Remappings = append(remappingJSON.Remappings, result)
        }
    }
	jsonData, err := json.Marshal(remappingJSON)
	if err != nil {
		return err
	}

	helper.WriteFile(string(jsonData), "./result/remapping/mythril_remappings.json")
	return nil
}