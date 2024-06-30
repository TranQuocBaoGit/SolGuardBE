package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func WriteFile(str string, path string){
    
    file, err := os.Create(path)
    CheckError(err)
    defer file.Close()

    file.WriteString(str)
}

func WriteFileExtra(str string, path string) {
	// Open the file in append mode
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	CheckError(err)
	defer file.Close()

	// Write the extra content to the file
	_, err = file.WriteString("\n" + str)
	CheckError(err)
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

func RemoveAfterFirstChar(input string, char string) string {
	index := strings.Index(input, char)
	if index == -1 {
		return input
	}
	return input[index:]
}

func RemoveAfterXChar(input string, char string, x int) string {
	for i := 0; i < x; i++ {
		input = RemoveAfterFirstChar(input, char)
	}
	return input
}

func WriteJSONToFile(data string, filename string) error {
	// Marshal the JSON data
	jsonBytes := []byte(data)

	// Write JSON data to a file
	if err := ioutil.WriteFile(filename, jsonBytes, 0644); err != nil {
		return err
	}

	log.Println("JSON data has been written to", filename)	
	return nil
}

func CleanupJSON(data []byte) []byte {
	var cleanedOutput []byte

	// Loop through the input data and remove invalid characters
	for _, b := range data {
		if b >= 32 && b <= 126 {
			cleanedOutput = append(cleanedOutput, b)
		}
	}

	// Convert the cleaned-up output to a string and trim any leading/trailing whitespace
	cleanedString := strings.TrimSpace(string(cleanedOutput))

	return []byte(cleanedString)
}

func HandleJSON(data interface{}) ([]byte, error) {
	// Use a custom buffer with a larger size
	buffer := bytes.NewBuffer(make([]byte, 1024 * 1024)) // Adjust buffer size as needed
	encoder := json.NewEncoder(buffer)
  
	// Set encoder options to handle large arrays without truncation
	encoder.SetIndent("", "") // Remove indentation for better space efficiency
	// You might consider other options like encoder.SetEscapeHTML(false)
  
	err := encoder.Encode(data)
	if err != nil {
	  return nil, fmt.Errorf("error encoding JSON: %w", err)
	}
	return buffer.Bytes(), nil
  }