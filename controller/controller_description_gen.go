package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"getContractDeployment/configs"
	"getContractDeployment/helper"
	"getContractDeployment/models"
	"io/ioutil"
	"net/http"
	"strings"
)

const apiURL = "https://api.openai.com/v1/chat/completions"

// RequestPayload defines the structure of the request payload for OpenAI API
type RequestPayload struct {
    Model    string      `json:"model"`
    Messages []Message `json:"messages"`
}

// Message defines the structure of each message in the request payload
type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

// ResponsePayload defines the structure of the response payload from OpenAI API
type ResponsePayload struct {
    Choices []Choice `json:"choices"`
}

// Choice defines the structure of each choice in the response payload
type Choice struct {
    Message Message `json:"message"`
}

func GenDescription(vulneName string) (models.Vulnerability, error) {
	config, err := configs.LoadConfig(".")
	helper.CheckError(err)

    apiKey := config.OPENAI_API_KEY
	// fmt.Print(apiKey)

    messages := []Message{
        {
            Role:    "system",
            Content: "You are a helpful assistant.",
        },
        {
            Role:    "user",
            Content: fmt.Sprintf(`please generate user friendly description for list of solidity vulnerability of smart contract : %s with requirement: - generate for general web3 user which know the basic of the blochchain is but only the basic like what blockchain is and what dApps and smart contract is but in normal term, like smart contract is just another type of transaction and its like a deal you made in real life - generated content should be minimum 3 sentences and at max 7 - generated content should be more practical for non tech user mate? for example "A reentrancy attack happens when an attacker tricks a smart contract into calling a function repeatedly...." can be turn into " an attacker will try to make a blockchain smart contract transaction repeatedly..." -generated content should also describe how this might affect user please return the result as a javascript array with object as for mat: {name:<vulnerability name>, description: <vulnerability description>}. Also don't include all the previous vulnerability and if it's the same vulnerability as before, generate a new description for it.`, vulneName) ,
        },
    }

    requestData := RequestPayload{
        Model:    "gpt-3.5-turbo",
        Messages: messages,
    }

    requestBody, err := json.Marshal(requestData)
    if err != nil {
        // fmt.Println("Error marshaling request data:", err)
        return models.Vulnerability{}, helper.MakeError(err, "(controller_description) Marshal request: ")
    }

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
    if err != nil {
        // fmt.Println("Error creating request:", err)
        return models.Vulnerability{}, helper.MakeError(err, "(controller_description) Creating request: ")
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + apiKey)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        // fmt.Println("Error sending request:", err)
        return models.Vulnerability{}, helper.MakeError(err, "(controller_description) Sending request: ")
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        // fmt.Println("Error reading response body:", err)
        return models.Vulnerability{}, helper.MakeError(err, "(controller_description) Reading respond body: ")
    }

    if resp.StatusCode != http.StatusOK {
        // fmt.Printf("Request failed with status %d: %s\n", resp.StatusCode, body)
        return models.Vulnerability{}, helper.MakeError(fmt.Errorf(resp.Status), "(controller_description) Status fail: ")
    }

    var responsePayload ResponsePayload
    if err := json.Unmarshal(body, &responsePayload); err != nil {
        // fmt.Println("Error unmarshaling response data:", err)
        return models.Vulnerability{}, helper.MakeError(err, "(controller_description) Unmarshaling response data: ")
    }

    if responsePayload.Choices != nil {
        content := responsePayload.Choices[0].Message.Content
            // Remove the surrounding ```javascript and ``` markers
        cleanStr := strings.TrimPrefix(content, "```javascript\n")
        cleanStr = strings.TrimSuffix(cleanStr, "```")

                // Fix the JSON format: replace keys without quotes to keys with quotes
        cleanStr = strings.ReplaceAll(cleanStr, "name:", "\"name\":")
        cleanStr = strings.ReplaceAll(cleanStr, "description:", "\"description\":")
        cleanStr = helper.RemoveTrailingCommas(cleanStr)
        fmt.Print(cleanStr)
        var returnResult []models.Vulnerability
        if err := json.Unmarshal([]byte(cleanStr), &returnResult); err != nil {
            return models.Vulnerability{}, helper.MakeError(err, "(controller_description) Unmarshaling json")
        }
        return returnResult[0], nil
    }

    return models.Vulnerability{}, helper.MakeError(fmt.Errorf("Empty Response"), "(controller_description) ")
}
