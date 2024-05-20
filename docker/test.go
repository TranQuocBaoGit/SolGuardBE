package docker

// func Test(file, contractFolder string, remappingJson bool) {
// 	start := time.Now()

// 	ctx := context.Background()

// 	// Create a Docker client
// 	cli, err := client.NewClientWithOpts(client.FromEnv)
// 	if err != nil {
// 		panic(helper.MakeError(err, "(solhint) new docker client"))
// 	}

// 	result, err := runSolhintContainer(ctx, cli, file, contractFolder, remappingJson)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// jsonData, err := json.Marshal(result)
// 	// if err != nil {
// 	// 	panic(err)
// 	// }

// 	end := time.Since(start)
// 	fmt.Println("Analysis done")
// 	fmt.Println("Time taken: ", end.Seconds())

// 	helper.WriteJSONToFile(result, "solhint.json")
// }

// func runSolhintContainer(ctx context.Context, cli *client.Client, file string, contractFolder string, remappingJSON bool) (string, error) {

// 	currentDir, err := os.Getwd()
// 	if err != nil {
// 		panic(err)
// 	}

// 	hostConfig := createHostConfig(currentDir, "/share")

// 	resp, err := createContainer(ctx, cli, "solhint", true, nil, &hostConfig, "")
// 	if err != nil {
// 		panic(helper.MakeError(err, "(solhint) create container"))
// 	}
// 	defer func() {
// 		err := cli.ContainerRemove(context.Background(), resp.ID, types.ContainerRemoveOptions{Force: true})
// 		if err != nil {
// 			panic(err)
// 		}
// 	}()

// 	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
// 	if err != nil {
// 		panic(helper.MakeError(err, "(solhint) container start"))
// 	}

// 	var setupCmd []string = []string{"sh", "-c", fmt.Sprintf("solhint --init && sed -i 's/\"extends\": \"solhint:default\"/\"extends\": \"solhint:recommended\"/' .solhint.json")}
// 	result1, err := performExec(cli, resp, setupCmd)
// 	fmt.Println(string(result1))
// 	if err != nil {
// 		panic(helper.MakeError(err, "(solhint) perform setup exec"))
// 	}

// 	var tempCmd []string = []string{"sh", "-c", fmt.Sprintf("cat .solhint.json")}
// 	temp, err := performExec(cli, resp, tempCmd)
// 	fmt.Println(string(temp))
// 	if err != nil {
// 		panic(helper.MakeError(err, "(solhint) perform temp exec"))
// 	}

// 	var cmd []string = []string{"sh", "-c", fmt.Sprintf("cd /share/result && solhint %s/%s -f json", contractFolder, file)}
// 	result2, err := performExec(cli, resp, cmd)
// 	if err != nil {
// 		panic(helper.MakeError(err, "(solhint) perform analyze exec"))
// 	}

// 	// var returnResult string
// 	// if err := json.Unmarshal(result, &returnResult); err != nil {
// 	// 	fmt.Print(string(result))
// 	// 	panic(helper.MakeError(err, "(solhint) json unmarshal"))
// 	// }

// 	return helper.RemoveAfterFirstChar(string(result2),"["), nil
// }