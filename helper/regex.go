package helper

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func UnescapeJSON(jsonString string) (string, error) {
	// Replace \" with "
	re := regexp.MustCompile(`\\\"`)
	unescaped := re.ReplaceAllString(jsonString, `"`)

	// Replace \\ with \
	re = regexp.MustCompile(`\\\\`)
	unescaped = re.ReplaceAllString(unescaped, `\`)

	return unescaped, nil
}

func PreprocessJSON(jsonData string) string {
	// Replace invalid characters with spaces
	validJSONData := strings.Map(func(r rune) rune {
		if r < 32 {
			return ' ' // Replace control characters with space
		}
		return r
	}, jsonData)

	return validJSONData
}

func FindImportPath(solidityCode string) []string {
	importRegex := regexp.MustCompile(`import\s+[^;]+;`)

	allImports := importRegex.FindAllStringSubmatch(solidityCode, -1)

	var result []string
	for _, eachImport := range allImports {
		pathRegex := regexp.MustCompile(`"+[^"]+"`)
		path := pathRegex.FindStringSubmatch(eachImport[0])
		result = append(result, path[0][1:len(path[0])-1])
	}

	return result
}

func ReplacePathWithFilename(input string) string {
	// Regular expression to match import statements with paths
	re := regexp.MustCompile(`import\s*"(.*?)"`)

	// Replace each match with only the file name
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract the path from the match
		submatches := re.FindStringSubmatch(match)
		if len(submatches) > 1 {
			path := submatches[1]

			// Find the last occurrence of "/" and extract the filename
			lastSlashIndex := len(path) - 1
			for i := len(path) - 1; i >= 0; i-- {
				if path[i] == '/' {
					lastSlashIndex = i
					break
				}
				if i == 0 {
					lastSlashIndex = -1
					break
				}
			}

			filename := path[lastSlashIndex+1:]

			// Replace the path with the filename
			return fmt.Sprintf(`import "./%s"`, filename)
		}

		return match
	})

	return result
}

func ChangeFromSolToJson(file string) string {
	// for i, path := range filePaths {
	if strings.HasSuffix(file, ".sol") {
		file = strings.TrimSuffix(file, ".sol") + ".json"
	}
	// }
	return file
}

func ChangeSolidityVersion(sourceCode, targetVersion string) (string, error) {
    // Define the regular expression to match the Solidity version pragma
    re := regexp.MustCompile(`pragma solidity\s+[^;]+;`)

    // Define the new pragma version string
    newPragma := fmt.Sprintf("pragma solidity %s;", targetVersion)

    // Replace the old pragma version with the new one
    modifiedSourceCode := re.ReplaceAllString(sourceCode, newPragma)

    // Check if any replacements were made (if no matches, it might be an error or simply no pragma found)
    if strings.Contains(modifiedSourceCode, newPragma) {
        return modifiedSourceCode, nil
    } else {
        return "", fmt.Errorf("no Solidity version pragma found in the source code")
    }
}

func EnsureCommas(jsonStr string) string {
    var buffer bytes.Buffer

    for i := 0; i < len(jsonStr); i++ {
        buffer.WriteByte(jsonStr[i])

        if jsonStr[i] == '}' && i+1 < len(jsonStr) && jsonStr[i+1] != ',' && jsonStr[i+1] != ']' && jsonStr[i+1] != ' ' {
            buffer.WriteByte(',')
        }
    }

    return buffer.String()
}

// Function to remove trailing commas before closing braces
func RemoveTrailingCommas(jsonStr string) string {
    re := regexp.MustCompile(`,\s*}`)
    return re.ReplaceAllString(jsonStr, "}")
}