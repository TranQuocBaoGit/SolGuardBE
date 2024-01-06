package helper

import (
	"os"
	"strings"
)

func WriteFile(str string, path string){
    
    file, err := os.Create(path)
    CheckError(err)
    defer file.Close()

    file.WriteString(str)
}

func GetLastFilePath(inputPath string) string {
	// Split the string by "/"
	parts := strings.Split(inputPath, "/")

	// Get the last part of the slice
	return parts[len(parts)-1]
}

func GetPathToFile(inputPath string) string {
	parts := strings.LastIndex(inputPath, "/")
	return inputPath[0:parts+1]
}
// ../../utils/StorageSlot.sol