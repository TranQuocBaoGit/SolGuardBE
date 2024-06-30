package docker

import (
	"context"
	"fmt"
	"getContractDeployment/helper"
	"io"
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

func createContainer(
	ctx context.Context, 
	cli *client.Client,
	imageName string,
	tty bool,
	cmd []string, 
	hostConfig *container.HostConfig,
	containerName string) (container.CreateResponse, error){
	resp, err := cli.ContainerCreate(
		ctx, 
		&container.Config{
		Image: imageName,
		Tty: tty,
		Cmd: cmd,
		},
		hostConfig,
		nil,
		nil,
		containerName,
	)
	return resp, err
}

func ListRunningContainerIDs() ([]string, error) {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return nil, err
    }

    ctx := context.Background()
    containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: false})
    if err != nil {
        return nil, err
    }
    var containerIDs []string
    for _, container := range containers {
        containerIDs = append(containerIDs, container.ID)
    }

    return containerIDs, nil
}

func RemoveContainerByID(containerID string) error {
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if err != nil {
        return err
    }
    ctx := context.Background()

    if err := cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true}); err != nil {
        return err
    }

    return nil
}

func listImage(){
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		fmt.Println(image.ID)
	}
}

func waitContainer(ctx context.Context, cli *client.Client, resp container.CreateResponse) error{
	// Wait for the container to finish
	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return err
		}
	case <-statusCh:
	}
	return nil
}

func createHostConfig(currentDir string, target string) container.HostConfig{
	hostConfig := container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: currentDir,
				Target: target,
			},
		},
	}
	return hostConfig
}

func retrieveContainerLogs(ctx context.Context, cli *client.Client, resp container.CreateResponse) ([]byte, error) {
	// Retrieve container logs
	out, err := cli.ContainerLogs(context.Background(), resp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil{
		return nil, err
	}
	defer out.Close()

	// Display container logs
	logBytes, err := ioutil.ReadAll(out)
	if err != nil{
		return nil, err
	}
	return logBytes, nil
}

func performExec(cli *client.Client, resp container.CreateResponse, cmd []string, tty bool) ([]byte, error){
	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
		Tty: tty,
	}

	respExec, err := cli.ContainerExecCreate(context.Background(), resp.ID, execConfig)
	if err != nil {
		return nil, helper.MakeError(err, "(general_docker) create execution")
	}
	execResp, err := cli.ContainerExecAttach(context.Background(), respExec.ID, types.ExecStartCheck{})
	if err != nil {
		return nil, helper.MakeError(err, "(general_docker) attach execution")
	}
	// return &execResp, nil
	defer execResp.Close()

	var buf strings.Builder

	_, err = io.Copy(&buf, execResp.Reader)
	if err != nil {
		return nil, helper.MakeError(err, "(general_docker) copy exec resp to string")
	}

	// var buf bytes.Buffer
	// _, err = io.Copy(&buf, execResp.Reader)
	// if err != nil {
	//   return nil, helper.MakeError(err, "(general_docker) copy exec resp to string")
	// }

	// result := []byte(helper.RemoveAfterFirstChar(buf.String(),"{"))
	// resultReturn, err := helper.HandleJSON(result)
	// if err != nil {
	// 	return nil, helper.MakeError(err, "(general_docker) handle json")
	// }

// 	var encoder *json.Encoder
//   // Use a custom buffer with a larger size
// 	buffer := bytes.NewBuffer(make([]byte, 1024 * 1024)) // Adjust buffer size as needed
// 	encoder = json.NewEncoder(buffer)
// 	encoder.SetIndent("", "")

// 	var data interface{}
// 	err = encoder.Encode(data)
// 	if err != nil {
// 		return nil, helper.MakeError(err, "(general_docker) copy exec resp to string")
// 	}


	// result = helper.CleanupJSON(result)

	return []byte(buf.String()), nil
}

// func deleteContainer(ctx context.Context, cli *client.Client, resp container.CreateResponse) (error){
// 	err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
// 	return err
// }