package helper

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

func CreateNewContractsFolder() (string, error) {
	number := rand.Intn(1001)

	folderName := fmt.Sprintf("contracts%d", number)

	newFolderPath := filepath.Join("result", folderName)

	_, err := os.Stat(newFolderPath)
	if !os.IsNotExist(err) {
		return "", fmt.Errorf("folder '%s' already exists", folderName)
	}

	err = os.Mkdir(newFolderPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create folder '%s': %v", folderName, err)
	}

	return folderName, nil
}

func DeleteContractsFolder(name string) {

	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		CheckError(fmt.Errorf("folder '%s' does not exist", name))
	}

	err = os.RemoveAll(name)
	if err != nil{
		CheckError(fmt.Errorf("failed to delete folder '%s': %v", name, err))
	}
}