package controller

import (
	"context"
	"getContractDeployment/configs"
	"getContractDeployment/docker"
	"getContractDeployment/models"
	"slices"
	"time"
)

func StandardizeResult(toolsResult []models.ToolResult) (models.StandardizeResults, error) {

	// vulneToServerity := make(map[string]string)
	standardizeMap := make(map[string]models.StandardizeResult)
	// var standardizeResults []models.StandardizeResult
	for _, toolResult := range toolsResult {
		
		for _, sumUp := range toolResult.SumUps {
			vulne, severity := IdentifyVulnerability(sumUp, toolResult.ToolName)
			if vulne == "" {
				continue
			}
			oneStandardize, exist := standardizeMap[vulne]
			if !exist {
				description, err := GetDescription(vulne)
				if err != nil {
					return models.StandardizeResults{}, err
				}
				standardizeMap[vulne] = models.StandardizeResult{
					Name: vulne,
					Severity: severity,
					GeneralDescription: description,
					AdvancedDescription: append([]string{}, sumUp.Description),
					Tools: []string{toolResult.ToolName},
					NoOccurrence: 1,
					Locations: []models.Location{sumUp.Location},
				}
			} else {
				newSeverity := GetHighestSeverity(oneStandardize.Severity, severity)
				oneStandardize.Severity = newSeverity
				if !slices.Contains(oneStandardize.Tools, toolResult.ToolName){
					oneStandardize.Tools = append(oneStandardize.Tools, toolResult.ToolName)
					oneStandardize.AdvancedDescription = append(oneStandardize.AdvancedDescription, sumUp.Description)
				}
				occurBefore := CompareLocation(oneStandardize.Locations, sumUp.Location)
				if !occurBefore{
					oneStandardize.NoOccurrence++
					oneStandardize.Locations = append(oneStandardize.Locations, sumUp.Location)
				}
				standardizeMap[vulne] = oneStandardize
			}
		}
	}

	standardizeResults := models.StandardizeResults{
		NoError: 0,
		Result:  []models.StandardizeResult{},
	}
	for _, oneStandardize := range standardizeMap {
		standardizeResults.NoError += oneStandardize.NoOccurrence
		standardizeResults.Result = append(standardizeResults.Result, oneStandardize)
	}

	return standardizeResults, nil
}

func GetHighestSeverity(sever1, sever2 string) string {
	if sever1 == "High" || sever2 == "High" {
		return "High"
	}
	if sever1 == "Medium" || sever2 == "Medium" {
		return "Medium"
	}
	return "Low"
}

func IdentifyVulnerability(sumup models.SumUp, tool string) (string, string) {
	if tool == "mythril" {
		return docker.MythrilStandardize(sumup)
	} else if tool == "slither" {
		return docker.SlitherStandardize(sumup)
	} else if tool == "solhint" {
		return docker.SolhintStandardize(sumup)
	} else if tool == "honeybadger" {
		return docker.HoneyBadgerStandardize(sumup)
	}
	return "", ""
}

func CheckSameLocation(location1, location2 models.Location) bool{
	if location1.Function == location2.Function{
		return true
	}
	if len(location1.Line) > len(location2.Line){
		return slices.Contains(location1.Line, location2.Line[0])
	}
	return slices.Contains(location2.Line, location1.Line[0])
}

func CompareLocation(currentLocations []models.Location, newLocation models.Location) bool {
	occurBefore := false
	for i, oldLocation := range currentLocations{
		if CheckSameLocation(oldLocation, newLocation) {
			occurBefore = true
			if len(oldLocation.Line) > len(newLocation.Line) {
				currentLocations[i].Line = newLocation.Line
			}
			if oldLocation.Function == newLocation.Function{
				continue
			}
			if oldLocation.Contract == "" && newLocation.Contract != ""{
				currentLocations[i].Contract = newLocation.Contract
			}
			if oldLocation.Function == "" && newLocation.Function != ""{
				currentLocations[i].Function = newLocation.Function
			}
		}
	}
	return occurBefore
}

func GetDescription(vulneName string) (string, error) {
	config, err := configs.LoadConfig(".")
	if err != nil{
		return "", err
	}
	client := configs.ConnectDB(config)
	ctx, _ := context.WithTimeout(context.Background(), 3600*time.Second)
	vulneCol := client.Database(config.DATABASE_NAME).Collection("vulnerability")

	description, err := GetVulnerabilityDescription(ctx, vulneCol, vulneName)
	if err != nil{
		return "", err
	}

	if description == "" {
		vulnerability, err := GenDescription(vulneName)
		if err != nil{
			return "", err
		}
		AddVulnerability(ctx, vulneCol, vulnerability)
		return vulnerability.Description, nil
	}
	return description, nil
}