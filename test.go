package main

import (
	"fmt"
	"path"
)

func getLastPathComponent(inputPath string) string {
	return path.Base(inputPath)
}

func main1() {
	// Your string
	inputString := "contracts/PoWERC20Factory.sol"

	// Get the last component
	lastComponent := getLastPathComponent(inputString)

	// Print the result
	fmt.Println(lastComponent)
}