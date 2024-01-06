package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func createContainer(
	ctx context.Context, 
	cli *client.Client, 
	imageName string, 
	cmd []string, 
	hostConfig *container.HostConfig) (container.CreateResponse, error){
	resp, err := cli.ContainerCreate(
		ctx, 
		&container.Config{
		Image: imageName,
		Cmd: cmd,
		},
		hostConfig,
		nil,
		nil,
		"",
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