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

func performExec(cli *client.Client, resp container.CreateResponse, cmd []string) ([]byte, error){
	execConfig := types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
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

	// decoder := json.NewDecoder(execResp.Reader)

	// return decoder, nil

	var buf strings.Builder

	// copy execResp to string buf
	_, err = io.Copy(&buf, execResp.Reader)
	if err != nil {
		return nil, helper.MakeError(err, "(general_docker) copy exec resp to string")
	}

	result := helper.CleanupJSON([]byte(buf.String()))

	return result, nil
}

// func deleteContainer(ctx context.Context, cli *client.Client, resp container.CreateResponse) (error){
// 	err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
// 	return err
// }