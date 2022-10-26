package env

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var nameRegex, _ = regexp.Compile("^[a-zA-Z_]+[a-zA-Z0-9_]*$")
var interpolationRegex, _ = regexp.Compile(`\${+([a-zA-Z_]+[a-zA-Z0-9_]*)\}`)
var multilineToken = "\"\"\""

func stripQuotes(value string) string {
	x := len(value) - 1

	if len(value) > 6 && value[0:3] == multilineToken && value[x-3:x] == multilineToken {
		value = strings.Trim(value, multilineToken)
	} else if value[0] == '\'' && value[x] == '\'' {
		value = strings.Trim(value, "'")
	} else if value[0] == '"' && value[x] == '"' {
		value = strings.Trim(value, "\"")
	}
	return value
}

func parseFile(content []byte) (map[string]string, error) {
	by, _ := os.ReadFile(".env")
	lines := strings.Split(string(by), "\n")

	vars := make(map[string]string)

	varName := ""
	varValue := ""
	multilineOpen := false
	for _, v := range lines {
		if !multilineOpen {
			// only lines beginning with # are omitted
			if len(v) <= 3 || strings.TrimSpace(v)[0] == '#' {
				continue
			}

			// line without value assign is omitted
			endNameIdx := strings.Index(v, "=")
			if endNameIdx == -1 {
				continue
			}

			varName = strings.TrimSpace(strings.ReplaceAll(v[0:endNameIdx], "export ", ""))

			// checking if name follows a standardized format
			if !nameRegex.MatchString(varName) {
				return nil, fmt.Errorf("cant parse file, wrong name format: %v", varName)
			}

			startValueIdx := endNameIdx + 1
			if len(v)-(startValueIdx+3) >= 0 && v[startValueIdx:startValueIdx+3] == multilineToken {
				multilineOpen = true
				startValueIdx += 3
			}

			varValue = v[startValueIdx:]
			if !multilineOpen {
				vars[varName] = varValue
				varName = ""
				varValue = ""
			}

		} else {
			multilineCloseIdx := strings.Index(v, multilineToken)
			// check for multiline termination
			if multilineCloseIdx == -1 {
				varValue += "\\n" + v
			} else {
				varValue += "\\n" + v[:multilineCloseIdx]
				vars[varName] = varValue
				multilineOpen = false
				varName = ""
				varValue = ""
			}
		}
	}

	for k, v := range vars {
		// dont interpolate values not enclosed by double quotes
		if v[0] != '"' || v[len(v)-1] != '"' {
			continue
		}

		val := vars[k]
		foundIdx := interpolationRegex.FindStringIndex(val)
		for len(foundIdx) == 2 {
			varName := val[foundIdx[0]+2 : foundIdx[1]-1]
			if varName == k {
				foundIdx = interpolationRegex.FindStringIndex(v)
				continue
			}
			swap := stripQuotes(vars[varName])
			newValue := val[0:foundIdx[0]] + swap + val[foundIdx[1]:]
			val = newValue
			vars[k] = val
			foundIdx = interpolationRegex.FindStringIndex(val)
		}

	}

	for k, v := range vars {
		vars[k] = stripQuotes(v)
	}

	return vars, nil
}
