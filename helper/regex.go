package helper

import (
	"fmt"
	"regexp"
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